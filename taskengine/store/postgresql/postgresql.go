package postgresql

import (
	"context"
	"database/sql"

	"github.com/MAD-py/go-taskengine/taskengine/store"
)

type DB interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

var _ store.Store = (*PostgresStore)(nil)

type PostgresStore struct {
	taskStore      *taskStore
	stateStore     *stateStore
	executionStore *executionStore
}

func (ps *PostgresStore) CreateTaskStore(ctx context.Context) error {
	return ps.taskStore.createStore(ctx)
}

func (ps *PostgresStore) CreateStateStore(ctx context.Context) error {
	return ps.stateStore.createStore(ctx)
}

func (ps *PostgresStore) CreateExecutionStore(ctx context.Context) error {
	return ps.executionStore.createStore(ctx)
}

func (ps *PostgresStore) CreateStores(ctx context.Context) error {
	if err := ps.taskStore.createStore(ctx); err != nil {
		return err
	}
	if err := ps.stateStore.createStore(ctx); err != nil {
		return err
	}
	if err := ps.executionStore.createStore(ctx); err != nil {
		return err
	}
	return nil
}

func (ps *PostgresStore) DeleteTaskStore(ctx context.Context) error {
	return ps.taskStore.deleteStore(ctx)
}

func (ps *PostgresStore) DeleteStateStore(ctx context.Context) error {
	return ps.stateStore.deleteStore(ctx)
}

func (ps *PostgresStore) DeleteExecutionStore(ctx context.Context) error {
	return ps.executionStore.deleteStore(ctx)
}

func (ps *PostgresStore) DeleteStores(ctx context.Context) error {
	if err := ps.executionStore.deleteStore(ctx); err != nil {
		return err
	}
	if err := ps.stateStore.deleteStore(ctx); err != nil {
		return err
	}
	if err := ps.taskStore.deleteStore(ctx); err != nil {
		return err
	}
	return nil
}

func (ps *PostgresStore) ClearTaskStore(ctx context.Context) error {
	return ps.taskStore.clearStore(ctx)
}

func (ps *PostgresStore) ClearStateStore(ctx context.Context) error {
	return ps.stateStore.clearStore(ctx)
}

func (ps *PostgresStore) ClearExecutionStore(ctx context.Context) error {
	return ps.executionStore.clearStore(ctx)
}

func (ps *PostgresStore) ClearStores(ctx context.Context) error {
	if err := ps.executionStore.clearStore(ctx); err != nil {
		return err
	}
	if err := ps.stateStore.clearStore(ctx); err != nil {
		return err
	}
	if err := ps.taskStore.clearStore(ctx); err != nil {
		return err
	}
	return nil
}

func (ps *PostgresStore) SaveTask(
	ctx context.Context, task *store.Task,
) error {
	return ps.taskStore.save(ctx, task)
}

func (ps *PostgresStore) SaveState(
	ctx context.Context, state *store.TaskState,
) error {
	taskID, err := ps.taskStore.getID(ctx, state.Name)
	if err != nil {
		return err
	}

	return ps.stateStore.save(ctx, taskID, state)
}

func (ps *PostgresStore) SaveExecution(
	ctx context.Context, execution *store.TaskExecution,
) error {
	taskID, err := ps.taskStore.getID(ctx, execution.Name)
	if err != nil {
		return err
	}

	return ps.executionStore.save(ctx, taskID, execution)
}

func (ps *PostgresStore) GetTask(
	ctx context.Context, name string,
) (*store.Task, error) {
	return ps.taskStore.get(ctx, name)
}

func (ps *PostgresStore) GetState(
	ctx context.Context, name string,
) (*store.TaskState, error) {
	return ps.stateStore.get(ctx, name)
}

func (ps *PostgresStore) RemoveTask(
	ctx context.Context, name string,
) error {
	return ps.taskStore.remove(ctx, name)
}

func (ps *PostgresStore) RemoveState(
	ctx context.Context, name string,
) error {
	return ps.stateStore.remove(ctx, name)
}

func NewStore(db DB) *PostgresStore {
	return &PostgresStore{
		taskStore:      newTaskStore(db),
		stateStore:     newStateStore(db),
		executionStore: newExecutionStore(db),
	}
}
