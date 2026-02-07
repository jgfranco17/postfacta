package main

import (
	"context"
	"flag"
	"os"

	"github.com/jgfranco17/postfacta/api/db"
	"github.com/jgfranco17/postfacta/api/logging"
	"github.com/jgfranco17/postfacta/api/router"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	_ "embed"
)

var (
	port    = flag.Int("port", 8080, "Port to listen on")
	devMode = flag.Bool("dev", true, "Run server in debug mode")
)

//go:embed specs.json
var embeddedConfig []byte

func main() {
	flag.Parse()
	logger := logging.New(os.Stderr, logrus.InfoLevel)

	if *devMode {
		logger.Infof("Running API server on port %d in dev mode", *port)
	} else {
		logger.Infof("Running API production server on port %d", *port)
		gin.SetMode(gin.ReleaseMode)
	}

	dbClient, err := db.NewClient()
	if err != nil {
		logger.Fatalf("Error initializing database client: %v", err)
	}

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()
	ctx = logging.AddToContext(ctx, logger)

	service, err := router.CreateNewService(ctx, *port, dbClient, embeddedConfig)
	if err != nil {
		logger.Fatalf("Error creating the server: %v", err)
	}
	err = service.Run(ctx)
	if err != nil {
		logrus.Fatalf("Error starting the server: %v", err)
	}
}
