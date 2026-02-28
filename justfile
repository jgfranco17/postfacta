# List out available commands
_default:
	@just --list --unsorted

# Launch API in debug mode
start port="8080":
	@echo "Running main app..."
	go run . --port {{ port }}

# Run unit tests
test:
    @echo "[TEST] Running unit tests..."
    @go clean -testcache
    @go test -cover ./...

# Run BDD-style integration tests
integration:
	@echo "[TEST] Running integration tests..."
	@ginkgo -v ./api/integration/...
