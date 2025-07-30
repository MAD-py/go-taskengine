package taskengine

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"time"

	"github.com/MAD-py/go-taskengine/taskengine/store"
)

type Job = func(ctx *Context) error

type Task struct {
	name string

	job     Job
	jobName string

	logger  Logger
	timeout time.Duration

	store store.Store
}

func (t *Task) Name() string { return t.name }

func (t *Task) Execute(parentCtx context.Context, tick *Tick) {
	startTime := time.Now()

	defer func() {
		if r := recover(); r != nil {
			endTime := time.Now()
			duration := endTime.Sub(startTime)

			t.logger.Errorf("PANIC in Task '%s' job: %v", t.name, r)

			err := t.store.SaveExecution(
				t.name,
				&store.ExecutionInfo{
					StartTime: startTime,
					EndTime:   endTime,
					Duration:  duration,
					Status:    store.ExecutionStatusPanic,
					ErrorMsg:  fmt.Sprintf("PANIC: %v", r),
				},
			)

			if err != nil {
				t.logger.Errorf(
					"Failed to save execution info for task '%s': %v",
					t.name, err,
				)
			}
		}
	}()

	var ctx context.Context
	var cancel context.CancelFunc

	if t.timeout > 0 {
		ctx, cancel = context.WithTimeout(parentCtx, t.timeout)
	} else {
		ctx, cancel = context.WithCancel(parentCtx)
	}

	defer cancel()

	ctxTask := Context{
		ctx:      ctx,
		tick:     tick,
		logger:   t.logger,
		taskName: t.name,
	}

	t.logger.Infof("Executing Task '%s'", t.name)

	err := t.job(&ctxTask)
	endTime := time.Now()
	duration := endTime.Sub(startTime)

	if err != nil {
		t.logger.Errorf("Task '%s' failed: %v", t.name, err)

		err := t.store.SaveExecution(
			t.name,
			&store.ExecutionInfo{
				StartTime: startTime,
				EndTime:   endTime,
				Duration:  duration,
				Status:    store.ExecutionStatusError,
				ErrorMsg:  err.Error(),
			},
		)
		if err != nil {
			t.logger.Errorf(
				"Failed to save execution info for task '%s': %v",
				t.name, err,
			)
		}
		return
	}

	t.logger.Infof("Task '%s' completed successfully", t.name)

	err = t.store.SaveExecution(
		t.name,
		&store.ExecutionInfo{
			StartTime: startTime,
			EndTime:   endTime,
			Duration:  duration,
			Status:    store.ExecutionStatusSuccess,
		},
	)
	if err != nil {
		t.logger.Errorf(
			"Failed to save execution info for task '%s': %v",
			t.name, err,
		)
	}
}

func NewTask(name string, job Job, options ...taskOption) (*Task, error) {
	if name == "" {
		return nil, errors.New("task name must be non-empty")
	}

	if job == nil {
		return nil, errors.New("job must be non-nil")
	}

	jobPtr := reflect.ValueOf(job).Pointer()
	jobName := runtime.FuncForPC(jobPtr).Name()

	task := &Task{
		job:     job,
		name:    name,
		jobName: jobName,
		logger:  newLogger(name),
	}

	for _, opt := range options {
		opt(task)
	}

	return task, nil
}

type taskOption func(*Task)

func WithTaskLogger(logger Logger) taskOption {
	return func(t *Task) {
		t.logger = logger
	}
}

func WithTimeout(timeout time.Duration) taskOption {
	return func(t *Task) {
		t.timeout = timeout
	}
}
