package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Helper functions
// Defining the Task struct
type Task struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}

// Saving the tasks to a file
func SaveTasks(filename string, tasks []Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

// Load the tasks from a file
func LoadTasks(filename string) ([]Task, error) {
	// Check if file exists then return an empty slice
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return []Task{}, nil
	}

	// Read the file
	data, err := os.ReadFile(filename)

	// If error occured during reading the file then return the error
	if err != nil {
		return nil, err
	}
	var task []Task
	err = json.Unmarshal(data, &task)
	return task, err
}

// ------------------------------------------------------------
const filename = "tasks.json"

func main() {
	tasks, err := LoadTasks(filename)
	if err != nil {
		fmt.Printf("Couldn't load file. \nError: %s", err)
		return
	}

	if len(os.Args) < 2 {
		fmt.Println("Usage: task-manager [command] [arguments]")
		return
	}

	command := os.Args[1]

	switch command {

	// -add <description>
	case "-add":
		if len(os.Args) < 3 {
			fmt.Println("Error: Please provide a task description")
			return
		}

		newTask := Task{
			ID:          len(tasks) + 1,
			Description: os.Args[2],
			Done:        false,
		}

		// Append newly added task to existing tasks
		tasks = append(tasks, newTask)

		// save taks to file
		SaveTasks(filename, tasks)
		fmt.Println("Task added")

	// -list
	// Lists all the tasks
	case "-list":
		fmt.Println("Listing all tasks...")
		for _, t := range tasks {
			status := " "
			if t.Done {
				status = "x"
			}
			fmt.Printf("[%s] %d: %s\n", status, t.ID, t.Description)
		}

	// -done <id> <y/n>
	case "-done":
		// fmt.Println("Listing all tasks...")
		if len(os.Args) < 4 {
			fmt.Println("-done requires two arguments <task_id> <y/n>")
			return
		}

		task_id, err := strconv.Atoi(os.Args[2])

		if err != nil {
			fmt.Println("Error: Invalid ID. Please provide a number.")
			return
		}

		response := os.Args[3]

		new_status := false
		if response == "y" || response == "Y" {
			new_status = true
		}

		found := false
		for i := range tasks {
			if tasks[i].ID == task_id {
				tasks[i].Done = new_status
				found = true
				break
			}
		}

		if found {
			SaveTasks(filename, tasks)
			var status string
			if new_status == true {
				status = "Completed"
			} else {
				status = "Pending"
			}
			fmt.Printf("Task %d marked as %s!\n", task_id, status)
		} else {
			fmt.Printf("Error: Task with ID %d not found.\n", task_id)
		}

	// -del <id>
	case "-del":
		var filtred_tasks []Task

		task_id, err := strconv.Atoi(os.Args[2])

		if err != nil {
			fmt.Println("Error: Invalid ID. Please provide a number.")
			return
		}
		found := false
		for _, t := range tasks {
			if task_id != t.ID {
				filtred_tasks = append(filtred_tasks, t)
			} else {
				found = true
			}
		}

		// Overwrite existing task slice
		tasks = filtred_tasks
		SaveTasks(filename, tasks)
		if found {
			fmt.Printf("Task with id %d deleted!\n", task_id)
		} else {
			fmt.Printf("Task with id %d not found!\n", task_id)
		}

	case "-search":
		// get string
		if len(os.Args) < 2 {
			fmt.Println("Serach requires one positional arg, -search <string>")
			return
		}
		query := os.Args[2]
		fmt.Printf("searching for \"%s\"...\n", query)

		// search logic
		cnt := 0
		for i := range tasks {
			if strings.Contains(strings.ToLower(tasks[i].Description), strings.ToLower(query)) {
				status := " "
				if tasks[i].Done {
					status = "x"
				}
				fmt.Printf("[%s] %d: %s\n", status, tasks[i].ID, tasks[i].Description)
				cnt++
			}
		}
		fmt.Printf("%d results found\n", cnt)

	default:
		fmt.Println("Unknown command")
	}
}
