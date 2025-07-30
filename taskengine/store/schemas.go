package store

import "time"

type TaskStatus string

const (
	TaskStatusIdle    TaskStatus = "idle"
	TaskStatusRunning TaskStatus = "running"
)

type ExecutionStatus string

const (
	ExecutionStatusPanic   ExecutionStatus = "panic"
	ExecutionStatusError   ExecutionStatus = "error"
	ExecutionStatusSuccess ExecutionStatus = "success"
	ExecutionStatusSkipped ExecutionStatus = "skipped"
)

type TaskSettings struct {
	Job     string `json:"job"`
	Policy  string `json:"policy"`
	Trigger string `json:"trigger"`
}

type ExecutionInfo struct {
	StartTime time.Time       `json:"start_time"`
	EndTime   time.Time       `json:"end_time"`
	Duration  time.Duration   `json:"duration"`
	Status    ExecutionStatus `json:"status"`
	ErrorMsg  string          `json:"error_msg,omitempty"`
}

type Execution struct {
	*ExecutionInfo

	TaskID    int `json:"task_id"`
	Iteration int `json:"iteration"`
}
