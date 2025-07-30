package postgresql

import (
	"github.com/MAD-py/go-taskengine/taskengine/store"
)

type taskStore struct {
	db DB
}

func (ts *taskStore) createStore() error {
	query := `
		CREATE TABLE IF NOT EXISTS tasks (
			id          SERIAL     PRIMARY KEY,
			name        TEXT       NOT NULL UNIQUE,
			job         TEXT       NOT NULL,
			trigger     TEXT       NOT NULL,
			policy      TEXT       NOT NULL,
			status		TEXT       NOT NULL DEFAULT 'idle',
			iteration   INT        NOT NULL DEFAULT 0,
			created_at  TIMESTAMP  NOT NULL DEFAULT NOW()
		);
	`
	_, err := ts.db.Exec(query)
	return err
}

func (ts *taskStore) deleteStore() error {
	query := "DROP TABLE IF EXISTS tasks;"
	_, err := ts.db.Exec(query)
	return err
}

func (ts *taskStore) clearStore() error {
	query := "TRUNCATE TABLE tasks RESTART IDENTITY;"
	_, err := ts.db.Exec(query)
	return err
}

func (ts *taskStore) exists(name string) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM tasks WHERE name = $1);"
	var exists bool
	err := ts.db.QueryRow(query, name).Scan(&exists)
	return exists, err
}

func (ts *taskStore) save(
	name string, settings *store.TaskSettings,
) error {
	query := `
		INSERT INTO tasks (name, job, trigger, policy)
		VALUES ($1, $2, $3, $4)
		RETURNING id;
	`
	return ts.db.QueryRow(
		query, name,
		settings.Job,
		settings.Trigger,
		settings.Policy,
	).Err()
}

func (ts *taskStore) getSettings(name string) (*store.TaskSettings, error) {
	query := `
		SELECT job, trigger, policy
		FROM tasks
		WHERE name = $1;
	`

	var settings store.TaskSettings
	err := ts.db.
		QueryRow(query, name).
		Scan(&settings.Job, &settings.Trigger, &settings.Policy)
	if err != nil {
		return nil, err
	}
	return &settings, nil
}

func (ts *taskStore) updateStatus(
	name string, status store.TaskStatus,
) error {
	query := "UPDATE tasks SET status = $2 WHERE name = $1;"
	_, err := ts.db.Exec(query, name, status)
	return err
}

func (ts *taskStore) increaseIteration(name string) (int, int, error) {
	query := `
		UPDATE tasks
		SET iteration = iteration + 1
		WHERE name = $1
		RETURNING id, iteration;
	`
	var id int
	var iteration int
	err := ts.db.QueryRow(query, name).Scan(&id, &iteration)
	if err != nil {
		return 0, 0, err
	}
	return id, iteration, nil
}

func newTaskStore(db DB) *taskStore {
	return &taskStore{db: db}
}
