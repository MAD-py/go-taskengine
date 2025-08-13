package taskengine

import (
	"reflect"
	"testing"
)

func TestCapacity(t *testing.T) {
	want := 10
	dispatcher := &dispatcher{queue: make(chan *Tick, want)}

	if got := dispatcher.Capacity(); want != got {
		t.Errorf("expected capacity %d, got %d", want, got)
	}
}

func TestSize(t *testing.T) {
	want := 0
	dispatcher := &dispatcher{queue: make(chan *Tick, 10)}

	if got := dispatcher.Size(); want != got {
		t.Errorf("expected size %d, got %d", want, got)
	}

	want = 1
	dispatcher.queue <- &Tick{}
	if got := dispatcher.Size(); want != got {
		t.Errorf("expected size %d, got %d", want, got)
	}
}

func TestEnqueue(t *testing.T) {
	dispatcher := &dispatcher{queue: make(chan *Tick, 1)}

	if err := dispatcher.Enqueue(&Tick{}); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if err := dispatcher.Enqueue(&Tick{}); err == nil {
		t.Error("expected error when queue is full, got nil")
	}
}

func TestDequeue(t *testing.T) {
	dispatcher := &dispatcher{queue: make(chan *Tick, 1)}
	tick := &Tick{}

	dispatcher.queue <- &Tick{}
	dequeued := <-dispatcher.Dequeue()
	if !reflect.DeepEqual(dequeued, tick) {
		t.Errorf("expected dequeued tick to be %v, got %v", tick, dequeued)
	}
}

func TestClose(t *testing.T) {
	dispatcher := &dispatcher{queue: make(chan *Tick, 10)}
	dispatcher.Close()

	select {
	case _, ok := <-dispatcher.Dequeue():
		if ok {
			t.Error("expected channel to be closed, but it is still open")
		}
	default:
		// Channel is closed, which is expected
	}
}

func TestNewDispatcher(t *testing.T) {
	want := 100
	dispatcher := newDispatcher(0)

	if got := dispatcher.Capacity(); want != got {
		t.Errorf("expected default capacity %d, got %d", want, got)
	}

	dispatcher = newDispatcher(10)
	if got := dispatcher.Capacity(); 10 != got {
		t.Errorf("expected capacity %d, got %d", 10, got)
	}
}
