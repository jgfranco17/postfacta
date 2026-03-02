package httperror

import (
	"context"
	"errors"
	"fmt"
)

// HttpError represents an HTTP-aware error with status code and request context.
// It implements the error interface and preserves request metadata for observability.
// The statusCode field determines the HTTP response status returned to the client.
// The cause field preserves the underlying error for error chain inspection.
type HttpError struct {
	message    string
	ctx        context.Context
	statusCode int
	cause      error
}

func (e HttpError) Error() string {
	return e.message
}

func (e HttpError) Context() context.Context {
	return e.ctx
}

func (e HttpError) Status() int {
	return e.statusCode
}

// Unwrap returns the underlying cause error, enabling errors.Is() and errors.As().
func (e HttpError) Unwrap() error {
	return e.cause
}

// New creates an HttpError with the given context, HTTP status code, and formatted message.
// The status code should follow RFC 7231 semantics.
//
// The context should contain request metadata (request ID, version) for error tracing.
// Invalid status codes (< 100 or > 599) are replaced with 500 during error handling.
//
// If the format string contains %w and wraps an error, that error is preserved as the cause,
// enabling errors.Is() and errors.As() to inspect the error chain.
func New(ctx context.Context, httpStatus int, format string, a ...any) HttpError {
	formattedErr := fmt.Errorf(format, a...)
	return HttpError{
		ctx:        ctx,
		statusCode: httpStatus,
		message:    formattedErr.Error(),
		cause:      errors.Unwrap(formattedErr),
	}
}
