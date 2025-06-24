package filestore

import (
	"os"

	"github.com/MAD-py/go-taskengine/taskengine/store"
)

type fileStore struct {
	states     *dataStore // data: map[string]*store.TaskState
	executions *dataStore // data: map[string][]*store.TaskExecution
}

// ====================================
// ========== [ Task State ] ==========
// ====================================

func (fs *fileStore) GetState(name string) (*store.TaskState, bool) {
	if state, exists := fs.states.Load(name); exists {
		return state.(*store.TaskState), true
	}
	return nil, false
}

func (fs *fileStore) SaveState(state *store.TaskState) error {
	fs.states.Store(state.Name, state)
	return fs.states.saveToFile()
}

func (fs *fileStore) DeleteState(name string) error {
	fs.states.Delete(name)
	return fs.states.saveToFile()
}

// ====================================
// ======== [ Task Execution ] ========
// ====================================

func (fs *fileStore) GetLastExecution(name string) (*store.TaskExecution, bool) {
	if executions, exists := fs.executions.Load(name); exists {
		executions := executions.([]*store.TaskExecution)
		if len(executions) > 0 {
			return executions[len(executions)-1], true
		}
	}
	return nil, false
}

func (fs *fileStore) GetExecutions(name string) ([]*store.TaskExecution, bool) {
	if executions, exists := fs.executions.Load(name); exists {
		return executions.([]*store.TaskExecution), true
	}
	return nil, false
}

func (fs *fileStore) SaveExecution(execution *store.TaskExecution) error {
	fs.executions.Store(execution.Name, execution)
	return fs.executions.saveToFile()
}

func (fs *fileStore) DeleteExecution(name string) error {
	fs.executions.Delete(name)
	return fs.executions.saveToFile()
}

func New(folderPath string) (*fileStore, error) {
	if err := os.MkdirAll(folderPath, 0755); err != nil {
		return nil, err
	}

	states, err := newDataStore(folderPath + "/states.json")
	if err != nil {
		return nil, err
	}

	executions, err := newDataStore(folderPath + "/executions.json")
	if err != nil {
		return nil, err
	}

	return &fileStore{
		states:     states,
		executions: executions,
	}, nil
}
