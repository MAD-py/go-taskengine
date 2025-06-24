package store

import "time"

type TaskState struct {
	Name string `json:"name"`

	LastTick  time.Time `json:"last_tick"`
	Iteration uint      `json:"iteration"`
}

type TaskExecution struct {
	Name string `json:"name"`

	StartedAt  time.Time     `json:"started_at"`
	FinishedAt time.Time     `json:"finished_at"`
	Duration   time.Duration `json:"duration"`

	Success  bool   `json:"success"`
	ErrorMsg string `json:"error_msg,omitempty"`
}
