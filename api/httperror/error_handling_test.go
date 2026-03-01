package httperror

import (
	"context"
	"fmt"
	"testing"

	"github.com/jgfranco17/postfacta/api/logging"

	"github.com/stretchr/testify/assert"
)

func TestHandleHTTPError(t *testing.T) {
	t.Run("Simple input error", func(t *testing.T) {
		inputErr := New(t.Context(), 500, "Some error")
		response := extractErrorResponse(t.Context(), inputErr)

		assert.Equal(t, 500, response.Status)
		assert.Equal(t, errorBody{
			Message: "Some error",
		}, response.Body)
	})

	t.Run("HTTP error wrapped in generic error", func(t *testing.T) {
		inputErr := New(t.Context(), 500, "Some error")

		err := fmt.Errorf("Outer error: %w", inputErr)
		response := extractErrorResponse(t.Context(), err)

		assert.Equal(t, 500, response.Status)
		assert.Equal(t, errorBody{
			Message: "Outer error: Some error",
		}, response.Body)
	})

	t.Run("HTTP error wrapping generic error", func(t *testing.T) {
		err := fmt.Errorf("Inner error")
		inputErr := New(t.Context(), 500, "Some error: %w", err)

		response := extractErrorResponse(t.Context(), inputErr)

		assert.Equal(t, 500, response.Status)
		assert.Equal(t, errorBody{
			Message: "Some error: Inner error",
		}, response.Body)
	})

	t.Run("HTTP error with requestId", func(t *testing.T) {
		ctx := context.WithValue(t.Context(), logging.RequestId, "4dfdcc88-2f3e-41ce-9757-4144cb3974a4")

		inputErr := New(ctx, 500, "Some error")
		response := extractErrorResponse(t.Context(), inputErr)

		assert.Equal(t, 500, response.Status)
		assert.Equal(t, errorBody{
			Message:   "Some error",
			RequestID: "4dfdcc88-2f3e-41ce-9757-4144cb3974a4",
		}, response.Body)
	})

	t.Run("HTTP error with service version", func(t *testing.T) {
		ctx := context.WithValue(t.Context(), logging.Version, "1.23.5")

		inputErr := New(ctx, 500, "Some error")
		response := extractErrorResponse(t.Context(), inputErr)

		assert.Equal(t, 500, response.Status)
		assert.Equal(t, errorBody{
			Message:        "Some error",
			ServiceVersion: "1.23.5",
		}, response.Body)
	})

	t.Run("Status code preservation for 400 Bad Request", func(t *testing.T) {
		inputErr := New(t.Context(), 400, "Invalid input")
		response := extractErrorResponse(t.Context(), inputErr)

		assert.Equal(t, 400, response.Status)
		assert.Equal(t, "Invalid input", response.Body.Message)
	})

	t.Run("Status code preservation for 404 Not Found", func(t *testing.T) {
		inputErr := New(t.Context(), 404, "Resource not found")
		response := extractErrorResponse(t.Context(), inputErr)

		assert.Equal(t, 404, response.Status)
		assert.Equal(t, "Resource not found", response.Body.Message)
	})

	t.Run("Status code preservation for 409 Conflict", func(t *testing.T) {
		inputErr := New(t.Context(), 409, "Resource already exists")
		response := extractErrorResponse(t.Context(), inputErr)

		assert.Equal(t, 409, response.Status)
		assert.Equal(t, "Resource already exists", response.Body.Message)
	})

	t.Run("Invalid status code defaults to 500", func(t *testing.T) {
		inputErr := New(t.Context(), 0, "Invalid status code")
		response := extractErrorResponse(t.Context(), inputErr)

		assert.Equal(t, 500, response.Status)
		assert.Equal(t, "Invalid status code", response.Body.Message)
	})

	t.Run("Negative status code defaults to 500", func(t *testing.T) {
		inputErr := New(t.Context(), -1, "Negative status")
		response := extractErrorResponse(t.Context(), inputErr)

		assert.Equal(t, 500, response.Status)
		assert.Equal(t, "Negative status", response.Body.Message)
	})

	t.Run("Out of range status code defaults to 500", func(t *testing.T) {
		inputErr := New(t.Context(), 600, "Out of range")
		response := extractErrorResponse(t.Context(), inputErr)

		assert.Equal(t, 500, response.Status)
		assert.Equal(t, "Out of range", response.Body.Message)
	})

	t.Run("Non-HttpError returns 500", func(t *testing.T) {
		err := fmt.Errorf("generic error")
		response := extractErrorResponse(t.Context(), err)

		assert.Equal(t, 500, response.Status)
		assert.Equal(t, "Internal Server Error", response.Body.Message)
	})
}
