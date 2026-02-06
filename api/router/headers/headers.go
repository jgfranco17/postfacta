package headers

import (
	"encoding/json"
	"net/http"

	"github.com/jgfranco17/postfacta/api/httperror"

	"github.com/gin-gonic/gin"
)

type OriginInfo struct {
	Origin  string `json:"origin"`
	Version string `json:"version"`
}

// Gets the origin info based on the received headers
func CreateOriginInfoHeader(c *gin.Context) (OriginInfo, error) {
	header := c.Request.Header["X-Origin-Info"]

	jsonHeader := OriginInfo{}

	if len(header) == 0 {
		return jsonHeader, httperror.New(c, http.StatusBadRequest, "X-Origin-Info header not found.")
	}

	if err := json.Unmarshal([]byte(header[0]), &jsonHeader); err != nil {
		return jsonHeader, httperror.New(c, http.StatusBadRequest, "Header schema validation: %s", err.Error())
	}
	return jsonHeader, nil
}
