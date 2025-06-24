package taskengine

import (
	"context"
	"errors"
	"time"
)

type Job = func(ctx *Context) error

type Task struct {
	job     Job
	name    string
	logger  Logger
	timeout time.Duration
}

func (t *Task) Name() string { return t.name }

func (t *Task) Execute(parentCtx context.Context, tick *Tick) {
	defer func() {
		if r := recover(); r != nil {
			t.logger.Errorf("PANIC in Task '%s' job: %v", t.name, r)
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
	if err := t.job(&ctxTask); err != nil {
		t.logger.Errorf("Task '%s' failed: %v", t.name, err)

	} else {
		t.logger.Infof("Task '%s' completed successfully", t.name)
	}
}

func NewTask(name string, job Job, timeout time.Duration) (*Task, error) {
	if name == "" {
		return nil, errors.New("task name must be non-empty")
	}

	if job == nil {
		return nil, errors.New("job must be non-nil")
	}

	return &Task{
		job:     job,
		name:    name,
		logger:  &stdLogger{},
		timeout: timeout,
	}, nil
}
