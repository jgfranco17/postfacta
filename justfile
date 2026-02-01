# List out available commands
_default:
	@just --list

# Launch API in debug mode
start:
	@echo "Running main app..."
	@uv run uvicorn postfacta.service:app --port 8000 --reload

# Clean unused files
clean:
	-@find ./ -name '*.pyc' -exec rm -f {} \;
	-@find ./ -name '__pycache__' -exec rm -rf {} \;
	-@find ./ -name 'Thumbs.db' -exec rm -f {} \;
	-@find ./ -name '*~' -exec rm -f {} \;
	-@rm -rf .pytest_cache
	-@rm -rf .cache
	-@rm -rf .mypy_cache
	-@rm -rf build
	-@rm -rf dist
	-@rm -rf *.egg-info
	-@rm -rf htmlcov
	-@rm -rf .tox/
	-@rm -rf docs/_build
	-@rm -rf .venv
	@echo "Cleaned out unused files and directories!"

# Run unit tests
pytest *args:
	@echo "Running unittest suite..."
	uv run pytest {{ args }}

# Run test coverage
coverage:
    @uv run coverage run -m pytest
    @uv run coverage report
