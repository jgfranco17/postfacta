# POSTFACTA

COVERAGE_OUTPUT := "coverage.out"

# List out available commands
_default:
	@just --list --unsorted

# Launch API in debug mode
start port="8080":
	@echo "Running main app..."
	go run . --port {{ port }}

# Run unit tests
test:
    #!/usr/bin/env bash
    echo "[TEST] Running unit tests..."
    go clean -testcache
    go test -cover -race -shuffle=on ./...

# Run BDD-style integration tests
integration:
    #!/usr/bin/env bash
    echo "[TEST] Running integration tests..."
    go clean -testcache
    ginkgo -v ./api/integration/...

# Generate total coverage report (unit + integration)
coverage threshold="70":
    #!/usr/bin/env bash
    echo "[TEST] Generating total coverage across all tests..."
    go clean -testcache
    go test -coverprofile={{ COVERAGE_OUTPUT }} -coverpkg="./api/..." ./...
    total=$(go tool cover -func={{ COVERAGE_OUTPUT }} | awk '/^total:/ {print $3}' | tr -d '%')
    if [[ -z "$total" ]]; then
        echo "[TEST] Failed to parse total coverage"
        exit 1
    fi
    threshold="{{ threshold }}"
    if (( $(echo "$total >= $threshold" | bc -l) )); then
        echo "[TEST] Coverage ${total}% meets threshold ${threshold}%"
    else
        echo "[TEST] Coverage ${total}% is below threshold ${threshold}%"
        exit 1
    fi
