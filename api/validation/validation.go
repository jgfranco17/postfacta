package validation

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jgfranco17/postfacta/api/httperror"
)

// BindRequest binds JSON request body to the target struct and validates it.
// It returns an HttpError with status 400 if binding or validation fails.
// The error message includes specific field-level validation failures.
func BindRequest(c *gin.Context, req interface{}) error {
	if err := c.ShouldBindJSON(req); err != nil {
		var validationErrs validator.ValidationErrors
		if errors.As(err, &validationErrs) {
			message := formatValidationErrors(validationErrs)
			return httperror.New(c, http.StatusBadRequest, "%s", message)
		}
		return httperror.New(c, http.StatusBadRequest, "Invalid request body: %s", err.Error())
	}
	return nil
}

// formatValidationErrors converts validator errors into a user-friendly message.
func formatValidationErrors(errs validator.ValidationErrors) string {
	var messages []string

	for _, err := range errs {
		field := strings.ToLower(err.Field())

		switch err.Tag() {
		case "required":
			messages = append(messages, fmt.Sprintf("%s is required", field))
		case "min":
			messages = append(messages, fmt.Sprintf("%s must be at least %s characters", field, err.Param()))
		case "max":
			messages = append(messages, fmt.Sprintf("%s must be at most %s characters", field, err.Param()))
		case "oneof":
			messages = append(messages, fmt.Sprintf("%s must be one of [%s]", field, err.Param()))
		default:
			messages = append(messages, fmt.Sprintf("%q validation failed", field))
		}
	}

	return fmt.Sprintf("Validation failed: %s", strings.Join(messages, ", "))
}
