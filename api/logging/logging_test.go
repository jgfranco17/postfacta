package logging

import (
	"bytes"
	"context"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestContextOperations(t *testing.T) {
	testCases := []struct {
		name     string
		testFunc func(t *testing.T)
	}{
		{
			name: "AddToContext and FromContext",
			testFunc: func(t *testing.T) {
				var buf bytes.Buffer
				logger := New(&buf, logrus.TraceLevel)
				ctx := AddToContext(context.Background(), logger)
				assert.Equal(t, logger, FromContext(ctx))
			},
		},
		{
			name: "FromContext panics when logger not set",
			testFunc: func(t *testing.T) {
				ctx := context.Background()
				assert.Panics(t, func() {
					FromContext(ctx)
				})
			},
		},
		{
			name: "FromContext with gin.Context",
			testFunc: func(t *testing.T) {
				var buf bytes.Buffer
				logger := New(&buf, logrus.InfoLevel)

				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)
				AddToRequestContext(c, logger)

				retrieved := FromContext(c)
				assert.Equal(t, logger, retrieved)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, tc.testFunc)
	}
}

func TestRequestContextOperations(t *testing.T) {
	testCases := []struct {
		name     string
		testFunc func(t *testing.T)
	}{
		{
			name: "AddToRequestContext",
			testFunc: func(t *testing.T) {
				var buf bytes.Buffer
				logger := New(&buf, logrus.DebugLevel)

				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				AddToRequestContext(c, logger)

				value, exists := c.Get(string(contextKey))
				require.True(t, exists)
				assert.Equal(t, logger, value)
			},
		},
		{
			name: "FromRequestContext",
			testFunc: func(t *testing.T) {
				var buf bytes.Buffer
				logger := New(&buf, logrus.WarnLevel)

				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)
				AddToRequestContext(c, logger)

				retrieved := FromRequestContext(c)
				assert.Equal(t, logger, retrieved)
			},
		},
		{
			name: "FromRequestContext panics when logger not set",
			testFunc: func(t *testing.T) {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)

				assert.Panics(t, func() {
					FromRequestContext(c)
				})
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, tc.testFunc)
	}
}

func TestFillFields(t *testing.T) {
	testCases := []struct {
		name     string
		metadata RequestMetadata
	}{
		{
			name: "with all fields populated",
			metadata: RequestMetadata{
				RequestID:   "req-123",
				Environment: "test",
				Version:     "1.0.0",
			},
		},
		{
			name: "with empty fields",
			metadata: RequestMetadata{
				RequestID:   "",
				Environment: "",
				Version:     "",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			FillFields(c, tc.metadata)

			requestID, exists := c.Get(RequestId)
			require.True(t, exists)
			assert.Equal(t, tc.metadata.RequestID, requestID)

			env, exists := c.Get(Environment)
			require.True(t, exists)
			assert.Equal(t, tc.metadata.Environment, env)

			version, exists := c.Get(Version)
			require.True(t, exists)
			assert.Equal(t, tc.metadata.Version, version)
		})
	}
}

func TestNewLogger_DifferentLevels(t *testing.T) {
	testCases := []struct {
		name  string
		level logrus.Level
	}{
		{"Trace level", logrus.TraceLevel},
		{"Debug level", logrus.DebugLevel},
		{"Info level", logrus.InfoLevel},
		{"Warn level", logrus.WarnLevel},
		{"Error level", logrus.ErrorLevel},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			logger := New(&buf, tc.level)

			assert.NotNil(t, logger)
			assert.Equal(t, tc.level, logger.Level)
		})
	}
}
