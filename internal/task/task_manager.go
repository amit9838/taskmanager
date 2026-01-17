package task

import (
	"fmt"
	"strings"
	"time"
)

// Repository interface for storage abstraction
type Repository interface {
	Load() ([]Task, error)
	Save(tasks []Task) error
}

type TaskManager struct {
	repo Repository
}

func NewTaskManager(repo Repository) (*TaskManager, error) {
	// NewTaskManager initializes a new TaskManager with the provided Repository.
	// It loads the existing tasks from the repository, and if the file doesn't exist,
	// it initializes the TaskManager with an empty task list.
	// If there is an issue loading the tasks, an error will be returned.
	tasks, err := repo.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load tasks: %w", err)
	}

	// Initialize with empty tasks if file doesn't exist
	if tasks == nil {
		tasks = []Task{}
	}

	return &TaskManager{repo: repo}, nil
}

// -------------------
func (tm *TaskManager) Add(description string) (int, error) {
	description = strings.TrimSpace(description)
	if description == "" {
		return 0, fmt.Errorf("description cannot be empty")
	}

	tasks, err := tm.repo.Load()
	if err != nil {
		return 0, err
	}

	// Generate new ID
	maxID := 0
	for _, t := range tasks {
		if t.ID > maxID {
			maxID = t.ID
		}
	}

	newTask := Task{
		ID:          maxID + 1,
		Description: description,
		Done:        false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	tasks = append(tasks, newTask)

	if err := tm.repo.Save(tasks); err != nil {
		return 0, err
	}

	return newTask.ID, nil
}

func (tm *TaskManager) List() ([]Task, error) {
	return tm.repo.Load()
}

func (tm *TaskManager) MarkDone(id int) error {
	tasks, err := tm.repo.Load()
	if err != nil {
		return err
	}

	found := false
	for i := range tasks {
		if tasks[i].ID == id {
			tasks[i].Done = true
			tasks[i].UpdatedAt = time.Now()
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("task with ID %d not found", id)
	}

	return tm.repo.Save(tasks)
}

func (tm *TaskManager) Delete(id int) error {
	tasks, err := tm.repo.Load()
	if err != nil {
		return err
	}

	found := false
	for i, t := range tasks {
		if t.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("task with ID %d not found", id)
	}

	return tm.repo.Save(tasks)
}

func (tm *TaskManager) Search(query string) ([]Task, error) {
	tasks, err := tm.repo.Load()
	if err != nil {
		return nil, err
	}

	query = strings.ToLower(query)
	var found []Task

	for _, t := range tasks {
		if strings.Contains(strings.ToLower(t.Description), query) {
			found = append(found, t)
		}
	}

	return found, nil
}
