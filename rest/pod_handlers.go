package rest

import (
	"awds/types"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// setupDeviceRouter setup http request router for device
func (adapter *RESTAdapter) setupPodRouter() {
	// any devices can call these APIs
	adapter.router.GET("/pods",adapter.handleListPods)
	adapter.router.GET("/pods/:id", adapter.handleGetPod)
	adapter.router.POST("/pods", adapter.handleRegisterPod)
	adapter.router.PATCH("/pods/:id", adapter.handleUpdatePod)
	adapter.router.DELETE("/pods/:id", adapter.handleDeletePod)
}

func (adapter *RESTAdapter) handleListPods(c *gin.Context) {
	logger := log.WithFields(log.Fields{
		"package":  "rest",
		"struct":   "RESTAdapter",
		"function": "handleListPods",
	})

	logger.Infof("access request to %s", c.Request.URL)

	type listOutput struct {
		Pods []types.Pod `json:"pods"`
	}

	output := listOutput{}

	pods, err := adapter.logic.ListPods()
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	output.Pods = pods

	// success
	c.JSON(http.StatusOK, output)
}

func (adapter *RESTAdapter) handleGetPod(c *gin.Context) {
	logger := log.WithFields(log.Fields{
		"package":  "rest",
		"struct":   "RESTAdapter",
		"function": "handleGetPod",
	})

	logger.Infof("access request to %s", c.Request.URL)

	podID := c.Param("id")

	err := types.ValidatePodID(podID)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	device, err := adapter.logic.GetPod(podID)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// success
	c.JSON(http.StatusOK, device)
}

func (adapter *RESTAdapter) handleRegisterPod(c *gin.Context) {
	logger := log.WithFields(log.Fields{
		"package":  "rest",
		"struct":   "RESTAdapter",
		"function": "handleRegisterPod",
	})

	logger.Infof("access request to %s", c.Request.URL)

	type podRegistrationRequest struct {
		Endpoint 	string	`json:"endpoint"`
		Description string `json:"description,omitempty"`
	}

	var input podRegistrationRequest
	
	err := c.BindJSON(&input)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pod := types.Pod{
		ID:          	   types.NewPodID(),
		Endpoint:          input.Endpoint,
		Description: 	   input.Description, // optional
	}

	err = adapter.logic.RegisterPod(&pod)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, pod)
}

func (adapter *RESTAdapter) handleUpdatePod(c *gin.Context) {
	logger := log.WithFields(log.Fields{
		"package":  "rest",
		"struct":   "RESTAdapter",
		"function": "handleUpdatePod",
	})

	logger.Infof("access request to %s", c.Request.URL)

	podID := c.Param("id")

	err := types.ValidatePodID(podID)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}


	type podUpdateRequest struct {
		EndPoint          string `json:"endpoint"`
		Description 	  string `json:"description,omitempty"`
	}

	var input podUpdateRequest

	err = c.BindJSON(&input)
	fmt.Println(input)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(input.EndPoint) > 0 {
		// update IP
		err = adapter.logic.UpdatePodEndpoint(podID, input.EndPoint)
		if err != nil {
			// fail
			logger.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	if len(input.Description) > 0 {
		// update password
		err = adapter.logic.UpdatePodDescription(podID, input.Description)
		if err != nil {
			// fail
			logger.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	pod, err := adapter.logic.GetPod(podID)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, pod)
}

func (adapter *RESTAdapter) handleDeletePod(c *gin.Context) {
	logger := log.WithFields(log.Fields{
		"package":  "rest",
		"struct":   "RESTAdapter",
		"function": "handleDeletePod",
	})

	logger.Infof("access request to %s", c.Request.URL)

	podID := c.Param("id")

	err := types.ValidatePodID(podID)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pod, err := adapter.logic.GetPod(podID)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Debugf("Deleting Pod ID: %s", podID)

	err = adapter.logic.DeletePod(podID)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, pod)
}