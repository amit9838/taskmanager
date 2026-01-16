package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"
)

const defaultFilename = "tasks.json"

// Domain Layer
// ------------------------------------------------------------

type Task struct {
	ID          int       `json:"id"`
	Description string    `json:"description"`
	Done        bool      `json:"done"`
	CreatedAt   time.Time `json:"created_at"`
}

type TaskManager struct {
	filename string
	tasks    []Task
}

// NewTaskManager initializes the manager and loads existing data
func NewTaskManager(filename string) (*TaskManager, error) {
	tm := &TaskManager{filename: filename}
	if err := tm.load(); err != nil {
		return nil, err
	}
	return tm, nil
}

// Storage Logic
// ------------------------------------------------------------

func (tm *TaskManager) load() error {
	if _, err := os.Stat(tm.filename); os.IsNotExist(err) {
		tm.tasks = []Task{}
		return nil
	}

	data, err := os.ReadFile(tm.filename)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	return json.Unmarshal(data, &tm.tasks)
}

func (tm *TaskManager) save() error {
	data, err := json.MarshalIndent(tm.tasks, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to encode tasks: %w", err)
	}
	// 0644 is standard readable/writable by owner, readable by group/others
	return os.WriteFile(tm.filename, data, 0644)
}

// Business Logic
// ------------------------------------------------------------

func (tm *TaskManager) Add(description string) int {
	// Fix: logic to prevent duplicate IDs if tasks are deleted
	maxID := 0
	for _, t := range tm.tasks {
		if t.ID > maxID {
			maxID = t.ID
		}
	}

	newTask := Task{
		ID:          maxID + 1,
		Description: description,
		Done:        false,
		CreatedAt:   time.Now(),
	}

	tm.tasks = append(tm.tasks, newTask)
	return newTask.ID
}

func (tm *TaskManager) List() {
	if len(tm.tasks) == 0 {
		fmt.Println("No tasks found.")
		return
	}

	// Use Tabwriter for clean columns
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tStatus\tDescription\tCreated")
	fmt.Fprintln(w, "--\t------\t-----------\t-------")

	for _, t := range tm.tasks {
		status := "[ ]"
		if t.Done {
			status = "[x]"
		}
		// Format date nicely
		dateStr := t.CreatedAt.Format("2006-01-02")
		fmt.Fprintf(w, "%d\t%s\t%s\t%s\n", t.ID, status, t.Description, dateStr)
	}
	w.Flush()
}

func (tm *TaskManager) MarkDone(id int) error {
	for i := range tm.tasks {
		if tm.tasks[i].ID == id {
			tm.tasks[i].Done = true
			return nil
		}
	}
	return fmt.Errorf("task with ID %d not found", id)
}

func (tm *TaskManager) Delete(id int) error {
	for i, t := range tm.tasks {
		if t.ID == id {
			// Idiomatic slice removal (order not preserved for efficiency,
			// use standard copy method if order matters)
			tm.tasks = append(tm.tasks[:i], tm.tasks[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("task with ID %d not found", id)
}

func (tm *TaskManager) Search(query string) {
	query = strings.ToLower(query)
	var found []Task
	for _, t := range tm.tasks {
		if strings.Contains(strings.ToLower(t.Description), query) {
			found = append(found, t)
		}
	}

	if len(found) == 0 {
		fmt.Println("No results found.")
		return
	}

	fmt.Printf("Found %d results:\n", len(found))
	for _, t := range found {
		status := "[ ]"
		if t.Done {
			status = "[x]"
		}
		fmt.Printf("%s %d: %s\n", status, t.ID, t.Description)
	}
}

// CLI / Presentation Layer
// ------------------------------------------------------------

func printUsage() {
	fmt.Println("Task Manager CLI")
	fmt.Println("\nUsage:")
	fmt.Println("  add \"<description>\"   Create a new task")
	fmt.Println("  list                  List all tasks")
	fmt.Println("  done <id>             Mark a task as completed")
	fmt.Println("  del <id>              Delete a task")
	fmt.Println("  search \"<term>\"       Search tasks")
	fmt.Println("")
}

func main() {
	// 1. Initialize Manager
	tm, err := NewTaskManager(defaultFilename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing task manager: %v\n", err)
		os.Exit(1)
	}

	// 2. Validate Args
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	// 3. Command Routing using Subcommands
	// Using NewFlagSet allows distinct parsing for each command
	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	doneCmd := flag.NewFlagSet("done", flag.ExitOnError)
	delCmd := flag.NewFlagSet("del", flag.ExitOnError)
	searchCmd := flag.NewFlagSet("search", flag.ExitOnError)

	switch os.Args[1] {
	case "add":
		addCmd.Parse(os.Args[2:])
		if addCmd.NArg() == 0 {
			fmt.Println("Error: Please provide a description.")
			return
		}
		// Join all remaining args so quotes aren't strictly required
		desc := strings.Join(addCmd.Args(), " ")
		id := tm.Add(desc)
		if err := tm.save(); err != nil {
			fmt.Printf("Error saving: %v\n", err)
			return
		}
		fmt.Printf("Task added with ID: %d\n", id)

	case "list":
		tm.List()

	case "done":
		doneCmd.Parse(os.Args[2:])
		if doneCmd.NArg() == 0 {
			fmt.Println("Error: Please provide a task ID.")
			return
		}
		id := parseID(doneCmd.Arg(0))
		if err := tm.MarkDone(id); err != nil {
			fmt.Println(err)
			return
		}
		tm.save()
		fmt.Printf("Task %d marked as done.\n", id)

	case "del":
		delCmd.Parse(os.Args[2:])
		if delCmd.NArg() == 0 {
			fmt.Println("Error: Please provide a task ID.")
			return
		}
		id := parseID(delCmd.Arg(0))
		if err := tm.Delete(id); err != nil {
			fmt.Println(err)
			return
		}
		tm.save()
		fmt.Printf("Task %d deleted.\n", id)

	case "search":
		searchCmd.Parse(os.Args[2:])
		if searchCmd.NArg() == 0 {
			fmt.Println("Error: Please provide a search term.")
			return
		}
		tm.Search(searchCmd.Arg(0))

	default:
		printUsage()
	}
}

// Helper to safely parse IDs
func parseID(s string) int {
	var id int
	_, err := fmt.Sscanf(s, "%d", &id)
	if err != nil {
		return -1
	}
	return id
}
