package store

import "time"

type ExecutionStatus int

const (
	ExecutionStatusFailed ExecutionStatus = iota
	ExecutionStatusSuccess
	ExecutionStatusSkipped
)

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

type TaskStatus int

const (
	TaskStatusIdle TaskStatus = iota
	TaskStatusRegistered
	TaskStatusRunning
)

func (ts TaskStatus) String() string {
	switch ts {
	case TaskStatusIdle:
		return "idle"
	case TaskStatusRegistered:
		return "registered"
	case TaskStatusRunning:
		return "running"
	default:
		return "unknown"
	}
}

type Task struct {
	Name      string    `json:"name"`
	Job       string    `json:"job"`
	Trigger   string    `json:"trigger"`
	Policy    string    `json:"policy"`
	CreatedAt time.Time `json:"created_at"`
}

type State struct {
	Name      string     `json:"name"`
	Status    TaskStatus `json:"status"`
	Iteration int        `json:"iteration"`

	LastExecution time.Time `json:"last_execution"`

	LastStatus   ExecutionStatus `json:"last_status"`
	LastErrorMsg string          `json:"last_error_msg,omitempty"`
}

type Execution struct {
	Name      string `json:"name"`
	Iteration int    `json:"iteration"`

	StartTime time.Time     `json:"start_time"`
	EndTime   time.Time     `json:"end_time"`
	Duration  time.Duration `json:"duration"`

	Status   ExecutionStatus `json:"status"`
	ErrorMsg string          `json:"error_msg,omitempty"`
}
