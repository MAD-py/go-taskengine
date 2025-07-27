package postgresql

import (
	"context"
	"database/sql"

	"github.com/MAD-py/go-taskengine/taskengine/store"
)

type stateStore struct {
	db DB
}

func (ss *stateStore) createStore(ctx context.Context) error {
	query := `
		CREATE TABLE IF NOT EXISTS task_states (
			task_id         INT 	   PRIMARY KEY REFERENCES tasks(id) ON DELETE CASCADE,
			iteration       INT        NOT NULL,
			last_execution  TIMESTAMP  NOT NULL,
			last_status     TEXT       NOT NULL,
			last_error_msg  TEXT
		);
	`

	_, err := ss.db.ExecContext(ctx, query)
	return err
}

func (ss *stateStore) deleteStore(ctx context.Context) error {
	query := "DROP TABLE IF EXISTS task_states;"
	_, err := ss.db.ExecContext(ctx, query)
	return err
}

func (ss *stateStore) clearStore(ctx context.Context) error {
	query := "TRUNCATE TABLE task_states RESTART IDENTITY;"
	_, err := ss.db.ExecContext(ctx, query)
	return err
}

func (ss *stateStore) save(
	ctx context.Context, taskID int, state *store.State,
) error {
	query := `
		INSERT INTO task_states (task_id, iteration, last_execution, last_status, last_error_msg)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (task_id) DO UPDATE SET
			iteration = EXCLUDED.iteration,
			last_execution = EXCLUDED.last_execution,
			last_status = EXCLUDED.last_status,
			last_error_msg = EXCLUDED.last_error_msg;
	`

	_, err := ss.db.ExecContext(ctx, query,
		taskID,
		state.Iteration,
		state.LastExecution,
		state.LastStatus.String(),
		state.LastErrorMsg,
	)
	return err
}

func (ss *stateStore) get(
	ctx context.Context, name string,
) (*store.State, error) {
	query := `
		SELECT t.name, iteration, last_execution, last_status, last_error_msg
		FROM task_states
		JOIN tasks t ON task_states.task_id = tasks.id 
		WHERE t.name = $1;
	`

	var state store.State
	err := ss.db.QueryRowContext(ctx, query, name).Scan(
		&state.Name,
		&state.Iteration,
		&state.LastExecution,
		&state.LastStatus,
		&state.LastErrorMsg,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &state, nil
}

func (ss *stateStore) remove(ctx context.Context, name string) error {
	query := `
		DELETE FROM task_states
		JOIN tasks t ON task_states.task_id = t.id
		WHERE t.name = $1;
	`

	_, err := ss.db.ExecContext(ctx, query, name)
	return err
}

func newStateStore(db DB) *stateStore {
	return &stateStore{db: db}
}
