package v0

import (
	"net/http"

	"github.com/jgfranco17/postfacta/api/db"
	"github.com/jgfranco17/postfacta/api/httperror"

	"github.com/gin-gonic/gin"
)

func getAllIncidents(dbClient db.DatabaseClient) func(c *gin.Context) error {
	return func(c *gin.Context) error {
		incidents, err := dbClient.GetAllIncidents(c)
		if err != nil {
			return httperror.New(c, http.StatusInternalServerError, "")
		}
		c.JSON(http.StatusOK, incidents)
		return nil
	}
}
