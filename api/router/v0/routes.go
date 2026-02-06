package v0

import (
	"github.com/jgfranco17/postfacta/api/db"
	"github.com/jgfranco17/postfacta/api/httperror"

	"github.com/gin-gonic/gin"
)

// Adds v0 routes to the router.
func SetRoutes(route *gin.RouterGroup, dbClient db.DatabaseClient) {
	v0 := route.Group("/v0")
	v0.GET("/incidents", httperror.WithErrorHandling(getAllIncidents(dbClient)))
}
