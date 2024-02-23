package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// setupRouter setup http request router
func (adapter *RESTAdapter) setupRouter() {
	adapter.router.GET("/ping", adapter.handlePing)

	adapter.setupDeviceRouter()
	adapter.setupJobRouter()
}

func (adapter *RESTAdapter) handlePing(c *gin.Context) {
	logger := log.WithFields(log.Fields{
		"package":  "rest",
		"struct":   "RESTAdapter",
		"function": "handlePing",
	})

	logger.Infof("access request to %s", c.Request.URL)

	type pingOutput struct {
		Message string `json:"message"`
	}

	output := pingOutput{
		Message: "pong",
	}
	c.JSON(http.StatusOK, output)
}
