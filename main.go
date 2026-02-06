package main

import (
	"flag"
	"time"

	"github.com/jgfranco17/postfacta/api/db"
	env "github.com/jgfranco17/postfacta/api/environment"
	"github.com/jgfranco17/postfacta/api/router"
	"github.com/jgfranco17/postfacta/api/router/system"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"

	_ "embed"
)

var (
	port    = flag.Int("port", 8080, "Port to listen on")
	devMode = flag.Bool("dev", true, "Run server in debug mode")
)

//go:embed specs.json
var embeddedConfig []byte

func init() {
	if env.IsLocalEnvironment() {
		logrus.SetFormatter(&logrus.TextFormatter{
			TimestampFormat: time.DateTime,
			DisableSorting:  true,
			PadLevelText:    true,
		})
		gin.SetMode(gin.DebugMode)
	} else {
		logrus.SetFormatter(&logrus.JSONFormatter{})
		gin.SetMode(gin.ReleaseMode)
	}
	prometheus.Register(system.HttpLastRequestReceivedTime)
	logrus.SetLevel(logrus.DebugLevel)
}

func main() {
	flag.Parse()
	if *devMode {
		logrus.Infof("Running API server on port %d in dev mode", *port)
	} else {
		logrus.Infof("Running API production server on port %d", *port)
		gin.SetMode(gin.ReleaseMode)
	}
	dbClient, err := db.NewClient()
	if err != nil {
		logrus.Fatalf("Error initializing database client: %v", err)
	}
	service, err := router.CreateNewService(*port, dbClient, embeddedConfig)
	if err != nil {
		logrus.Fatalf("Error creating the server: %v", err)
	}
	err = service.Run()
	if err != nil {
		logrus.Fatalf("Error starting the server: %v", err)
	}
}
