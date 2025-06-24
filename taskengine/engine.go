package taskengine

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Engine struct {
	mu sync.Mutex

	ctx             context.Context
	shutdownTimeout time.Duration

	supervisors []*WorkerSupervisor

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
		go func(s *WorkerSupervisor) { defer wg.Done(); s.Shutdown() }(s)
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
			go func() { defer close(done); s.Shutdown() }()

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
) {
	e.logger.Infof("Registering task '%s' with policy '%s'", task.Name(), policy)

	dispatcher := newDispatcher(maxExecutionLag)
	worker := newWorker(task, dispatcher, policy, e.logger)
	scheduler := newScheduler(trigger, dispatcher, catchUpEnabled, e.logger)

	ws := newWorkerSupervisor(worker, scheduler, dispatcher, e.logger)

	e.mu.Lock()
	defer e.mu.Unlock()
	e.supervisors = append(e.supervisors, ws)

	e.logger.Infof("Task '%s' registered successfully", task.Name())
}

func (e *Engine) RemoveTask(name string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	for i, s := range e.supervisors {
		if s.worker.task.name == name {
			e.supervisors = append(e.supervisors[:i], e.supervisors[i+1:]...)
			e.logger.Infof("Task '%s' removed successfully", name)
			return nil
		}
	}
	e.logger.Warnf("Task %s not found", name)
	return errors.New("task not found")
}

func New() *Engine {
	return &Engine{
		ctx:             context.Background(),
		logger:          &stdLogger{},
		supervisors:     make([]*WorkerSupervisor, 0),
		shutdownTimeout: 30 * time.Second, // Default shutdown timeout
	}
}
