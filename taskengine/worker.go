package taskengine

import (
	"context"
	"sync"
	"sync/atomic"
)

type workerPolicy int

const (
	WorkerPolicyParallel workerPolicy = iota
	WorkerPolicySerial
	WorkerPolicySkipIfBusy
)

func (p workerPolicy) String() string {
	switch p {
	case WorkerPolicyParallel:
		return "parallel"
	case WorkerPolicySerial:
		return "serial"
	case WorkerPolicySkipIfBusy:
		return "skip_if_busy"
	default:
		return "unknown"
	}
}

type workerState int

const (
	workerIdle workerState = iota
	workerRunning
)

type Worker struct {
	wg sync.WaitGroup

	task *Task

	dispatcher *Dispatcher

	policy workerPolicy

	state   atomic.Value
	running atomic.Bool

	logger Logger
}

func (w *Worker) Status() workerState { return w.state.Load().(workerState) }

func (w *Worker) Run(ctx context.Context) {
	if w.state.Load().(workerState) != workerIdle {
		return
	}

	defer w.state.Store(workerIdle)
	w.state.Store(workerRunning)

	for {
		select {
		case tick, ok := <-w.dispatcher.Dequeue():
			if !ok {
				w.logger.Infof(
					"Worker for task '%s' stopped: dispatcher queue closed",
					w.task.Name(),
				)
				w.wg.Wait()
				return
			}

			switch w.policy {
			case WorkerPolicyParallel:
				// TODO: maximum concurrency limit handling
				w.wg.Add(1)
				go func() { defer w.wg.Done(); w.task.Execute(ctx, tick) }()
			case WorkerPolicySerial:
				// TODO: timeout handling for serial execution
				w.task.Execute(ctx, tick)
			case WorkerPolicySkipIfBusy:
				if w.running.CompareAndSwap(false, true) {
					w.wg.Add(1)
					go func() {
						defer w.running.Store(false)
						defer w.wg.Done()
						w.task.Execute(ctx, tick)
					}()
				} else {
					w.logger.Warnf(
						"Skipping execution of task '%s': already running",
						w.task.Name(),
					)
				}
			}
		case <-ctx.Done():
			w.logger.Infof("Shutting down worker for task '%s'", w.task.Name())
			w.wg.Wait()
			w.logger.Infof("Worker for task '%s' shutdown complete", w.task.Name())
			return
		}
	}
}

func newWorker(
	task *Task,
	dispatcher *Dispatcher,
	policy workerPolicy,
	logger Logger,
) *Worker {
	w := &Worker{
		task:       task,
		policy:     policy,
		logger:     logger,
		dispatcher: dispatcher,
	}
	w.state.Store(workerIdle)
	return w
}
