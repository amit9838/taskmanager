package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/amit9838/taskmanager/internal/task"
	"github.com/amit9838/taskmanager/pkg/display"
)

type Command interface {
	Execute(manager *task.TaskManager, args []string) error
}

// AddCommand
type AddCommand struct{}

func (c *AddCommand) Execute(manager *task.TaskManager, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("please provide a description")
	}

	desc := strings.Join(args, " ")
	id, err := manager.Add(desc)
	if err != nil {
		return err
	}

	fmt.Printf("Task added with ID: %d\n", id)
	return nil
}

// ListCommand
type ListCommand struct{}

func (c *ListCommand) Execute(manager *task.TaskManager, args []string) error {
	tasks, err := manager.List()
	if err != nil {
		return err
	}

	if len(tasks) == 0 {
		fmt.Println("No tasks found.")
		return nil
	}

	display.PrintTasks(tasks)
	return nil
}

// DoneCommand
type DoneCommand struct{}

func (c *DoneCommand) Execute(manager *task.TaskManager, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("please provide a task ID")
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid task ID: %v", err)
	}

	if err := manager.MarkDone(id); err != nil {
		return err
	}

	fmt.Printf("Task %d marked as done.\n", id)
	return nil
}

// DeleteCommand
type DeleteCommand struct{}

func (c *DeleteCommand) Execute(manager *task.TaskManager, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("please provide a task ID")
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid task ID: %v", err)
	}

	if err := manager.Delete(id); err != nil {
		return err
	}

	fmt.Printf("Task %d deleted.\n", id)
	return nil
}

// SearchCommand
type SearchCommand struct{}

func (c *SearchCommand) Execute(manager *task.TaskManager, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("please provide a search term")
	}

	query := strings.Join(args, " ")
	tasks, err := manager.Search(query)
	if err != nil {
		return err
	}

	if len(tasks) == 0 {
		fmt.Println("No results found.")
		return nil
	}

	fmt.Printf("Found %d results:\n", len(tasks))
	display.PrintTasks(tasks)
	return nil
}

// HelpCommand
type HelpCommand struct{}

func (c *HelpCommand) Execute(manager *task.TaskManager, args []string) error {
	printUsage()
	return nil
}

func printUsage() {
	fmt.Println("Task Manager CLI")
	fmt.Println("\nUsage:")
	fmt.Println("  add \"<description>\"   Create a new task")
	fmt.Println("  list                  List all tasks")
	fmt.Println("  done <id>             Mark a task as completed")
	fmt.Println("  del <id>              Delete a task")
	fmt.Println("  search \"<term>\"       Search tasks")
	fmt.Println("  help                  Show this help message")
	fmt.Println("")
}
