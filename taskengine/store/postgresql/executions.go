package postgresql

import (
	"context"

	"github.com/MAD-py/go-taskengine/taskengine/store"
)

type executionStore struct {
	db DB
}

func (es *executionStore) createStore(ctx context.Context) error {
	query := `
		CREATE TABLE IF NOT EXISTS task_executions (
			id          INT        SERIAL PRIMARY KEY,
			task_id     INT        NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
			iteration   INT        NOT NULL,
			start_time  TIMESTAMP  NOT NULL,
			end_time    TIMESTAMP  NOT NULL,
			duration    BIGINT     NOT NULL,
			status      INT        NOT NULL,
			error_msg   TEXT
		);
	`

	_, err := es.db.ExecContext(ctx, query)
	return err
}

func (es *executionStore) deleteStore(ctx context.Context) error {
	_, err := es.db.ExecContext(ctx, "DROP TABLE IF EXISTS task_executions;")
	return err
}

func (es *executionStore) clearStore(ctx context.Context) error {
	_, err := es.db.ExecContext(
		ctx, "TRUNCATE TABLE task_executions RESTART IDENTITY;",
	)
	return err
}

func (es *executionStore) save(
	ctx context.Context, taskID int, execution *store.TaskExecution,
) error {
	query := `
		INSERT INTO task_executions
		(task_id, iteration, start_time, end_time, duration, status, error_msg)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := es.db.ExecContext(ctx, query,
		taskID,
		execution.Iteration,
		execution.StartTime,
		execution.EndTime,
		execution.Duration.Milliseconds(),
		execution.Status.String(),
		execution.ErrorMsg,
	)
	return err
}

func newExecutionStore(db DB) *executionStore {
	return &executionStore{db: db}
}
