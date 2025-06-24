package store

type Store interface {
	GetState(name string) (*TaskState, bool)
	SaveState(state *TaskState) error
	DeleteState(name string) error

	GetExecutions(name string) ([]*TaskExecution, bool)
	SaveExecution(state *TaskExecution) error
	GetLastExecution(name string) (*TaskExecution, bool)
	DeleteExecution(name string) error
}
