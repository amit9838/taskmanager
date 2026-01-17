package display

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/amit9838/taskmanager/internal/task"
)

func PrintTasks(tasks []task.Task) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tStatus\tDescription\tCreated\tUpdated")
	fmt.Fprintln(w, "--\t------\t-----------\t-------\t-------")

	for _, t := range tasks {
		status := "[ ]"
		if t.Done {
			status = "[x]"
		}

		createdStr := t.CreatedAt.Format("2006-01-02")
		updatedStr := ""
		if !t.UpdatedAt.IsZero() {
			updatedStr = t.UpdatedAt.Format("2006-01-02")
		}

		fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\n",
			t.ID, status, t.Description, createdStr, updatedStr)
	}
	w.Flush()
}

func PrintTasksSimple(tasks []task.Task) {
	for _, t := range tasks {
		status := "[ ]"
		if t.Done {
			status = "[x]"
		}
		fmt.Printf("%s %d: %s\n", status, t.ID, t.Description)
	}
}
