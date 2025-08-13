package taskengine

import (
	"errors"
	"time"
)

type Tick struct {
	lastTick    time.Time
	currentTick time.Time
}

type Dispatcher interface {
	Size() int
	Close()
	Enqueue(tick *Tick) error
	Dequeue() <-chan *Tick
	Capacity() int
}

type dispatcher struct {
	queue chan *Tick
}

func (d *dispatcher) Capacity() int {
	return cap(d.queue)
}

func (d *dispatcher) Size() int {
	return len(d.queue)
}

func (d *dispatcher) Enqueue(tick *Tick) error {
	select {
	case d.queue <- tick:
		return nil
	default:
		return errors.New("dispatcher queue is full")
	}
}

func (d *dispatcher) Dequeue() <-chan *Tick {
	return d.queue
}

func (d *dispatcher) Close() {
	close(d.queue)
}

func newDispatcher(capacity int) Dispatcher {
	if capacity <= 0 {
		capacity = 100 // Default capacity
	}
	return &dispatcher{
		queue: make(chan *Tick, capacity),
	}
}
