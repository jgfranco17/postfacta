package validation

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jgfranco17/postfacta/api/httperror"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	gin.SetMode(gin.TestMode)
}

type testRequest struct {
	Name     string `json:"name" binding:"required,min=3,max=10"`
	Email    string `json:"email" binding:"required,email"`
	Age      int    `json:"age" binding:"required,min=18,max=100"`
	Category string `json:"category" binding:"required,oneof=A B C"`
}

func createTestContext(t *testing.T, body interface{}) (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	if body != nil {
		jsonData, err := json.Marshal(body)
		require.NoError(t, err)
		c.Request = httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader(jsonData))
		c.Request.Header.Set("Content-Type", "application/json")
	} else {
		c.Request = httptest.NewRequest(http.MethodPost, "/test", nil)
	}

	return c, w
}

func TestBindRequestValidData(t *testing.T) {
	validData := testRequest{
		Name:     "Alice",
		Email:    "alice@example.com",
		Age:      25,
		Category: "A",
	}

	c, _ := createTestContext(t, validData)
	var req testRequest

	err := BindRequest(c, &req)

	assert.NoError(t, err)
	assert.Equal(t, "Alice", req.Name)
	assert.Equal(t, "alice@example.com", req.Email)
	assert.Equal(t, 25, req.Age)
	assert.Equal(t, "A", req.Category)
}

func TestBindRequestValidationErrors(t *testing.T) {
	testCases := []struct {
		name              string
		data              interface{}
		expectedErrorMsgs []string
	}{
		{
			name: "missing required field",
			data: map[string]interface{}{
				"email":    "test@example.com",
				"age":      30,
				"category": "B",
			},
			expectedErrorMsgs: []string{"Validation failed", "name is required"},
		},
		{
			name: "min length violation",
			data: testRequest{
				Name:     "Al",
				Email:    "al@example.com",
				Age:      25,
				Category: "A",
			},
			expectedErrorMsgs: []string{"name must be at least 3 characters"},
		},
		{
			name: "max length violation",
			data: testRequest{
				Name:     "AliceWonderland",
				Email:    "alice@example.com",
				Age:      25,
				Category: "A",
			},
			expectedErrorMsgs: []string{"name must be at most 10 characters"},
		},
		{
			name: "oneof violation",
			data: testRequest{
				Name:     "Alice",
				Email:    "alice@example.com",
				Age:      25,
				Category: "D",
			},
			expectedErrorMsgs: []string{"category must be one of"},
		},
		{
			name: "min value violation",
			data: testRequest{
				Name:     "Alice",
				Email:    "alice@example.com",
				Age:      15,
				Category: "A",
			},
			expectedErrorMsgs: []string{"age must be at least 18"},
		},
		{
			name: "max value violation",
			data: testRequest{
				Name:     "Alice",
				Email:    "alice@example.com",
				Age:      150,
				Category: "A",
			},
			expectedErrorMsgs: []string{"age must be at most 100"},
		},
		{
			name: "multiple validation errors",
			data: map[string]interface{}{
				"name": "Al",
				"age":  15,
			},
			expectedErrorMsgs: []string{"Validation failed", "name must be at least 3", "email is required"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c, _ := createTestContext(t, tc.data)
			var req testRequest

			err := BindRequest(c, &req)

			assert.Error(t, err)
			var httpErr httperror.HttpError
			assert.ErrorAs(t, err, &httpErr)
			assert.Equal(t, http.StatusBadRequest, httpErr.Status())

			errorMsg := err.Error()
			for _, expectedMsg := range tc.expectedErrorMsgs {
				assert.Contains(t, errorMsg, expectedMsg)
			}
		})
	}
}

func TestBindRequestInvalidJSON(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader([]byte("{invalid json")))
	c.Request.Header.Set("Content-Type", "application/json")

	var req testRequest
	err := BindRequest(c, &req)

	assert.Error(t, err)
	var httpErr httperror.HttpError
	assert.ErrorAs(t, err, &httpErr)
	assert.Equal(t, http.StatusBadRequest, httpErr.Status())
	assert.Contains(t, err.Error(), "Invalid request body")
}

func TestBindRequestEmptyBody(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/test", nil)
	c.Request.Header.Set("Content-Type", "application/json")

	var req testRequest
	err := BindRequest(c, &req)

	assert.Error(t, err)
	var httpErr httperror.HttpError
	assert.ErrorAs(t, err, &httpErr)
	assert.Equal(t, http.StatusBadRequest, httpErr.Status())
}
