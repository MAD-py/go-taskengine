package taskengine

import (
	"errors"
	"fmt"
	"time"

	"github.com/adhocore/gronx"
)

type Trigger interface {
	Next(lastRun time.Time) (time.Time, error)
}

type intervalTrigger struct {
	interval time.Duration

	runOnStart bool
}

func (t *intervalTrigger) String() string {
	return fmt.Sprintf(
		"Interval(interval=%s, runOnStart=%v)",
		t.interval,
		t.runOnStart,
	)
}

func (t *intervalTrigger) Next(lastRun time.Time) (time.Time, error) {
	if lastRun.IsZero() {
		if t.runOnStart {
			return time.Now(), nil
		}
		return time.Now().Add(t.interval), nil
	}
	return lastRun.Add(t.interval), nil
}

func NewIntervalTrigger(interval time.Duration, runOnStart bool) (Trigger, error) {
	if interval <= 0 {
		return nil, errors.New("interval must be a positive duration")
	}
	return &intervalTrigger{interval: interval, runOnStart: runOnStart}, nil
}

type cronTrigger struct {
	expr string
}

func (t *cronTrigger) String() string {
	return fmt.Sprintf("Cron(expr=%s)", t.expr)
}

func (t *cronTrigger) Next(lastRun time.Time) (time.Time, error) {
	if lastRun.IsZero() {
		lastRun = time.Now()
	}

	next, err := gronx.NextTickAfter(t.expr, lastRun, false)
	if err != nil {
		return time.Time{}, err
	}
	return next, nil
}

func NewCronTrigger(expr string) (Trigger, error) {
	if !gronx.IsValid(expr) {
		return nil, errors.New("invalid cron expression")
	}
	return &cronTrigger{expr: expr}, nil
}
