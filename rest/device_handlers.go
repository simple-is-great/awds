package rest

import (
	"awds/types"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// setupDeviceRouter setup http request router for device
func (adapter *RESTAdapter) setupDeviceRouter() {
	// any devices can call these APIs
	adapter.router.GET("/devices", adapter.handleListDevices)
	adapter.router.GET("/devices/:id", adapter.handleGetDevice)
	adapter.router.POST("/devices", adapter.handleRegisterDevice)
	adapter.router.PATCH("/devices/:id", adapter.handleUpdateDevice)
	adapter.router.DELETE("/devices/:id", adapter.handleDeleteDevice)
}

func (adapter *RESTAdapter) handleListDevices(c *gin.Context) {
	logger := log.WithFields(log.Fields{
		"package":  "rest",
		"struct":   "RESTAdapter",
		"function": "handleListDevices",
	})

	logger.Infof("access request to %s", c.Request.URL)

	type listOutput struct {
		Devices []types.Device `json:"devices"`
	}

	output := listOutput{}
	
	devices, err := adapter.logic.ListDevices()
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

		output.Devices = devices
	

	// success
	c.JSON(http.StatusOK, output)
}

func (adapter *RESTAdapter) handleGetDevice(c *gin.Context) {
	logger := log.WithFields(log.Fields{
		"package":  "rest",
		"struct":   "RESTAdapter",
		"function": "handleGetDevice",
	})

	logger.Infof("access request to %s", c.Request.URL)

	deviceID := c.Param("id")

	err := types.ValidateDeviceID(deviceID)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	device, err := adapter.logic.GetDevice(deviceID)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// success
	c.JSON(http.StatusOK, device)
}

func (adapter *RESTAdapter) handleRegisterDevice(c *gin.Context) {
	logger := log.WithFields(log.Fields{
		"package":  "rest",
		"struct":   "RESTAdapter",
		"function": "handleRegisterDevice",
	})

	logger.Infof("access request to %s", c.Request.URL)

	type deviceRegistrationRequest struct {
		IP						string			`json:"ip"`
		Port					string			`json:"port"`
		Endpoint				string			`json:"endpoint"`
		Description 			string    		`json:"description,omitempty"`
	}

	var input deviceRegistrationRequest
	
	err := c.BindJSON(&input)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(input.IP) == 0 {
		remoteAddrFields := strings.Split(c.Request.RemoteAddr, ":")
		if len(remoteAddrFields) > 0 {
			input.IP = remoteAddrFields[0]
		}
	}

	if len(input.Port) == 0 {
		remoteAddrFields := strings.Split(c.Request.RemoteAddr, ":")
		if len(remoteAddrFields) > 0 {
			input.Port = remoteAddrFields[1]
		}
	}

	device := types.Device{
		ID:          types.NewDeviceID(),
		IP:			 input.IP,
		Port: 		 input.Port,
		Endpoint: 	 input.Endpoint,
		Description: input.Description, // optional
	}

	err = adapter.logic.CreateDevice(&device)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, device)
}

func (adapter *RESTAdapter) handleUpdateDevice(c *gin.Context) {
	logger := log.WithFields(log.Fields{
		"package":  "rest",
		"struct":   "RESTAdapter",
		"function": "handleUpdateDevice",
	})

	logger.Infof("access request to %s", c.Request.URL)
	deviceID := c.Param("id")

	err := types.ValidateDeviceID(deviceID)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	type deviceUpdateRequest struct {
		Endpoint          string `json:"endpoint"`
		Description 	  string `json:"description,omitempty"`
	}

	var input deviceUpdateRequest

	err = c.BindJSON(&input)
	fmt.Println(input)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(input.Endpoint) > 0 {
		// update IP
		err = adapter.logic.UpdateDeviceEndpoint(deviceID, input.Endpoint)
		if err != nil {
			// fail
			logger.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	if len(input.Description) > 0 {
		// update password
		err = adapter.logic.UpdateDeviceDescription(deviceID, input.Description)
		if err != nil {
			// fail
			logger.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	device, err := adapter.logic.GetDevice(deviceID)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, device)
}

func (adapter *RESTAdapter) handleDeleteDevice(c *gin.Context) {
	logger := log.WithFields(log.Fields{
		"package":  "rest",
		"struct":   "RESTAdapter",
		"function": "handleDeleteDevice",
	})

	logger.Infof("access request to %s", c.Request.URL)

	deviceID := c.Param("id")

	err := types.ValidateDeviceID(deviceID)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	device, err := adapter.logic.GetDevice(deviceID)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Debugf("Deleting Device ID: %s", deviceID)

	err = adapter.logic.DeleteDevice(deviceID)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, device)
}

