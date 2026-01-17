package task

import (
	"errors"
	"testing"
	"time"
)

// MockRepository implements the Repository interface for testing
type MockRepository struct {
	tasks      []Task
	loadError  error
	saveError  error
	loadCalled int
	saveCalled int
	lastSaved  []Task
}

// Load implements Repository interface
func (m *MockRepository) Load() ([]Task, error) {
	m.loadCalled++
	if m.loadError != nil {
		return nil, m.loadError
	}
	return m.tasks, nil
}

// Save implements Repository interface
func (m *MockRepository) Save(tasks []Task) error {
	m.saveCalled++
	m.lastSaved = tasks

	if m.saveError != nil {
		return m.saveError
	}

	m.tasks = tasks
	return nil
}

// TestHelper: Creates a task with specific fields
func createTestTask(id int, description string, done bool) Task {
	now := time.Now()
	return Task{
		ID:          id,
		Description: description,
		Done:        done,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// TestNewTaskManager tests constructor
func TestNewTaskManager(t *testing.T) {
	t.Run("Successfully creates TaskManager", func(t *testing.T) {
		mockRepo := &MockRepository{
			tasks: []Task{createTestTask(1, "Test task", false)},
		}

		tm, err := NewTaskManager(mockRepo)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if tm == nil {
			t.Fatal("Expected TaskManager to be created")
		}

		if mockRepo.loadCalled != 1 {
			t.Errorf("Expected Load to be called once, got %d", mockRepo.loadCalled)
		}
	})

	t.Run("Handles Load error", func(t *testing.T) {
		mockRepo := &MockRepository{
			loadError: errors.New("load failed"),
		}

		tm, err := NewTaskManager(mockRepo)
		if err == nil {
			t.Fatal("Expected error but got none")
		}
		if tm != nil {
			t.Fatal("Expected TaskManager to be nil on error")
		}
		// Note: We don't try to use tm here because it's nil
	})

	t.Run("Initializes with empty tasks when repository returns nil", func(t *testing.T) {
		mockRepo := &MockRepository{
			tasks: nil, // Simulating empty file
		}

		tm, err := NewTaskManager(mockRepo)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if tm == nil {
			t.Fatal("Expected TaskManager to be created")
		}
	})
}

// TestAdd tests the Add method
func TestAdd(t *testing.T) {
	t.Run("Adds task successfully", func(t *testing.T) {
		mockRepo := &MockRepository{}
		tm, err := NewTaskManager(mockRepo)
		if err != nil {
			t.Fatalf("Failed to create TaskManager: %v", err)
		}

		id, err := tm.Add("Buy groceries")
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if id != 1 {
			t.Errorf("Expected ID 1, got %d", id)
		}

		if mockRepo.saveCalled != 1 {
			t.Errorf("Expected Save to be called once, got %d", mockRepo.saveCalled)
		}

		if len(mockRepo.lastSaved) != 1 {
			t.Errorf("Expected 1 task saved, got %d", len(mockRepo.lastSaved))
		}

		savedTask := mockRepo.lastSaved[0]
		if savedTask.ID != 1 || savedTask.Description != "Buy groceries" || savedTask.Done {
			t.Errorf("Task not saved correctly: %+v", savedTask)
		}
	})

	t.Run("Increments ID correctly", func(t *testing.T) {
		mockRepo := &MockRepository{
			tasks: []Task{
				createTestTask(1, "Task 1", false),
				createTestTask(3, "Task 3", false), // Gap in IDs
			},
		}
		tm, err := NewTaskManager(mockRepo)
		if err != nil {
			t.Fatalf("Failed to create TaskManager: %v", err)
		}

		id, err := tm.Add("New Task")
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if id != 4 { // Should be maxID + 1 = 3 + 1 = 4
			t.Errorf("Expected ID 4, got %d", id)
		}
	})

	t.Run("Rejects empty description", func(t *testing.T) {
		mockRepo := &MockRepository{}
		tm, err := NewTaskManager(mockRepo)
		if err != nil {
			t.Fatalf("Failed to create TaskManager: %v", err)
		}

		_, err = tm.Add("   ") // Only whitespace
		if err == nil {
			t.Fatal("Expected error for empty description")
		}

		if mockRepo.saveCalled > 0 {
			t.Error("Save should not be called when description is empty")
		}
	})

	t.Run("Handles Save error", func(t *testing.T) {
		mockRepo := &MockRepository{
			saveError: errors.New("save failed"),
		}
		tm, err := NewTaskManager(mockRepo)
		if err != nil {
			t.Fatalf("Failed to create TaskManager: %v", err)
		}

		_, err = tm.Add("Should fail")
		if err == nil {
			t.Fatal("Expected error but got none")
		}
	})
}

// TestList tests the List method
func TestList(t *testing.T) {
	t.Run("Returns all tasks", func(t *testing.T) {
		expectedTasks := []Task{
			createTestTask(1, "Task 1", false),
			createTestTask(2, "Task 2", true),
		}
		mockRepo := &MockRepository{tasks: expectedTasks}
		tm, err := NewTaskManager(mockRepo)
		if err != nil {
			t.Fatalf("Failed to create TaskManager: %v", err)
		}

		tasks, err := tm.List()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if len(tasks) != len(expectedTasks) {
			t.Errorf("Expected %d tasks, got %d", len(expectedTasks), len(tasks))
		}

		// Load called twice: once in NewTaskManager, once in List
		if mockRepo.loadCalled != 2 {
			t.Errorf("Expected Load to be called twice, got %d", mockRepo.loadCalled)
		}
	})

	t.Run("Handles Load error", func(t *testing.T) {
		mockRepo := &MockRepository{}
		tm, err := NewTaskManager(mockRepo)
		if err != nil {
			t.Fatalf("Failed to create TaskManager: %v", err)
		}

		// Now set the load error for the next Load call
		mockRepo.loadError = errors.New("load failed")

		_, err = tm.List()
		if err == nil {
			t.Fatal("Expected error but got none")
		}
	})
}

// TestMarkDone tests the MarkDone method
func TestMarkDone(t *testing.T) {
	t.Run("Marks task as done", func(t *testing.T) {
		mockRepo := &MockRepository{
			tasks: []Task{
				createTestTask(1, "Task 1", false),
				createTestTask(2, "Task 2", false),
			},
		}
		tm, err := NewTaskManager(mockRepo)
		if err != nil {
			t.Fatalf("Failed to create TaskManager: %v", err)
		}

		err = tm.MarkDone(1)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if mockRepo.saveCalled != 1 {
			t.Errorf("Expected Save to be called once, got %d", mockRepo.saveCalled)
		}

		if !mockRepo.lastSaved[0].Done {
			t.Error("Task 1 should be marked as done")
		}

		if mockRepo.lastSaved[1].Done {
			t.Error("Task 2 should not be marked as done")
		}
	})

	t.Run("Returns error for non-existent task", func(t *testing.T) {
		mockRepo := &MockRepository{
			tasks: []Task{createTestTask(1, "Task 1", false)},
		}
		tm, err := NewTaskManager(mockRepo)
		if err != nil {
			t.Fatalf("Failed to create TaskManager: %v", err)
		}

		err = tm.MarkDone(999)
		if err == nil {
			t.Fatal("Expected error but got none")
		}

		if mockRepo.saveCalled > 0 {
			t.Error("Save should not be called for non-existent task")
		}
	})

	t.Run("Handles Save error", func(t *testing.T) {
		mockRepo := &MockRepository{
			tasks: []Task{createTestTask(1, "Task 1", false)},
		}
		tm, err := NewTaskManager(mockRepo)
		if err != nil {
			t.Fatalf("Failed to create TaskManager: %v", err)
		}

		// Set save error for the next Save call
		mockRepo.saveError = errors.New("save failed")

		err = tm.MarkDone(1)
		if err == nil {
			t.Fatal("Expected error but got none")
		}
	})
}

// TestDelete tests the Delete method
func TestDelete(t *testing.T) {
	t.Run("Deletes task successfully", func(t *testing.T) {
		mockRepo := &MockRepository{
			tasks: []Task{
				createTestTask(1, "Task 1", false),
				createTestTask(2, "Task 2", false),
				createTestTask(3, "Task 3", false),
			},
		}
		tm, err := NewTaskManager(mockRepo)
		if err != nil {
			t.Fatalf("Failed to create TaskManager: %v", err)
		}

		err = tm.Delete(2)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if len(mockRepo.lastSaved) != 2 {
			t.Errorf("Expected 2 tasks after deletion, got %d", len(mockRepo.lastSaved))
		}

		// Check Task 1 still exists
		foundTask1 := false
		for _, task := range mockRepo.lastSaved {
			if task.ID == 1 {
				foundTask1 = true
			}
			if task.ID == 2 {
				t.Error("Task 2 should have been deleted")
			}
		}
		if !foundTask1 {
			t.Error("Task 1 should not have been deleted")
		}
	})

	t.Run("Returns error for non-existent task", func(t *testing.T) {
		mockRepo := &MockRepository{
			tasks: []Task{createTestTask(1, "Task 1", false)},
		}
		tm, err := NewTaskManager(mockRepo)
		if err != nil {
			t.Fatalf("Failed to create TaskManager: %v", err)
		}

		err = tm.Delete(999)
		if err == nil {
			t.Fatal("Expected error but got none")
		}
	})

	t.Run("Handles Save error", func(t *testing.T) {
		mockRepo := &MockRepository{
			tasks: []Task{createTestTask(1, "Task 1", false)},
		}
		tm, err := NewTaskManager(mockRepo)
		if err != nil {
			t.Fatalf("Failed to create TaskManager: %v", err)
		}

		// Set save error for the next Save call
		mockRepo.saveError = errors.New("save failed")

		err = tm.Delete(1)
		if err == nil {
			t.Fatal("Expected error but got none")
		}
	})
}

// TestSearch tests the Search method
func TestSearch(t *testing.T) {
	tasks := []Task{
		createTestTask(1, "Buy grocery", false),
		createTestTask(2, "Clean house", false),
		createTestTask(3, "Go grocery shopping", false),
	}

	t.Run("Finds matching tasks (case-insensitive)", func(t *testing.T) {
		mockRepo := &MockRepository{tasks: tasks}
		tm, err := NewTaskManager(mockRepo)
		if err != nil {
			t.Fatalf("Failed to create TaskManager: %v", err)
		}

		results, err := tm.Search("grocery")
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if len(results) != 2 {
			t.Errorf("Expected 2 results, got %d", len(results))
		}

		// Check both grocery-related tasks are found
		foundTask1 := false
		foundTask3 := false
		for _, task := range results {
			if task.ID == 1 {
				foundTask1 = true
			}
			if task.ID == 3 {
				foundTask3 = true
			}
		}
		if !foundTask1 || !foundTask3 {
			t.Error("Expected both grocery tasks to be found")
		}
	})

	t.Run("Returns empty slice for no matches", func(t *testing.T) {
		mockRepo := &MockRepository{tasks: tasks}
		tm, err := NewTaskManager(mockRepo)
		if err != nil {
			t.Fatalf("Failed to create TaskManager: %v", err)
		}

		results, err := tm.Search("nonexistent")
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if len(results) != 0 {
			t.Errorf("Expected 0 results, got %d", len(results))
		}
	})

	t.Run("Handles Load error", func(t *testing.T) {
		mockRepo := &MockRepository{}
		tm, err := NewTaskManager(mockRepo)
		if err != nil {
			t.Fatalf("Failed to create TaskManager: %v", err)
		}

		// Set load error for the next Load call
		mockRepo.loadError = errors.New("load failed")

		_, err = tm.Search("test")
		if err == nil {
			t.Fatal("Expected error but got none")
		}
	})
}

// Integration-style test showing multiple operations
func TestTaskManagerIntegration(t *testing.T) {
	mockRepo := &MockRepository{}
	tm, err := NewTaskManager(mockRepo)
	if err != nil {
		t.Fatalf("Failed to create TaskManager: %v", err)
	}

	// Add tasks
	id1, err := tm.Add("First task")
	if err != nil {
		t.Fatalf("Failed to add first task: %v", err)
	}

	id2, err := tm.Add("Second task")
	if err != nil {
		t.Fatalf("Failed to add second task: %v", err)
	}

	// List tasks
	tasks, err := tm.List()
	if err != nil {
		t.Fatalf("Failed to list tasks: %v", err)
	}

	if len(tasks) != 2 {
		t.Errorf("Expected 2 tasks, got %d", len(tasks))
	}

	// Mark first task as done
	err = tm.MarkDone(id1)
	if err != nil {
		t.Fatalf("Failed to mark task as done: %v", err)
	}

	// Verify task is done
	tasks, err = tm.List()
	if err != nil {
		t.Fatalf("Failed to list tasks: %v", err)
	}

	for _, task := range tasks {
		if task.ID == id1 && !task.Done {
			t.Error("Task 1 should be marked as done")
		}
		if task.ID == id2 && task.Done {
			t.Error("Task 2 should not be marked as done")
		}
	}

	// Search for tasks
	results, err := tm.Search("second")
	if err != nil {
		t.Fatalf("Failed to search tasks: %v", err)
	}

	if len(results) != 1 || results[0].ID != id2 {
		t.Error("Search should find second task")
	}

	// Delete a task
	err = tm.Delete(id1)
	if err != nil {
		t.Fatalf("Failed to delete task: %v", err)
	}

	tasks, err = tm.List()
	if err != nil {
		t.Fatalf("Failed to list tasks: %v", err)
	}

	if len(tasks) != 1 || tasks[0].ID != id2 {
		t.Error("After deletion, only task 2 should remain")
	}
}
