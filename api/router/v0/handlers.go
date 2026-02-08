package v0

import (
	"net/http"

	"github.com/jgfranco17/postfacta/api/core"
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

type incidentStartResponse struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

func startIncident(dbClient db.DatabaseClient) func(c *gin.Context) error {
	return func(c *gin.Context) error {
		var requestBody core.IncidentRequest
		if err := c.ShouldBindJSON(&requestBody); err != nil {
			return httperror.New(c, http.StatusBadRequest, "Invalid request body: %s", err.Error())
		}

		incident := core.NewIncident(requestBody)
		if err := dbClient.StoreIncident(c, incident); err != nil {
			return httperror.New(c, http.StatusInternalServerError, "")
		}

		c.JSON(http.StatusCreated, incidentStartResponse{
			ID:      incident.ID,
			Message: "Incident started successfully",
		})
		return nil
	}
}
