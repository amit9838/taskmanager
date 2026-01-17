package cli

import (
	"flag"
	"fmt"

	"github.com/amit9838/taskmanager/internal/task"
)

func ExecuteCommand(manager *task.TaskManager, args []string) error {
	if len(args) == 0 {
		printUsage()
		return nil
	}

	command := args[0]
	remainingArgs := args[1:]

	var cmd Command

	switch command {
	case "add":
		cmd = &AddCommand{}
		fs := flag.NewFlagSet("add", flag.ContinueOnError)
		if err := fs.Parse(remainingArgs); err != nil {
			return err
		}
		remainingArgs = fs.Args()

	case "list":
		cmd = &ListCommand{}

	case "done":
		cmd = &DoneCommand{}
		fs := flag.NewFlagSet("done", flag.ContinueOnError)
		if err := fs.Parse(remainingArgs); err != nil {
			return err
		}
		remainingArgs = fs.Args()

	case "del":
		cmd = &DeleteCommand{}
		fs := flag.NewFlagSet("del", flag.ContinueOnError)
		if err := fs.Parse(remainingArgs); err != nil {
			return err
		}
		remainingArgs = fs.Args()

	case "search":
		cmd = &SearchCommand{}
		fs := flag.NewFlagSet("search", flag.ContinueOnError)
		if err := fs.Parse(remainingArgs); err != nil {
			return err
		}
		remainingArgs = fs.Args()

	case "help":
		cmd = &HelpCommand{}

	default:
		return fmt.Errorf("unknown command: %s\nUse 'help' to see available commands", command)
	}

	return cmd.Execute(manager, remainingArgs)
}
