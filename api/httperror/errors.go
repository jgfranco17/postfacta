package httperror

import (
	"context"
	"fmt"
)

// HttpError represents an HTTP-aware error with status code and request context.
// It implements the error interface and preserves request metadata for observability.
// The statusCode field determines the HTTP response status returned to the client.
type HttpError struct {
	message    string
	ctx        context.Context
	statusCode int
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

// New creates an HttpError with the given context, HTTP status code, and formatted message.
// The status code should follow RFC 7231 semantics.
//
// The context should contain request metadata (request ID, version) for error tracing.
// Invalid status codes (< 100 or > 599) are replaced with 500 during error handling.
func New(ctx context.Context, httpStatus int, format string, a ...any) HttpError {
	return HttpError{
		ctx:        ctx,
		statusCode: httpStatus,
		message:    fmt.Errorf(format, a...).Error(),
	}
}
