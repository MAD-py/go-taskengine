package postgresql

import "github.com/MAD-py/go-taskengine/taskengine/store"

type executionStore struct {
	db DB
}

func (es *executionStore) createStore() error {
	query := `
		CREATE TABLE IF NOT EXISTS executions (
			id          SERIAL     PRIMARY KEY,
			task_id     INT        NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
			iteration   INT        NOT NULL,
			start_time  TIMESTAMP  NOT NULL,
			end_time    TIMESTAMP  NOT NULL,
			duration    BIGINT     NOT NULL,
			status      TEXT       NOT NULL,
			error_msg   TEXT
		);
	`

	_, err := es.db.Exec(query)
	return err
}

func (es *executionStore) deleteStore() error {
	query := "DROP TABLE IF EXISTS executions;"
	_, err := es.db.Exec(query)
	return err
}

func (es *executionStore) clearStore() error {
	query := "TRUNCATE TABLE executions RESTART IDENTITY;"
	_, err := es.db.Exec(query)
	return err
}

func (es *executionStore) save(execution *store.Execution) error {
	query := `
		INSERT INTO executions (task_id, iteration, start_time, end_time, duration, status, error_msg)
		VALUES ($1, $2, $3, $4, $5, $6, $7);
	`

	var errorMsg any
	if execution.ErrorMsg != "" {
		errorMsg = execution.ErrorMsg
	}

	_, err := es.db.Exec(
		query,
		execution.TaskID,
		execution.Iteration,
		execution.StartTime,
		execution.EndTime,
		execution.Duration.Milliseconds(),
		execution.Status,
		errorMsg,
	)
	return err
}

func newExecutionStore(db DB) *executionStore {
	return &executionStore{db: db}
}
