package httperror

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jgfranco17/postfacta/api/logging"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestWithErrorHandling(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		name           string
		handler        func(c *gin.Context) error
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "handler returns no error",
			handler: func(c *gin.Context) error {
				c.JSON(http.StatusOK, gin.H{"status": "success"})
				return nil
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"status":"success"}`,
		},
		{
			name: "handler returns HttpError 400",
			handler: func(c *gin.Context) error {
				return New(c, http.StatusBadRequest, "invalid input")
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"invalid input"}`,
		},
		{
			name: "handler returns HttpError 404",
			handler: func(c *gin.Context) error {
				return New(c, http.StatusNotFound, "not found")
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"message":"not found"}`,
		},
		{
			name: "handler returns generic error",
			handler: func(c *gin.Context) error {
				return errors.New("something went wrong")
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"message":"Internal Server Error"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)

			// Set up logging context to avoid nil pointer
			logger := logging.New(io.Discard, logrus.DebugLevel)
			logging.AddToRequestContext(c, logger)

			handlerFunc := WithErrorHandling(tc.handler)
			handlerFunc(c)

			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.JSONEq(t, tc.expectedBody, w.Body.String())
		})
	}
}

func TestRespondWithError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		name           string
		status         int
		format         string
		args           []any
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "simple error message",
			status:         http.StatusBadRequest,
			format:         "invalid input",
			args:           nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"invalid input"}`,
		},
		{
			name:           "formatted error message",
			status:         http.StatusNotFound,
			format:         "user %s not found",
			args:           []any{"john"},
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"message":"user john not found"}`,
		},
		{
			name:           "multiple format arguments",
			status:         http.StatusConflict,
			format:         "conflict between %s and %s",
			args:           []any{"item1", "item2"},
			expectedStatus: http.StatusConflict,
			expectedBody:   `{"message":"conflict between item1 and item2"}`,
		},
		{
			name:           "internal server error",
			status:         http.StatusInternalServerError,
			format:         "database connection failed",
			args:           nil,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"message":"database connection failed"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)

			var buf bytes.Buffer
			logger := logging.New(&buf, logrus.DebugLevel)
			logging.AddToRequestContext(c, logger)

			RespondWithError(c, tc.status, tc.format, tc.args...)

			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.JSONEq(t, tc.expectedBody, w.Body.String())
		})
	}
}

func TestHandleError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		name           string
		setupContext   func(c *gin.Context) context.Context
		error          error
		expectedStatus int
		validateBody   func(t *testing.T, body string)
	}{
		{
			name: "HttpError with no metadata",
			setupContext: func(c *gin.Context) context.Context {
				return c
			},
			error:          New(context.Background(), http.StatusBadRequest, "bad request"),
			expectedStatus: http.StatusBadRequest,
			validateBody: func(t *testing.T, body string) {
				assert.JSONEq(t, `{"message":"bad request"}`, body)
			},
		},
		{
			name: "HttpError with request ID",
			setupContext: func(c *gin.Context) context.Context {
				return context.WithValue(c, logging.RequestId, "test-request-id")
			},
			error:          nil, // will be created with enriched context
			expectedStatus: http.StatusNotFound,
			validateBody: func(t *testing.T, body string) {
				assert.JSONEq(t, `{"message":"not found","requestId":"test-request-id"}`, body)
			},
		},
		{
			name: "HttpError with service version",
			setupContext: func(c *gin.Context) context.Context {
				return context.WithValue(c, logging.Version, "1.0.0")
			},
			error:          nil, // will be created with enriched context
			expectedStatus: http.StatusInternalServerError,
			validateBody: func(t *testing.T, body string) {
				assert.JSONEq(t, `{"message":"internal error","serviceVersion":"1.0.0"}`, body)
			},
		},
		{
			name: "generic error",
			setupContext: func(c *gin.Context) context.Context {
				return c
			},
			error:          errors.New("something went wrong"),
			expectedStatus: http.StatusInternalServerError,
			validateBody: func(t *testing.T, body string) {
				assert.JSONEq(t, `{"message":"Internal Server Error"}`, body)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)

			var buf bytes.Buffer
			logger := logging.New(&buf, logrus.DebugLevel)
			logging.AddToRequestContext(c, logger)

			// Set up context with metadata
			ctx := tc.setupContext(c)

			// Create error with enriched context if needed
			var err error
			if tc.error != nil {
				err = tc.error
			} else {
				// Create error with the enriched context
				switch tc.expectedStatus {
				case http.StatusNotFound:
					err = New(ctx, http.StatusNotFound, "not found")
				case http.StatusInternalServerError:
					err = New(ctx, http.StatusInternalServerError, "internal error")
				}
			}

			require.NotNil(t, err)
			handleError(c, err)

			assert.Equal(t, tc.expectedStatus, w.Code)
			tc.validateBody(t, w.Body.String())
		})
	}
}
