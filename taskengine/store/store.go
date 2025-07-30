package store

type Store interface {
	CreateStores() error
	DeleteStores() error
	ClearStores() error

	SaveTask(name string, settings *TaskSettings) error
	TaskExists(name string) (bool, error)
	SaveExecution(name string, info *ExecutionInfo) error
	GetTaskSettings(name string) (*TaskSettings, error)
	UpdateTaskStatus(name string, status TaskStatus) error
}
