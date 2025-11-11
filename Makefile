# Project settings
APP_NAME=api
CMD_DIR=./cmd/api

# Default target
all: run

# Run the app
run:
	go run $(CMD_DIR)

# Build the binary into ./bin/api
build:
	mkdir -p bin
	go build -o bin/$(APP_NAME) $(CMD_DIR)

# Clean up build files
clean:
	rm -rf bin

# Run tests (if you have tests)
test:
	go test ./...

# Format all Go code
fmt:
	go fmt ./...

# Tidy up dependencies
tidy:
	go mod tidy
