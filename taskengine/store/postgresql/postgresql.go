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
	executionStore *executionStore
}

func (ps *PostgresStore) CreateStores(ctx context.Context) error {
	if err := ps.taskStore.createStore(ctx); err != nil {
		return err
	}
	if err := ps.executionStore.createStore(ctx); err != nil {
		return err
	}
	return nil
}

func (ps *PostgresStore) DeleteStores(ctx context.Context) error {
	if err := ps.executionStore.deleteStore(ctx); err != nil {
		return err
	}
	if err := ps.taskStore.deleteStore(ctx); err != nil {
		return err
	}
	return nil
}

func (ps *PostgresStore) ClearStores(ctx context.Context) error {
	if err := ps.executionStore.clearStore(ctx); err != nil {
		return err
	}
	if err := ps.taskStore.clearStore(ctx); err != nil {
		return err
	}
	return nil
}

func (ps *PostgresStore) TaskExists(ctx context.Context, name string) (bool, error) {
	return ps.taskStore.exists(ctx, name)
}

func (ps *PostgresStore) SaveTask(ctx context.Context, name string, settings *store.TaskSettings) error {
	return ps.taskStore.save(ctx, name, settings)
}

func (ps *PostgresStore) GetTaskSettings(ctx context.Context, name string) (*store.TaskSettings, error) {
	return ps.taskStore.getSettings(ctx, name)
}

func (ps *PostgresStore) UpdateTaskStatus(ctx context.Context, name string, status store.TaskStatus) error {
	return ps.taskStore.updateStatus(ctx, name, status)
}

func (ps *PostgresStore) SaveExecution(ctx context.Context, name string, info *store.ExecutionInfo) error {
	taskID, iteration, err := ps.taskStore.increaseIteration(ctx, name)
	if err != nil {
		return err
	}

	execution := &store.Execution{
		ExecutionInfo: info,
		TaskID:        taskID,
		Iteration:     iteration,
	}

	return ps.executionStore.save(ctx, execution)
}

func NewStore(db DB) *PostgresStore {
	return &PostgresStore{
		taskStore:      newTaskStore(db),
		executionStore: newExecutionStore(db),
	}
}
