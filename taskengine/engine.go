package taskengine

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/MAD-py/go-taskengine/taskengine/store"
)

type Engine struct {
	mu sync.Mutex

	ctx             context.Context
	shutdownTimeout time.Duration

	supervisors map[string]*WorkerSupervisor

	store  store.Store
	logger Logger
}

func (e *Engine) Run() error {
	ctxSignal, cancelSignal := signal.NotifyContext(e.ctx, os.Interrupt, syscall.SIGTERM)
	defer cancelSignal()

	e.Start()
	<-ctxSignal.Done()
	return e.Shutdown()
}

func (e *Engine) Start() {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.logger.Info("Starting Task Engine...")
	for _, s := range e.supervisors {
		err := e.store.UpdateTaskStatus(
			e.ctx, s.worker.task.name, store.TaskStatusRunning,
		)
		if err != nil {
			e.logger.Errorf(
				"Failed run task '%s': %v", s.worker.task.name, err,
			)
			continue
		}
		s.Start(e.ctx)
	}
	e.logger.Infof("Task Engine started with %d supervisors", len(e.supervisors))
}

func (e *Engine) Shutdown() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.logger.Info("Shutting down Task Engine...")
	var wg sync.WaitGroup
	for _, s := range e.supervisors {
		wg.Add(1)
		go func(s *WorkerSupervisor) {
			defer wg.Done()
			s.Shutdown()

			err := e.store.UpdateTaskStatus(
				e.ctx, s.worker.task.name, store.TaskStatusIdle,
			)
			if err != nil {
				e.logger.Errorf(
					"Failed to update task '%s' status: %v",
					s.worker.task.name, err,
				)
			}
		}(s)
	}

	done := make(chan struct{})
	go func() { defer close(done); wg.Wait() }()

	select {
	case <-done:
		e.logger.Info("Task Engine shutdown complete")
		return nil
	case <-time.After(e.shutdownTimeout):
		e.logger.Errorf("Task Engine shutdown timed out after %s", e.shutdownTimeout)
		return errors.New("shutdown timed out")
	}
}

func (e *Engine) StartTask(name string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	for _, s := range e.supervisors {
		if s.worker.task.name == name {
			s.Start(e.ctx)
			return nil
		}
	}
	e.logger.Warnf("Task %s not found", name)
	return errors.New("task not found")
}

func (e *Engine) ShutdownTask(name string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	for _, s := range e.supervisors {
		if s.worker.task.name == name {
			done := make(chan struct{})
			go func() {
				defer close(done)
				s.Shutdown()

				err := e.store.UpdateTaskStatus(
					e.ctx, s.worker.task.name, store.TaskStatusIdle,
				)
				if err != nil {
					e.logger.Errorf(
						"Failed to update task '%s' status: %v",
						s.worker.task.name, err,
					)
				}
			}()

			select {
			case <-done:
				e.logger.Infof("Task %s shutdown complete", name)
				return nil
			case <-time.After(e.shutdownTimeout):
				e.logger.Errorf("Shutdown timed out for task %s after %s", name, e.shutdownTimeout)
				return errors.New("shutdown timed out for task: " + name)
			}
		}
	}
	e.logger.Warnf("Task %s not found", name)
	return errors.New("task not found")
}

func (e *Engine) RegisterTask(
	task *Task,
	policy workerPolicy,
	trigger Trigger,
	catchUpEnabled bool,
	maxExecutionLag int,
) error {
	e.mu.Lock()
	if _, exists := e.supervisors[task.name]; exists {
		e.logger.Warnf("Task '%s' is already registered", task.name)
		return nil
	}
	e.mu.Unlock()

	e.logger.Infof("Registering task '%s' with policy '%s'", task.name, policy)

	exists, err := e.store.TaskExists(e.ctx, task.name)
	if err != nil {
		return err
	}

	if exists {
		err := e.validateTaskSettings(task.name, task.jobName, policy, trigger)
		if err != nil {
			return err
		}
	} else {
		err := e.store.SaveTask(e.ctx, task.name, &store.TaskSettings{
			Job:     task.jobName,
			Policy:  policy.String(),
			Trigger: trigger.String(),
		})
		if err != nil {
			return err
		}
	}

	dispatcher := newDispatcher(maxExecutionLag)
	worker := newWorker(task, dispatcher, policy, e.logger)
	scheduler := newScheduler(trigger, dispatcher, catchUpEnabled, e.logger)

	ws := newWorkerSupervisor(worker, scheduler, dispatcher, e.logger)

	e.mu.Lock()
	e.supervisors[task.name] = ws
	e.mu.Unlock()

	err = e.store.UpdateTaskStatus(e.ctx, task.name, store.TaskStatusRegistered)
	if err != nil {
		e.mu.Lock()
		delete(e.supervisors, task.name)
		e.mu.Unlock()
		return err
	}

	e.logger.Infof("Task '%s' registered successfully", task.name)

	return nil
}

func (e *Engine) validateTaskSettings(
	taskName, jobName string, policy workerPolicy, trigger Trigger,
) error {
	taskSettings, err := e.store.GetTaskSettings(e.ctx, taskName)
	if err != nil {
		return err
	}

	if taskSettings.Job != jobName {
		return ErrorJobNameMismatch
	}

	if taskSettings.Policy != policy.String() {
		return ErrorPolicyMismatch
	}

	if taskSettings.Trigger != trigger.String() {
		return ErrorTriggerMismatch
	}

	return nil
}

func (e *Engine) RemoveTask(name string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if supervisor, exists := e.supervisors[name]; exists {
		supervisor.Shutdown()
		delete(e.supervisors, name)

		err := e.store.UpdateTaskStatus(e.ctx, name, store.TaskStatusIdle)
		if err != nil {
			return err
		}

		e.logger.Infof("Task '%s' removed successfully", name)
		return nil
	}

	e.logger.Warnf("Task %s not found", name)
	return errors.New("task not found")
}

func New(store store.Store) (*Engine, error) {
	ctx := context.Background()
	if err := store.CreateStores(ctx); err != nil {
		return nil, err
	}

	return &Engine{
		ctx:             ctx,
		store:           store,
		logger:          &stdLogger{},
		supervisors:     make(map[string]*WorkerSupervisor),
		shutdownTimeout: 30 * time.Second, // Default shutdown timeout
	}, nil
}
