package taskengine

import (
	"errors"
	"time"
)

type Tick struct {
	lastTick    time.Time
	currentTick time.Time
}

type Dispatcher struct {
	queue chan *Tick
}

func (d *Dispatcher) Capacity() int {
	return cap(d.queue)
}

func (d *Dispatcher) Size() int {
	return len(d.queue)
}

func (d *Dispatcher) Enqueue(tick *Tick) error {
	select {
	case d.queue <- tick:
		return nil
	default:
		return errors.New("dispatcher queue is full")
	}
}

func (d *Dispatcher) Dequeue() <-chan *Tick {
	return d.queue
}

func (d *Dispatcher) Close() {
	close(d.queue)
}

func newDispatcher(capacity int) *Dispatcher {
	if capacity <= 0 {
		capacity = 100 // Default capacity
	}
	return &Dispatcher{
		queue: make(chan *Tick, capacity),
	}
}
