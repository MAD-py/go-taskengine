package postgresql

import (
	"database/sql"

	"github.com/MAD-py/go-taskengine/taskengine/store"
)

type DB interface {
	Exec(query string, args ...any) (sql.Result, error)
	QueryRow(query string, args ...any) *sql.Row
}

var _ store.Store = (*PostgresStore)(nil)

type PostgresStore struct {
	taskStore      *taskStore
	executionStore *executionStore
}

func (ps *PostgresStore) CreateStores() error {
	if err := ps.taskStore.createStore(); err != nil {
		return err
	}
	if err := ps.executionStore.createStore(); err != nil {
		return err
	}
	return nil
}

func (ps *PostgresStore) DeleteStores() error {
	if err := ps.executionStore.deleteStore(); err != nil {
		return err
	}
	if err := ps.taskStore.deleteStore(); err != nil {
		return err
	}
	return nil
}

func (ps *PostgresStore) ClearStores() error {
	if err := ps.executionStore.clearStore(); err != nil {
		return err
	}
	if err := ps.taskStore.clearStore(); err != nil {
		return err
	}
	return nil
}

func (ps *PostgresStore) TaskExists(name string) (bool, error) {
	return ps.taskStore.exists(name)
}

func (ps *PostgresStore) SaveTask(name string, settings *store.TaskSettings) error {
	return ps.taskStore.save(name, settings)
}

func (ps *PostgresStore) GetTaskSettings(name string) (*store.TaskSettings, error) {
	return ps.taskStore.getSettings(name)
}

func (ps *PostgresStore) UpdateTaskStatus(name string, status store.TaskStatus) error {
	return ps.taskStore.updateStatus(name, status)
}

func (ps *PostgresStore) SaveExecution(name string, info *store.ExecutionInfo) error {
	taskID, iteration, err := ps.taskStore.increaseIteration(name)
	if err != nil {
		return err
	}

	execution := &store.Execution{
		ExecutionInfo: info,
		TaskID:        taskID,
		Iteration:     iteration,
	}

	return ps.executionStore.save(execution)
}

func NewStore(db DB) *PostgresStore {
	return &PostgresStore{
		taskStore:      newTaskStore(db),
		executionStore: newExecutionStore(db),
	}
}
