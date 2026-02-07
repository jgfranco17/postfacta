package system

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var startTime time.Time

func init() {
	startTime = time.Now()
}

func SetSystemRoutes(route *gin.Engine, data []byte) error {
	startTime = time.Now()

	if data != nil {
		specs, err := getCodebaseSpec(data)
		if err != nil {
			return err
		}
		route.GET("/service-info", ServiceInfoHandler(specs, startTime))
	}
	for _, homeRoute := range []string{"", "/home"} {
		route.GET(homeRoute, HomeHandler)
	}

	route.GET("/healthz", HealthCheckHandler())
	route.GET("/metrics", gin.WrapH(promhttp.Handler()))
	route.NoRoute(NotFoundHandler)

	return nil
}
