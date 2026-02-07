package logging

import (
	"context"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/jgfranco17/postfacta/api/environment"
	"github.com/sirupsen/logrus"
)

type contextLogKey string

const contextKey contextLogKey = "logger"

// List of keys needed by the core functionality
const (
	RequestId   string = "requestId"
	Version     string = "version"
	Environment string = "environment"
	Origin      string = "origin"
)

func New(stream io.Writer, level logrus.Level) *logrus.Logger {
	logger := logrus.New()
	logger.SetOutput(stream)
	logger.SetLevel(level)

	formatter := environment.GetLogFormatter()
	logger.SetFormatter(formatter)

	return logger
}

type RequestMetadata struct {
	RequestID   string
	Environment string
	Version     string
}

func FillFields(c *gin.Context, fields RequestMetadata) {
	// Golang recommends contexts to use custom types instead
	// of strings, but gin defines key as a string.
	c.Set(string(RequestId), fields.RequestID)
	c.Set(string(Environment), fields.Environment)
	c.Set(string(Version), fields.Version)
}

func AddToContext(ctx context.Context, logger *logrus.Logger) context.Context {
	return context.WithValue(ctx, contextKey, logger)
}

func FromContext(ctx context.Context) *logrus.Logger {
	if logger, ok := ctx.Value(contextKey).(*logrus.Logger); ok {
		return logger
	}

	// Try gin.Context if it's that type
	if ginCtx, ok := ctx.(*gin.Context); ok {
		if logger, ok := ginCtx.Get(string(contextKey)); ok {
			if logrusLogger, ok := logger.(*logrus.Logger); ok {
				return logrusLogger
			}
		}
	}

	panic("no logger set in context")
}

// AddToRequestContext adds logger to gin.Context
func AddToRequestContext(c *gin.Context, logger *logrus.Logger) {
	c.Set(string(contextKey), logger)
}

// FromRequestContext retrieves logger from gin.Context
func FromRequestContext(c *gin.Context) *logrus.Logger {
	if logger, ok := c.Get(string(contextKey)); ok {
		if logrusLogger, ok := logger.(*logrus.Logger); ok {
			return logrusLogger
		}
	}
	panic("no logger set in gin context")
}
