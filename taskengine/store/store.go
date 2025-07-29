package store

import "context"

type Store interface {
	CreateStores(ctx context.Context) error
	DeleteStores(ctx context.Context) error
	ClearStores(ctx context.Context) error

	SaveTask(ctx context.Context, name string, settings *TaskSettings) error
	TaskExists(ctx context.Context, name string) (bool, error)
	SaveExecution(ctx context.Context, name string, info *ExecutionInfo) error
	GetTaskSettings(ctx context.Context, name string) (*TaskSettings, error)
	UpdateTaskStatus(ctx context.Context, name string, status TaskStatus) error
}
