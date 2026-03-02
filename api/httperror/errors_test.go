package httperror

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type customError struct {
	code int
}

func (ce *customError) Error() string {
	return fmt.Sprintf("custom error with code %d", ce.code)
}

func TestHttpErrorNewSimpleError(t *testing.T) {
	const inputErrorMessage string = "This is an input error"

	err := New(context.Background(), 400, inputErrorMessage)

	var expectedError HttpError
	assert.ErrorAs(t, err, &expectedError)
	assert.Equal(t, inputErrorMessage, err.Error())
}

func TestHttpErrorNewWrappedError(t *testing.T) {
	const rootMessage string = "This is the root"
	inputErrorMessage := "This is an input error: %v"

	rootError := fmt.Errorf(rootMessage)

	err := New(context.Background(), 500, inputErrorMessage, rootError)

	var expectedError HttpError
	assert.ErrorAs(t, err, &expectedError)
	assert.Equal(t, "This is an input error: This is the root", err.Error())
}

func TestHttpErrorPreservesCauseWithWrapping(t *testing.T) {
	rootCause := fmt.Errorf("database connection failed")
	wrapped := New(context.Background(), 500, "Failed to query: %w", rootCause)

	assert.Equal(t, "Failed to query: database connection failed", wrapped.Error())
	assert.Equal(t, rootCause, errors.Unwrap(wrapped))
	assert.ErrorAs(t, wrapped, &rootCause)
}

func TestHttpErrorErrorsIsWorksAcrossHttpError(t *testing.T) {
	var sentinelErr = errors.New("not found")
	httpErr := New(context.Background(), 404, "Resource not found: %w", sentinelErr)

	assert.ErrorAs(t, httpErr, &sentinelErr)
}

func TestHttpErrorErrorsAsWorksAcrossHttpError(t *testing.T) {
	customErr := &customError{code: 42}
	wrappedErr := fmt.Errorf("custom error occurred: %w", customErr)
	httpErr := New(context.Background(), 500, "Handler failed: %w", wrappedErr)

	var target *customError
	assert.ErrorAs(t, httpErr, &target)
	assert.Equal(t, 42, target.code)
}

func TestHttpErrorUnwrapReturnsNilForSimpleError(t *testing.T) {
	simpleErr := New(context.Background(), 400, "simple error without cause")

	assert.Nil(t, errors.Unwrap(simpleErr))
}

func TestHttpErrorUnwrapReturnsNilForNonWrappingFormat(t *testing.T) {
	rootErr := fmt.Errorf("root cause")
	nonWrapping := New(context.Background(), 500, "Error occurred: %v", rootErr)

	assert.Nil(t, errors.Unwrap(nonWrapping))
	assert.Equal(t, "Error occurred: root cause", nonWrapping.Error())
}

func TestHttpErrorMultiLayerWrapping(t *testing.T) {
	layer1 := errors.New("original error")
	layer2 := fmt.Errorf("wrapped once: %w", layer1)
	layer3 := fmt.Errorf("wrapped twice: %w", layer2)
	httpErr := New(context.Background(), 500, "HTTP wrapper: %w", layer3)

	assert.ErrorAs(t, httpErr, &layer1)
	assert.ErrorAs(t, httpErr, &layer2)
	assert.ErrorAs(t, httpErr, &layer3)
}
