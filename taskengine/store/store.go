package store

import "context"

type TaskStore interface {
	CreateTaskStore(ctx context.Context) error
	DeleteTaskStore(ctx context.Context) error
	ClearTaskStore(ctx context.Context) error

	GetTask(ctx context.Context, name string) (*Task, error)
	SaveTask(ctx context.Context, task *Task) error
	RemoveTask(ctx context.Context, name string) error
}

type StateStore interface {
	CreateStateStore(ctx context.Context) error
	DeleteStateStore(ctx context.Context) error
	ClearStateStore(ctx context.Context) error

	GetState(ctx context.Context, name string) (*State, error)
	SaveState(ctx context.Context, state *State) error
	RemoveState(ctx context.Context, name string) error
}

type ExecutionStore interface {
	CreateExecutionStore(ctx context.Context) error
	DeleteExecutionStore(ctx context.Context) error
	ClearExecutionStore(ctx context.Context) error

	SaveExecution(ctx context.Context, execution *Execution) error
}

type Store interface {
	TaskStore
	StateStore
	ExecutionStore

	CreateStores(ctx context.Context) error
	DeleteStores(ctx context.Context) error
	ClearStores(ctx context.Context) error
}
