package taskengine

import (
	"context"
	"sync/atomic"
	"time"
)

type schedulerControlCommand int

const (
	schedulerPause schedulerControlCommand = iota
	schedulerResume
)

type schedulerState int

const (
	schedulerIdle schedulerState = iota
	schedulerPaused
	schedulerRunning
)

type Scheduler struct {
	trigger Trigger

	dispatcher *Dispatcher

	control chan schedulerControlCommand

	state          atomic.Value
	catchUpEnabled bool

	logger Logger
}

func (s *Scheduler) Status() schedulerState { return s.state.Load().(schedulerState) }

func (s *Scheduler) Pause() { s.control <- schedulerPause }

func (s *Scheduler) Resume() { s.control <- schedulerResume }

func (s *Scheduler) Run(ctx context.Context) error {
	if s.state.Load().(schedulerState) != schedulerIdle {
		s.logger.Error("Scheduler is already running or paused, cannot start again")
		return nil
	}

	s.logger.Info("Starting Scheduler...")

	var lastTick time.Time
	s.state.Store(schedulerRunning)

Run:
	for {
		now := time.Now()
		nextTick, err := s.trigger.Next(lastTick)
		if err != nil {
			s.logger.Errorf("Error getting next tick: %v", err)
			return err
		}

		if nextTick.Before(now) && !s.catchUpEnabled {
			s.logger.Warnf("Next tick %s is in the past, skipping", nextTick.Format("2006-01-02 15:04:05"))
			lastTick = nextTick
			continue
		}

		select {
		case <-time.After(time.Until(nextTick)):
			tick := Tick{
				lastTick:    lastTick,
				currentTick: nextTick,
			}

			s.logger.Infof("Dispatching tick at %s", nextTick.Format("2006-01-02 15:04:05"))
			err := s.dispatcher.Enqueue(&tick)
			if err != nil {
				s.logger.Errorf("Error dispatching tick: %v", err)
				return err
			}

			lastTick = nextTick
		case cmd := <-s.control:
			switch cmd {
			case schedulerPause:
				s.logger.Info("Scheduler paused")
				s.state.Store(schedulerPaused)
				break Run
			default:
			}
		case <-ctx.Done():
			s.logger.Info("Scheduler shutdown complete")
			return nil
		}
	}

	for {
		select {
		case cmd := <-s.control:
			switch cmd {
			case schedulerResume:
				s.logger.Info("Scheduler resumed")
				s.state.Store(schedulerRunning)
				goto Run
			default:
			}
		case <-ctx.Done():
			s.logger.Info("Scheduler shutdown complete")
			return nil
		}
	}
}

func newScheduler(
	trigger Trigger,
	dispatcher *Dispatcher,
	catchUpEnabled bool,
	logger Logger,
) *Scheduler {
	s := &Scheduler{
		logger:         logger,
		trigger:        trigger,
		control:        make(chan schedulerControlCommand, 1),
		dispatcher:     dispatcher,
		catchUpEnabled: catchUpEnabled,
	}
	s.state.Store(schedulerIdle)
	return s
}
