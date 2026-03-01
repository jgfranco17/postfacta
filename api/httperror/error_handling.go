package httperror

import (
	"context"
	"errors"

	"github.com/jgfranco17/postfacta/api/logging"

	"github.com/gin-gonic/gin"
)

// ServiceError represents the structure of an error response sent to clients.
// It includes the HTTP status code and a body containing the error message
// and additional optional metadata.
type ServiceError struct {
	Status int
	Body   errorBody
}

type errorBody struct {
	Message        string `json:"message,omitempty"`
	RequestID      string `json:"requestId,omitempty"`
	ServiceVersion string `json:"serviceVersion,omitempty"`
}

// WithErrorHandling wraps a handler function that returns an error, automatically
// handling any returned errors by logging and sending an appropriate HTTP response.
func WithErrorHandling(handler func(c *gin.Context) error) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := handler(c)
		if err != nil {
			handleError(c, err)
		}
	}
}

func getContextField(ctx context.Context, fieldName string) string {
	value, ok := ctx.Value(fieldName).(string)
	if !ok {
		return ""
	}
	return value
}

func getErrorMetadataFromContext(ctx context.Context) errorBody {
	requestId := getContextField(ctx, logging.RequestId)
	serviceVersion := getContextField(ctx, logging.Version)

	return errorBody{
		RequestID:      requestId,
		ServiceVersion: serviceVersion,
	}
}

// extractErrorResponse converts an error into an HTTP error response.
// If the error is an HttpError (or wraps one), it extracts the status code and request metadata.
// Status codes outside the valid HTTP range (100-599) are replaced with 500.
// Non-HttpError instances are treated as internal server errors (500) with a generic message.
func extractErrorResponse(ctx context.Context, err error) ServiceError {
	errorMessage := err.Error()

	var httpErrInstance HttpError
	if errors.As(err, &httpErrInstance) {
		body := getErrorMetadataFromContext(httpErrInstance.Context())
		body.Message = errorMessage
		status := httpErrInstance.Status()
		if status < 100 || status > 599 {
			status = 500
		}
		return ServiceError{Status: status, Body: body}
	}

	body := getErrorMetadataFromContext(ctx)
	body.Message = "Internal Server Error"
	return ServiceError{Status: 500, Body: body}
}

// handleError logs the error and sends the appropriate HTTP response to the client.
func handleError(c *gin.Context, err error) {
	log := logging.FromContext(c)
	log.Error(err)
	errorResponse := extractErrorResponse(c, err)
	c.JSON(errorResponse.Status, errorResponse.Body)
}
