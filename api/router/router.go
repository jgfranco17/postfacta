package router

import (
	"context"
	"fmt"
	"os"

	"github.com/jgfranco17/postfacta/api/db"
	env "github.com/jgfranco17/postfacta/api/environment"
	"github.com/jgfranco17/postfacta/api/logging"
	"github.com/jgfranco17/postfacta/api/router/headers"
	system "github.com/jgfranco17/postfacta/api/router/system"
	v0 "github.com/jgfranco17/postfacta/api/router/v0"
	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Service struct {
	Router *gin.Engine
	Port   int
}

func (s *Service) Run(ctx context.Context) error {
	logger := logging.FromContext(ctx)
	logger.WithFields(logrus.Fields{"port": s.Port}).Infof("Starting service")

	if err := s.Router.Run(fmt.Sprintf(":%d", s.Port)); err != nil {
		return fmt.Errorf("Failed to start service on port %v: %w", s.Port, err)
	}
	return nil
}

// Add the fields we want to expose in the logger to the request context
func addLoggerFields() gin.HandlerFunc {
	level := env.GetLogLevel()
	logger := logging.New(os.Stderr, level)

	return func(c *gin.Context) {
		logging.AddToRequestContext(c, logger)

		if !env.IsRunningLocally() {
			requestID := uuid.NewString()
			environment := os.Getenv(env.ENV_KEY_ENVIRONMENT)
			version := os.Getenv(env.ENV_KEY_VERSION)

			logging.FillFields(c, logging.RequestMetadata{
				RequestID:   requestID,
				Environment: environment,
				Version:     version,
			})

			originInfo, err := headers.CreateOriginInfoHeader(c)

			if err == nil && originInfo.Origin != "" {
				c.Set(string(logging.Origin), fmt.Sprintf("%s@%s", originInfo.Origin, originInfo.Version))
			}
		}
		c.Next()
	}
}

// Log the start and completion of a request
func logRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := logging.FromRequestContext(c)

		origin := c.Request.Header.Get("Origin")
		log.Infof("Request Started: [%s] %s from %s", c.Request.Method, c.Request.URL, origin)
		c.Next()
		log.Infof("Request Completed: [%s] %s", c.Request.Method, c.Request.URL)
	}
}

// Configure the router adding routes and middlewares
func getRouter(ctx context.Context, dbClient db.DatabaseClient, metadata []byte) (*gin.Engine, error) {
	logger := logging.FromContext(ctx)
	router := gin.Default()

	router.Use(addLoggerFields())
	router.Use(logRequest())
	router.Use(GetCors())
	router.Use(system.PrometheusMiddleware())
	if err := system.SetSystemRoutes(router, metadata); err != nil {
		return nil, err
	}

	apiBaseGroup := router.Group("/api")
	v0.SetRoutes(apiBaseGroup, dbClient)

	logger.Trace("Router configured")
	return router, nil
}

// CreateNewService configures a new service instance.
func CreateNewService(ctx context.Context, port int, dbClient db.DatabaseClient, metadata []byte) (*Service, error) {
	router, err := getRouter(ctx, dbClient, metadata)
	if err != nil {
		return nil, fmt.Errorf("Failed to create new service instance: %w", err)
	}

	return &Service{
		Router: router,
		Port:   port,
	}, nil
}
