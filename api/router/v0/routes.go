package v0

import (
	"github.com/jgfranco17/postfacta/api/db"
	"github.com/jgfranco17/postfacta/api/httperror"

	"github.com/gin-gonic/gin"
)

// Adds v0 routes to the router.
func SetRoutes(route *gin.Engine, dbClient db.DatabaseClient) error {
	v0 := route.Group("/v0")
	{
		testExecutionRoutes := v0.Group("/tests")
		{
			testExecutionRoutes.POST("/run", httperror.WithErrorHandling(runTests(dbClient)))
		}
	}
	return nil
}
