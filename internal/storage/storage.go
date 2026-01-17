package storage

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/amit9838/taskmanager/internal/task"
)

type JSONStorage struct {
	filename string
}

// NewJSONStorage creates a new JSONStorage instance with the given filename.
// The filename is used to store and load tasks from a local file.
// If the file does not exist, it will be created.
// If the file is empty, an empty task list will be returned.
// The filename is relative to the current working directory when running locally,
// and is relative to the home directory of the current user when building a binary.
func NewJSONStorage(filename string) *JSONStorage {

	// home directory of the current user
	// Use this when building binary
	// home, _ := os.UserHomeDir()
	// path := filepath.Join(home, filename)

	// relative path when running locally
	path := "tasks.json"
	return &JSONStorage{filename: path}
}

// Load reads tasks from the JSON file specified during initialization of the storage.
// If the file does not exist, it will be created with an empty task list.
// If the file is empty, an empty task list will be returned.
// The function returns an error if there was an issue reading or decoding the file.
// The error will contain more information about the issue.
func (s *JSONStorage) Load() ([]task.Task, error) {
	if _, err := os.Stat(s.filename); os.IsNotExist(err) {
		initialData := []byte("[]")
		err := os.WriteFile(s.filename, initialData, 0644)
		if err != nil {
			return nil, fmt.Errorf("could not create file: %w", err)
		}
		return []task.Task{}, nil
	}

	data, err := os.ReadFile(s.filename)

	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	if len(data) == 0 {
		return []task.Task{}, nil
	}

	var tasks []task.Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, fmt.Errorf("failed to decode tasks: %w", err)
	}

	return tasks, nil
}

// Save writes the provided tasks to the JSON file specified during initialization of the storage.
// If an error occurs while encoding or writing the tasks, an error will be returned.
// The error will contain more information about the issue.
func (s *JSONStorage) Save(tasks []task.Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to encode tasks: %w", err)
	}

	// Write to temporary file first
	tempFile := s.filename + ".tmp"
	if err := os.WriteFile(tempFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write temporary file: %w", err)
	}

	// Atomic rename
	if err := os.Rename(tempFile, s.filename); err != nil {
		return fmt.Errorf("failed to rename temporary file: %w", err)
	}

	return nil
}
