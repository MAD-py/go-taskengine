package taskengine

import (
	"strings"
	"testing"
	"time"
)

func TestIntervalTriggerString(t *testing.T) {
	tests := []struct {
		name       string
		interval   time.Duration
		runOnStart bool
		expected   string
	}{
		{
			name:       "runOnStart false",
			interval:   10 * time.Second,
			runOnStart: false,
			expected:   "Interval(interval=10s, runOnStart=false)",
		},
		{
			name:       "runOnStart true",
			interval:   5 * time.Minute,
			runOnStart: true,
			expected:   "Interval(interval=5m0s, runOnStart=true)",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			trigger := &intervalTrigger{
				interval:   tc.interval,
				runOnStart: tc.runOnStart,
			}

			if got := trigger.String(); got != tc.expected {
				t.Errorf("expected %s, got %s", tc.expected, got)
			}
		})
	}
}

func TestIntervalTriggerNext(t *testing.T) {
	tests := []struct {
		name       string
		interval   time.Duration
		runOnStart bool
		lastRun    time.Time
	}{
		{
			name:       "zero lastRun with runOnStart true",
			interval:   30 * time.Second,
			runOnStart: true,
			lastRun:    time.Time{},
		},
		{
			name:       "zero lastRun with runOnStart false",
			interval:   30 * time.Second,
			runOnStart: false,
			lastRun:    time.Time{},
		},
		{
			name:       "with lastRun",
			interval:   15 * time.Minute,
			runOnStart: true, // Should not matter when lastRun is not zero
			lastRun:    time.Now().Add(-10 * time.Minute),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			trigger := &intervalTrigger{
				interval:   tc.interval,
				runOnStart: tc.runOnStart,
			}

			before := time.Now()
			next, err := trigger.Next(tc.lastRun)

			if err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			if tc.lastRun.IsZero() {
				if tc.runOnStart {
					// Should be around now
					after := time.Now()
					if next.Before(before) || next.After(after.Add(time.Second)) {
						t.Errorf("expected next time to be around now, got %v", next)
					}
				} else {
					// Should be now + interval
					expectedNext := before.Add(tc.interval)
					if next.Before(expectedNext.Add(-time.Second)) || next.After(expectedNext.Add(time.Second)) {
						t.Errorf("expected next time to be around %v, got %v", expectedNext, next)
					}
				}
			} else {
				// Should be lastRun + interval
				expected := tc.lastRun.Add(tc.interval)
				if next != expected {
					t.Errorf("expected next time to be %v, got %v", expected, next)
				}
			}
		})
	}
}

func TestNewIntervalTrigger(t *testing.T) {
	tests := []struct {
		name        string
		interval    time.Duration
		runOnStart  bool
		expectError bool
	}{
		{
			name:        "valid interval",
			interval:    10 * time.Second,
			runOnStart:  true,
			expectError: false,
		},
		{
			name:        "zero interval",
			interval:    0,
			runOnStart:  true,
			expectError: true,
		},
		{
			name:        "negative interval",
			interval:    -1 * time.Second,
			runOnStart:  false,
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			trigger, err := NewIntervalTrigger(tc.interval, tc.runOnStart)

			if tc.expectError {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				if trigger != nil {
					t.Errorf("expected nil trigger, got %v", trigger)
				}
				if err != nil && !strings.Contains(err.Error(), "interval must be a positive duration") {
					t.Errorf("expected specific error message, got %v", err)
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
				if trigger == nil {
					t.Error("expected trigger to be created, got nil")
				}
				// Test that it implements Trigger interface
				var _ Trigger = trigger
			}
		})
	}
}

func TestCronTriggerString(t *testing.T) {
	tests := []struct {
		name     string
		expr     string
		expected string
	}{
		{
			name:     "weekday morning",
			expr:     "0 9 * * 1-5",
			expected: "Cron(expr=0 9 * * 1-5)",
		},
		{
			name:     "every minute",
			expr:     "* * * * *",
			expected: "Cron(expr=* * * * *)",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			trigger := &cronTrigger{expr: tc.expr}

			if got := trigger.String(); got != tc.expected {
				t.Errorf("expected %s, got %s", tc.expected, got)
			}
		})
	}
}

