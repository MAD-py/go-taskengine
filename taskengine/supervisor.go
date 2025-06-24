package taskengine

import (
	"context"
	"sync"
	"sync/atomic"
)

type workerSupervisorState int

const (
	workerSupervisorIdle workerSupervisorState = iota
	workerSupervisorRunning
)

type WorkerSupervisor struct {
	wg sync.WaitGroup

	worker     *Worker
	scheduler  *Scheduler
	dispatcher *Dispatcher

	shutdown context.CancelFunc

	state atomic.Value

	logger Logger
}

func (ws *WorkerSupervisor) WorkerStatus() workerState { return ws.worker.Status() }

func (ws *WorkerSupervisor) SchedulerStatus() schedulerState { return ws.scheduler.Status() }

func (ws *WorkerSupervisor) PauseScheduler() { ws.scheduler.Pause() }

func (ws *WorkerSupervisor) ResumeScheduler() { ws.scheduler.Resume() }

func (ws *WorkerSupervisor) Shutdown() {
	if ws.shutdown != nil {
		ws.shutdown()
		ws.shutdown = nil
		ws.dispatcher.Close()

		ws.wg.Wait()
		ws.state.Store(workerSupervisorIdle)
	}
}

func (ws *WorkerSupervisor) Start(ctx context.Context) {
	if ws.state.Load().(workerSupervisorState) != workerSupervisorIdle {
		ws.logger.Error("WorkerSupervisor is already running, cannot start again")
		return
	}

	ctx, cancel := context.WithCancel(ctx)
	ws.shutdown = cancel

	ws.wg.Add(1)
	go func() { defer ws.wg.Done(); ws.worker.Run(ctx) }()

	ws.wg.Add(1)
	go func() { defer ws.wg.Done(); ws.scheduler.Run(ctx) }()
}

func newWorkerSupervisor(
	worker *Worker, scheduler *Scheduler, dispatcher *Dispatcher, logger Logger,
) *WorkerSupervisor {
	ws := &WorkerSupervisor{
		logger:     logger,
		worker:     worker,
		scheduler:  scheduler,
		dispatcher: dispatcher,
	}
	ws.state.Store(workerSupervisorIdle)
	return ws
}
