# List out available commands
_default:
	@just --list --unsorted

# Launch API in debug mode
start port="8080":
	@echo "Running main app..."
	go run . --port {{ port }}

test:
    @echo "[TEST] Running unit tests..."
    @go clean -testcache
    @go test -cover ./...
