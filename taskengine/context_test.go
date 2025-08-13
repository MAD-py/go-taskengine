package taskengine

import (
	"context"
	"testing"
	"time"
)

// mockLogger is a simple logger implementation for testing
type mockLogger struct{}

func (m *mockLogger) Info(msg string)                        {}
func (m *mockLogger) Infof(format string, args ...any)      {}
func (m *mockLogger) Warn(msg string)                        {}
func (m *mockLogger) Warnf(format string, args ...any)      {}
func (m *mockLogger) Error(msg string)                       {}
func (m *mockLogger) Errorf(format string, args ...any)     {}

func TestContextLogger(t *testing.T) {
	logger := &mockLogger{}
	ctx := &Context{logger: logger}

	if got := ctx.Logger(); got != logger {
		t.Errorf("expected logger %v, got %v", logger, got)
	}
}

func TestContextTaskName(t *testing.T) {
	want := "test-task"
	ctx := &Context{taskName: want}

	if got := ctx.TaskName(); got != want {
		t.Errorf("expected task name %s, got %s", want, got)
	}
}

func TestContextLastTick(t *testing.T) {
	want := time.Now().Add(-time.Hour)
	tick := &Tick{lastTick: want}
	ctx := &Context{tick: tick}

	if got := ctx.LastTick(); got != want {
		t.Errorf("expected last tick %v, got %v", want, got)
	}
}

func TestContextCurrentTick(t *testing.T) {
	want := time.Now()
	tick := &Tick{currentTick: want}
	ctx := &Context{tick: tick}

	if got := ctx.CurrentTick(); got != want {
		t.Errorf("expected current tick %v, got %v", want, got)
	}
}

func TestContextDeadline(t *testing.T) {
	deadline := time.Now().Add(time.Hour)
	baseCtx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()

	ctx := &Context{ctx: baseCtx}

	gotDeadline, ok := ctx.Deadline()
	if !ok {
		t.Error("expected deadline to be set")
	}

	if gotDeadline != deadline {
		t.Errorf("expected deadline %v, got %v", deadline, gotDeadline)
	}
}

func TestContextDone(t *testing.T) {
	baseCtx, cancel := context.WithCancel(context.Background())
	ctx := &Context{ctx: baseCtx}

	// Should not be done initially
	select {
	case <-ctx.Done():
		t.Error("context should not be done initially")
	default:
		// Expected
	}

	// Cancel and check it's done
	cancel()
	select {
	case <-ctx.Done():
		// Expected
	default:
		t.Error("context should be done after cancel")
	}
}

func TestContextErr(t *testing.T) {
	baseCtx, cancel := context.WithCancel(context.Background())
	ctx := &Context{ctx: baseCtx}

	// Should have no error initially
	if err := ctx.Err(); err != nil {
		t.Errorf("expected no error initially, got %v", err)
	}

	// Cancel and check error
	cancel()
	if err := ctx.Err(); err != context.Canceled {
		t.Errorf("expected context.Canceled error, got %v", err)
	}
}

func TestContextValue(t *testing.T) {
	baseCtx := context.WithValue(context.Background(), "test-key", "test-value")
	ctx := &Context{ctx: baseCtx}

	// Test custom ContextKey returns self
	if got := ctx.Value(ContextKey); got != ctx {
		t.Errorf("expected context itself for ContextKey, got %v", got)
	}

	// Test value from base context
	want := "test-value"
	if got := ctx.Value("test-key"); got != want {
		t.Errorf("expected value %s, got %v", want, got)
	}

	// Test non-existent key
	if got := ctx.Value("non-existent"); got != nil {
		t.Errorf("expected nil for non-existent key, got %v", got)
	}
}

func TestContextValueWithNilBaseContext(t *testing.T) {
	ctx := &Context{ctx: nil}

	// Test custom ContextKey returns self even with nil base context
	if got := ctx.Value(ContextKey); got != ctx {
		t.Errorf("expected context itself for ContextKey, got %v", got)
	}

	// Test other keys return nil when base context is nil
	if got := ctx.Value("any-key"); got != nil {
		t.Errorf("expected nil for any key when base context is nil, got %v", got)
	}
}