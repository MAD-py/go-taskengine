package store

import "time"

type ExecutionStatus int

func (es ExecutionStatus) String() string {
	switch es {
	case ExecutionStatusFailed:
		return "failed"
	case ExecutionStatusSuccess:
		return "success"
	case ExecutionStatusSkipped:
		return "skipped"
	default:
		return "unknown"
	}
}

const (
	ExecutionStatusFailed ExecutionStatus = iota
	ExecutionStatusSuccess
	ExecutionStatusSkipped
)

type Task struct {
	Name      string    `json:"name"`
	Trigger   string    `json:"trigger"`
	Policy    string    `json:"policy"`
	CreatedAt time.Time `json:"created_at"`
}

type TaskState struct {
	Name      string `json:"name"`
	Iteration int    `json:"iteration"`

	LastExecution time.Time `json:"last_execution"`

	LastStatus   ExecutionStatus `json:"last_status"`
	LastErrorMsg string          `json:"last_error_msg,omitempty"`
}

type TaskExecution struct {
	Name      string `json:"name"`
	Iteration int    `json:"iteration"`

	StartTime time.Time     `json:"start_time"`
	EndTime   time.Time     `json:"end_time"`
	Duration  time.Duration `json:"duration"`

	Status   ExecutionStatus `json:"status"`
	ErrorMsg string          `json:"error_msg,omitempty"`
}
