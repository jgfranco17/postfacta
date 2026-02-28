package integration_test

import (
	"context"
	"io"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/jgfranco17/postfacta/api/db"
	"github.com/jgfranco17/postfacta/api/logging"
	"github.com/jgfranco17/postfacta/api/router"
	"github.com/sirupsen/logrus"
)

func NewMockServer(dbClient db.DatabaseClient) (*httptest.Server, error) {
	gin.SetMode(gin.TestMode)
	logger := logging.New(io.Discard, logrus.InfoLevel)
	ctx := logging.AddToContext(context.Background(), logger)

	const testMetadata string = `{"author":"test","repository":"","version":"0.0.0","license":"","languages":["Go"],"active":true}`
	service, err := router.CreateNewService(ctx, 0, dbClient, []byte(testMetadata))
	if err != nil {
		return nil, err
	}

	return httptest.NewServer(service.Router), nil
}
