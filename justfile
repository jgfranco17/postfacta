# List out available commands
_default:
	@just --list --unsorted

# Launch API in debug mode
start:
	@echo "Running main app..."
	@go run .

# Run unit tests
pytest *args:
	@echo "Running unittest suite..."
	@uv run pytest {{ args }}

# Run integration tests
integration-pytest *args:
	#!/usr/bin/env bash
	echo "Running integration test suite..."
	RUN_INTEGRATION="true" uv run pytest -m integration {{ args }}

# Run test coverage
coverage:
    @uv run coverage run -m pytest
    @uv run coverage report

# Run UV Python
python *args:
    @uv run python {{ args }}

# Tidy and lint code
lint:
    @echo "Running code formatter (black)..."
    @uv run black .
    @echo "Running linter (flake8)..."
    @uv run flake8 .
    @echo "Running type checker (mypy)..."
    @uv run mypy .
    @echo "Code checking complete!"
