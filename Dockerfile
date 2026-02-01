# syntax=docker/dockerfile:1

ARG PYTHON_VERSION=3.12.11
FROM python:${PYTHON_VERSION}-slim AS build

# Configure env for headless Python service
ENV PYTHONDONTWRITEBYTECODE=1
ENV PYTHONUNBUFFERED=1
COPY --from=ghcr.io/astral-sh/uv:0.9.28 /uv /uvx /bin/

WORKDIR /app

# Create a non-privileged user that the app will run under.
# See https://docs.docker.com/go/dockerfile-user-best-practices/
ARG UID=10001
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    appuser

# Install dependencies
COPY ./pyproject.toml ./uv.lock /app/
RUN uv sync --locked


FROM build AS service

# Switch to the non-privileged user to run the application.
USER appuser

COPY ./postfacta /app/postfacta

EXPOSE 8000

CMD [ "uv", "run", "uvicorn", "postfacta.service:app", "--port", "8000", "--host", "0.0.0.0" ]
