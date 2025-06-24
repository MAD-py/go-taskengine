package filestore

import (
	"encoding/json"
	"os"
	"sync"

	"github.com/MAD-py/go-taskengine/taskengine/store"
)

type dataStore struct {
	mu sync.Mutex

	filePath string

	data *sync.Map
}

func (ds *dataStore) Delete(key string) { ds.data.Delete(key) }

func (ds *dataStore) Load(key string) (any, bool) { return ds.data.Load(key) }

func (ds *dataStore) Store(key string, value any) { ds.data.Store(key, value) }

func (ds *dataStore) toMap() map[string]any {
	data := make(map[string]any)
	ds.data.Range(func(key, value any) bool {
		data[key.(string)] = value
		return true
	})
	return data
}

func (ds *dataStore) saveToFile() error {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	bytes, err := json.MarshalIndent(ds.toMap(), "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(ds.filePath, bytes, 0644)
}

func newDataStore(filePath string) (*dataStore, error) {
	s := &dataStore{
		filePath: filePath,
		data:     &sync.Map{},
	}

	bytes, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return s, nil
		}
		return nil, err
	}

	var data map[string]*store.TaskState
	if err := json.Unmarshal(bytes, &data); err != nil {
		return nil, err
	}

	for key, value := range data {
		s.data.Store(key, value)
	}

	return s, nil
}
