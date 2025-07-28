package postgresql

import (
	"context"
	"database/sql"

	"github.com/MAD-py/go-taskengine/taskengine/store"
)

type taskStore struct {
	db DB
}

func (ts *taskStore) createStore(ctx context.Context) error {
	query := `
		CREATE TABLE IF NOT EXISTS tasks (
			id          SERIAL     PRIMARY KEY,
			name        TEXT       NOT NULL UNIQUE,
			job         TEXT       NOT NULL,
			trigger     TEXT       NOT NULL,
			policy      TEXT       NOT NULL,
			created_at  TIMESTAMP  NOT NULL DEFAULT NOW()
		);
	`
	_, err := ts.db.ExecContext(ctx, query)
	return err
}

func (ts *taskStore) deleteStore(ctx context.Context) error {
	query := "DROP TABLE IF EXISTS tasks;"
	_, err := ts.db.ExecContext(ctx, query)
	return err
}

func (ts *taskStore) clearStore(ctx context.Context) error {
	query := "TRUNCATE TABLE tasks RESTART IDENTITY;"
	_, err := ts.db.ExecContext(ctx, query)
	return err
}

func (ts *taskStore) save(ctx context.Context, task *store.Task) error {
	query := `
		INSERT INTO tasks (name, job, status, trigger, policy)
		VALUES ($1, $2, $3, $4, $5)
	`
	return ts.db.QueryRowContext(
		ctx, query, task.Name, task.Job,
		task.Status, task.Trigger, task.Policy,
	).Err()
}

func (ts *taskStore) getID(ctx context.Context, name string) (int, error) {
	query := "SELECT id FROM tasks WHERE name = $1;"

	var id int
	err := ts.db.QueryRowContext(ctx, query, name).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, store.ErrorTaskNotFound
		}
		return 0, err
	}

	return id, nil
}

func (ts *taskStore) get(ctx context.Context, name string) (*store.Task, error) {
	query := `
		SELECT name, job, status, trigger, policy
		FROM tasks
		WHERE name = $1;
	`

	var task store.Task
	err := ts.db.QueryRowContext(ctx, query, name).Scan(
		&task.Name,
		&task.Job,
		&task.Status,
		&task.Trigger,
		&task.Policy,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrorTaskNotFound
		}
		return nil, err
	}

	return &task, nil
}

func (ts *taskStore) remove(ctx context.Context, name string) error {
	query := "DELETE FROM tasks WHERE name = $1;"
	_, err := ts.db.ExecContext(ctx, query, name)
	return err
}

func newTaskStore(db DB) *taskStore {
	return &taskStore{db: db}
}
