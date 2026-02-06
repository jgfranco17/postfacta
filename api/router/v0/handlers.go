package v0

import (
	"net/http"

	"github.com/jgfranco17/postfacta/api/db"

	"github.com/gin-gonic/gin"
)

func runTests(dbClient db.DatabaseClient) func(c *gin.Context) error {
	return func(c *gin.Context) error {
		c.JSON(http.StatusOK, gin.H{
			"message": "Test execution started",
		})
		return nil
	}
}
