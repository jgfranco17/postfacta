package system

import (
	"fmt"
	"net/http"
	"time"

	"github.com/jgfranco17/postfacta/api/environment"
	"github.com/jgfranco17/postfacta/api/httperror"

	"github.com/gin-gonic/gin"
)

func HomeHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Welcome to the PostFacta API!",
	})
}

func ServiceInfoHandler(codebaseSpec *ProjectCodebase, startTime time.Time) func(c *gin.Context) {
	return func(c *gin.Context) {
		timeSinceStart := time.Since(startTime)
		uptimeSeconds := fmt.Sprintf("%ds", int(timeSinceStart.Seconds()))
		c.JSON(http.StatusOK, ServiceInfo{
			Name:        "PostFacta API",
			Codebase:    *codebaseSpec,
			Environment: environment.GetApplicationEnv(),
			Uptime:      uptimeSeconds,
		})
	}
}

func HealthCheckHandler() func(c *gin.Context) {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, HealthStatus{
			Status: "healthy",
		})
	}
}

func NotFoundHandler(c *gin.Context) {
	httperror.RespondWithError(c, http.StatusNotFound, "Endpoint '%s' does not exist", c.Request.URL.Path)
}
