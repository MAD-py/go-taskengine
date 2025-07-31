package taskengine

import (
	"testing"
	"time"
)

func TestIntervalTriggerString(t *testing.T) {
	trigger := &intervalTrigger{
		interval:   10 * time.Second,
		runOnStart: true,
	}

	want := "Interval(interval=10s, runOnStart=true)"
	if got := trigger.String(); want != got {
		t.Errorf("expected %s, got %s", want, got)
	}
}
