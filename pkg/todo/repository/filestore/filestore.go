package filestore

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"

	"github.com/snirkop89/todo-cli/pkg/todo"
)

// Filestore is the object representing a file based storage
type FileStore struct {
	filename string
}

func New(filename string) *FileStore {
	return &FileStore{filename}
}

func (store *FileStore) Load() ([]*todo.Todo, error) {
	content, err := os.ReadFile(store.filename)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}

	if len(content) == 0 {
		return nil, nil
	}

	var todos []*todo.Todo
	if err := json.Unmarshal(content, &todos); err != nil {
		return nil, err
	}

	return todos, nil
}

func (store *FileStore) Save(todos []*todo.Todo) error {
	content, err := json.MarshalIndent(todos, "", "  ")
	if err != nil {
		return err
	}
	err = os.WriteFile(store.filename, content, 0644)
	if err != nil {
		return err
	}
	return nil
}