func TestNewCronTrigger(t *testing.T) {
	tests := []struct {
		name        string
		expr        string
		expectError bool
	}{
		{
			name:        "every minute",
			expr:        "* * * * *",
			expectError: false,
		},
		{
			name:        "weekdays 9 AM",
			expr:        "0 9 * * 1-5",
			expectError: false,
		},
		{
			name:        "first day of month",
			expr:        "0 0 1 * *",
			expectError: false,
		},
		{
			name:        "empty expression",
			expr:        "",
			expectError: true,
		},
		{
			name:        "invalid expression",
			expr:        "invalid",
			expectError: true,
		},
		{
			name:        "too few fields",
			expr:        "* * * *",
			expectError: true,
		},
		{
			name:        "invalid minute",
			expr:        "60 * * * *",
			expectError: true,
		},
		{
			name:        "invalid hour",
			expr:        "* 25 * * *",
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			trigger, err := NewCronTrigger(tc.expr)

			if tc.expectError {
				if err == nil {
					t.Errorf("expected error for expression %s, got nil", tc.expr)
				}
				if trigger != nil {
					t.Errorf("expected nil trigger for expression %s, got %v", tc.expr, trigger)
				}
				if err != nil && !strings.Contains(err.Error(), "invalid cron expression") {
					t.Errorf("expected specific error message for %s, got %v", tc.expr, err)
				}
			} else {
				if err != nil {
					t.Errorf("expected no error for expression %s, got %v", tc.expr, err)
				}
				if trigger == nil {
					t.Errorf("expected trigger to be created for expression %s, got nil", tc.expr)
				}
				// Test that it implements Trigger interface
				var _ Trigger = trigger
			}
		})
	}
}

func TestCronTriggerNext(t *testing.T) {
	tests := []struct {
		name        string
		expr        string
		lastRun     time.Time
		expectError bool
	}{
		{
			name:        "every minute with zero lastRun",
			expr:        "* * * * *",
			lastRun:     time.Time{},
			expectError: false,
		},
		{
			name:        "every minute with lastRun",
			expr:        "* * * * *",
			lastRun:     time.Now().Add(-30 * time.Second),
			expectError: false,
		},
		{
			name:        "weekdays 9 AM",
			expr:        "0 9 * * 1-5",
			lastRun:     time.Time{},
			expectError: false,
		},
		{
			name:        "first day of month",
			expr:        "0 0 1 * *",
			lastRun:     time.Now().Add(-24 * time.Hour),
			expectError: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			trigger := &cronTrigger{expr: tc.expr}

			next, err := trigger.Next(tc.lastRun)

			if tc.expectError {
				if err == nil {
					t.Errorf("expected error for expression %s, got nil", tc.expr)
				}
			} else {
				if err != nil {
					t.Errorf("expected no error for expression %s, got %v", tc.expr, err)
				}

				if next.IsZero() {
					t.Errorf("expected next time to be set, got zero time")
				}

				// For zero lastRun, next should be after now
				if tc.lastRun.IsZero() {
					if next.Before(time.Now()) {
						t.Errorf("expected next time to be after now for zero lastRun, got %v", next)
					}
				} else {
					// For non-zero lastRun, next should be after lastRun
					if next.Before(tc.lastRun) {
						t.Errorf("expected next time to be after lastRun (%v), got %v", tc.lastRun, next)
					}
				}
			}
		})
	}
}

func TestCronTriggerNextError(t *testing.T) {
	// Test error case by directly creating cronTrigger with invalid expression
	// that bypasses the NewCronTrigger validation but fails in NextTickAfter()
	trigger := &cronTrigger{expr: "invalid expression"}

	_, err := trigger.Next(time.Now())

	if err == nil {
		t.Error("expected error for invalid cron expression in NextTickAfter, got nil")
	}
}