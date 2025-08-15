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
		name       string
		expr       string
		runOnStart bool
		expected   string
	}{
		{
			name:       "weekday morning without runOnStart",
			expr:       "0 9 * * 1-5",
			runOnStart: false,
			expected:   "Cron(expr=0 9 * * 1-5, runOnStart=false)",
		},
		{
			name:       "every minute with runOnStart",
			expr:       "* * * * *",
			runOnStart: true,
			expected:   "Cron(expr=* * * * *, runOnStart=true)",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			trigger := &cronTrigger{expr: tc.expr, runOnStart: tc.runOnStart}

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
		runOnStart  bool
		expectError bool
	}{
		{
			name:        "every minute without runOnStart",
			expr:        "* * * * *",
			runOnStart:  false,
			expectError: false,
		},
		{
			name:        "weekdays 9 AM with runOnStart",
			expr:        "0 9 * * 1-5",
			runOnStart:  true,
			expectError: false,
		},
		{
			name:        "first day of month",
			expr:        "0 0 1 * *",
			runOnStart:  false,
			expectError: false,
		},
		{
			name:        "empty expression",
			expr:        "",
			runOnStart:  false,
			expectError: true,
		},
		{
			name:        "invalid expression",
			expr:        "invalid",
			runOnStart:  true,
			expectError: true,
		},
		{
			name:        "too few fields",
			expr:        "* * * *",
			runOnStart:  false,
			expectError: true,
		},
		{
			name:        "invalid minute",
			expr:        "60 * * * *",
			runOnStart:  true,
			expectError: true,
		},
		{
			name:        "invalid hour",
			expr:        "* 25 * * *",
			runOnStart:  false,
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			trigger, err := NewCronTrigger(tc.expr, tc.runOnStart)

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

				// Test that runOnStart is set correctly
				cronTrig, ok := trigger.(*cronTrigger)
				if !ok {
					t.Errorf("expected *cronTrigger, got %T", trigger)
				}
				if cronTrig.runOnStart != tc.runOnStart {
					t.Errorf("expected runOnStart=%v, got %v", tc.runOnStart, cronTrig.runOnStart)
				}
			}
		})
	}
}

func TestCronTriggerNext(t *testing.T) {
	tests := []struct {
		name        string
		expr        string
		runOnStart  bool
		lastRun     time.Time
		expectError bool
	}{
		{
			name:        "every minute with zero lastRun and runOnStart false",
			expr:        "* * * * *",
			runOnStart:  false,
			lastRun:     time.Time{},
			expectError: false,
		},
		{
			name:        "every minute with zero lastRun and runOnStart true",
			expr:        "* * * * *",
			runOnStart:  true,
			lastRun:     time.Time{},
			expectError: false,
		},
		{
			name:        "every minute with lastRun",
			expr:        "* * * * *",
			runOnStart:  true,
			lastRun:     time.Now().Add(-30 * time.Second),
			expectError: false,
		},
		{
			name:        "weekdays 9 AM with runOnStart true",
			expr:        "0 9 * * 1-5",
			runOnStart:  true,
			lastRun:     time.Time{},
			expectError: false,
		},
		{
			name:        "first day of month with runOnStart false",
			expr:        "0 0 1 * *",
			runOnStart:  false,
			lastRun:     time.Now().Add(-24 * time.Hour),
			expectError: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			trigger := &cronTrigger{expr: tc.expr, runOnStart: tc.runOnStart}

			before := time.Now()
			next, err := trigger.Next(tc.lastRun)
			after := time.Now()

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

				// Test runOnStart behavior when lastRun is zero
				if tc.lastRun.IsZero() && tc.runOnStart {
					// Should return time.Now() when runOnStart=true
					if next.Before(before) || next.After(after.Add(time.Second)) {
						t.Errorf("expected next time to be around now with runOnStart=true, got %v", next)
					}
				} else if tc.lastRun.IsZero() && !tc.runOnStart {
					// Should calculate next cron time from now when runOnStart=false
					if next.Before(before) {
						t.Errorf("expected next time to be after now with runOnStart=false, got %v", next)
					}
				} else {
					// For non-zero lastRun, should calculate next cron time after lastRun
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
	trigger := &cronTrigger{expr: "invalid expression", runOnStart: false}

	_, err := trigger.Next(time.Now())

	if err == nil {
		t.Error("expected error for invalid cron expression in NextTickAfter, got nil")
	}
}