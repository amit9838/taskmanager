# Task Manager (tm)

A lightweight, CLI-based task management application written in Go. This project demonstrates a clean Go project structure, local package imports, testing and JSON-based persistent storage.

## 📁 Project Structure

```text
.
├── cmd/
│   └── taskmanager/
│       └── main.go         # Entry point of the application
├── internal/               # Private project code
│   ├── storage/
│   │   └── storage.go      # JSON persistence logic
│   └── task/
│       ├── task.go         # Task struct definition
│       └── task_manager.go # Task list manipulation logic
├── pkg/                    # Public library code
│   ├── cli/
│   │   └── commands.go     # CLI argument parsing
│   └── display/
│       └── display.go      # Terminal output formatting
├── go.mod                  # Go module definition
├── Makefile                # Automation scripts
└── tasks.json              # Data storage (auto-generated)

```

---

## 🚀 Getting Started

### Prerequisites

* **Go**: 1.25+
* **Make**: (Optional, for using the Makefile)

### Installation

1. **Clone the repository:**
```bash
git clone https://github.com/amit9838/taskmanager.git
cd taskmanager

```


2. **Install the binary globally:**
```bash
make install

```


*Note: This installs the binary as `tm` in your `$(go env GOPATH)/bin`.*
3. **Verify Installation:**
```bash
tm --help

```



---

## 🛠 Usage

You can run the application using the `tm` command.

```shell
# Add a new task
tm add "Buy groceries"
tm add "Finish project report"
tm add "Call mom"

# List all tasks
tm list

# Mark task as done
tm done 1

# Delete a task
tm del 2

# Search for tasks
tm search "groceries"

# Show help
tm help
```

Task list
---
run `tm list`.
```
ID  Status  Description                 Created
--  ------  -----------                 -------
1   [x]     Buy groceries               2024-01-15
2   [ ]     Finish project report       2024-01-15
3   [ ]     Call mom                    2024-01-15
```

Search Project
---
```shell
tm search project
# Output: Found 1 results:
# [ ] 2: Finish project report
```

## 💾 Storage Logic

The application stores data in a `tasks.json` file.

* **Auto-Initialization:** If the file does not exist, the application will automatically create it with an empty list `[]`.
* **Resilience:** The application handles empty files and whitespace gracefully to prevent JSON decoding errors.

---

## 🏗 Development

### Makefile Commands

We use a `Makefile` to simplify common tasks:

* `make build`: Compiles the binary to the current directory.
* `make run`: Runs the application directly from source.
* `make clean`: Removes binaries and build artifacts.
* `make uninstall`: Removes the `tm` binary from your system path.

* `go test -v ./internal/task` : Vourbose testing on a particular file.

---
