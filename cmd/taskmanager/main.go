package main

import (
	"fmt"
	"os"

	"github.com/amit9838/taskmanager/internal/storage"
	"github.com/amit9838/taskmanager/internal/task"
	"github.com/amit9838/taskmanager/pkg/cli"
)

func main() {
	// Initialize storage
	fileStorage := storage.NewJSONStorage("tasks.json")

	// Initialize task manager
	taskManager, err := task.NewTaskManager(fileStorage)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing task manager: %v\n", err)
		os.Exit(1)
	}

	// Parse and execute command
	if err := cli.ExecuteCommand(taskManager, os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
