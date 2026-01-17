# Variables
BINARY_NAME=tm
SRC_PATH=./cmd/taskmanager

.PHONY: all build install clean uninstall run

# Default action: build the binary locally
all: build

build:
	@echo "Building binary..."
	go build -o $(BINARY_NAME) $(SRC_PATH)

install:
	@echo "Installing to $(shell go env GOPATH)/bin..."
	go install $(SRC_PATH)
	@# If the folder name is 'taskmanager', 'go install' names it 'taskmanager'.
	@# This next line renames it to 'tm' in the bin folder if necessary.
	@if [ -f $(shell go env GOPATH)/bin/taskmanager ]; then \
		mv $(shell go env GOPATH)/bin/taskmanager $(shell go env GOPATH)/bin/$(BINARY_NAME); \
	fi
	@echo "Done! You can now run '$(BINARY_NAME)'"

uninstall:
	@echo "Removing binary..."
	@rm -f $(shell go env GOPATH)/bin/$(BINARY_NAME)
	@rm -f $(BINARY_NAME)
	@echo "Uninstalled successfully."

clean:
	@echo "Cleaning build cache and local binaries..."
	go clean
	rm -f $(BINARY_NAME)

run:
	go run $(SRC_PATH)


test:
	go test ./...