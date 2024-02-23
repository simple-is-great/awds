package rest

import (
	"awds/types"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// setupAppRouter setup http request router for app
func (adapter *RESTAdapter) setupJobRouter() {
	// any devices can call these APIs
	adapter.router.GET("/jobs", adapter.handleListJobs)
	adapter.router.GET("/jobs/:id", adapter.handleGetJob)
	adapter.router.POST("/jobs", adapter.handleCreateJob)
	adapter.router.PATCH("/jobs/:id", adapter.handleUpdateJob)
	adapter.router.DELETE("/jobs/:id", adapter.handleDeleteJob)

	adapter.router.POST("/schedules/:id", adapter.handleScheduleJob)
	// adapter.router.DELETE("/schedules/:id", adapter.handleUnscheduleJob)

}

func (adapter *RESTAdapter) handleListJobs(c *gin.Context) {
	logger := log.WithFields(log.Fields{
		"package":  "rest",
		"struct":   "RESTAdapter",
		"function": "handleListJobs",
	})

	logger.Infof("access request to %s", c.Request.URL)

	type listOutput struct {
		Jobs []types.Job `json:"jobs"`
	}

	output := listOutput{}

	jobs, err := adapter.logic.ListJobs()
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	output.Jobs = jobs

	// success
	c.JSON(http.StatusOK, output)
}

func (adapter *RESTAdapter) handleGetJob(c *gin.Context) {
	logger := log.WithFields(log.Fields{
		"package":  "rest",
		"struct":   "RESTAdapter",
		"function": "handleGetJob",
	})

	logger.Infof("access request to %s", c.Request.URL)

	jobID := c.Param("id")

	err := types.ValidateJobID(jobID)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	app, err := adapter.logic.GetJob(jobID)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// success
	c.JSON(http.StatusOK, app)
}

func (adapter *RESTAdapter) handleCreateJob(c *gin.Context) {
	logger := log.WithFields(log.Fields{
		"package":  "rest",
		"struct":   "RESTAdapter",
		"function": "handleCreateJob",
	})

	logger.Infof("access request to %s", c.Request.URL)


	type jobCreationRequest struct {
		DeviceIDList []string 			`json:"device_id_list"`
		EndIndex 	 int				`json:"end_index"`
	}

	var input jobCreationRequest

	err := c.BindJSON(&input)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	job := types.Job{
		ID:          	types.NewJobID(),
		DeviceIDList:   input.DeviceIDList,
		EndIndex: 		input.EndIndex,
	}
	
	err = adapter.logic.InsertJob(&job)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, job)
}

func (adapter *RESTAdapter) handleUpdateJob(c *gin.Context) {
	logger := log.WithFields(log.Fields{
		"package":  "rest",
		"struct":   "RESTAdapter",
		"function": "handleUpdateJob",
	})

	logger.Infof("access request to %s", c.Request.URL)

	jobID := c.Param("id")

	err := types.ValidateJobID(jobID)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	type jobUpdateRequest struct {
		// TODO: decide how to change deviceID list in the job
		DeviceIDList	[]string 	  	`json:"device_id_list"`
		EndIndex		int				`json:"end_index"`
	}

	var input jobUpdateRequest

	err = c.BindJSON(&input)
	fmt.Println(input)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(input.DeviceIDList) > 0 {
		// update DeviceID
		err = adapter.logic.UpdateDeviceIDList(jobID, input.DeviceIDList)
		if err != nil {
			// fail
			logger.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	if input.EndIndex > 0 {
		// update InputSize
		err = adapter.logic.UpdateEndIndex(jobID, input.EndIndex)
		if err != nil {
			// fail
			logger.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	job, err := adapter.logic.GetJob(jobID)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, job)
}

func (adapter *RESTAdapter) handleDeleteJob(c *gin.Context) {
	logger := log.WithFields(log.Fields{
		"package":  "rest",
		"struct":   "RESTAdapter",
		"function": "handleDeleteJob",
	})

	logger.Infof("access request to %s", c.Request.URL)

	jobID := c.Param("id")

	err := types.ValidateJobID(jobID)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	job, err := adapter.logic.GetJob(jobID)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Debugf("Deleting Job ID: %s", jobID)

	err = adapter.logic.DeleteJob(jobID)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, job)
}

func (adapter *RESTAdapter) handleScheduleJob(c *gin.Context) {
	logger := log.WithFields(log.Fields{
		"package":  "rest",
		"struct":   "RESTAdapter",
		"function": "handleScheduleJob",
	})

	logger.Infof("access request to %s", c.Request.URL)
	
	jobID := c.Param("id")

	err := types.ValidateJobID(jobID)
	if err != nil {
		logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	job, err := adapter.logic.GetJob(jobID)
	if err != nil {
		logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = adapter.logic.ScheduleJob(jobID)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logger.Debugf("Scheduling Job ID: %s", jobID)

	c.JSON(http.StatusOK, job)
}

// func (adapter *RESTAdapter) handleUnscheduleJob(c *gin.Context) {
// 	logger := log.WithFields(log.Fields{
// 		"package":  "rest",
// 		"struct":   "RESTAdapter",
// 		"function": "handleUnmountVolume",
// 	})

// 	logger.Infof("access request to %s", c.Request.URL)

// 	jobID := c.Param("id")

// 	err := types.ValidateJobID(jobID)
// 	if err != nil {
// 		// fail
// 		logger.Error(err)
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	job, err := adapter.logic.GetJob(jobID)
// 	if err != nil {
// 		// fail
// 		logger.Error(err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	logger.Debugf("Unscheduling Job ID: %s", jobID)

// 	err = adapter.logic.UnscheduleJob(jobID)
// 	if err != nil {
// 		// fail
// 		logger.Error(err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, job)
// }
