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
		DeviceID string 	  `json:"device_id"`
		PodID string 		  `json:"pod_id"`
		InputSize int		  `json:"input_size"`
	}

	var input jobCreationRequest
	
	// need to get batch time from func
	// do not manually get input

	err := c.BindJSON(&input)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	job := types.Job{
		ID:          	types.NewJobID(),
		DeviceID:       input.DeviceID,
		PodID:  		input.PodID,
		InputSize: 		input.InputSize,
	}
	
	// TODO: need to notify finished job
	err = adapter.logic.CreateJob(&job)
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
		DeviceID			string 	  		`json:"device_id"`
		PodID 				string 		  	`json:"pod_id"`
		InputSize 			int		  		`json:"input_size"`
		PartitionRate		float64			`json:"partition_rate"`
		Completed			bool			`json:"completed"`
		DeviceStartIndex 	int				`json:"device_start_index"` 
		DeviceEndIndex 		int				`json:"device_end_index"` 
		PodStartIndex 		int				`json:"pod_start_index"` 
		PodEndIndex 		int				`json:"pod_end_index"` 
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

	if len(input.DeviceID) > 0 {
		// update DeviceID
		err = adapter.logic.UpdateJobDevice(jobID, input.DeviceID)
		if err != nil {
			// fail
			logger.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	if len(input.PodID) > 0 {
		// update PodID
		err = adapter.logic.UpdateJobPod(jobID, input.PodID)
		if err != nil {
			// fail
			logger.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	if input.InputSize > 0 {
		// update InputSize
		err = adapter.logic.UpdateInputSize(jobID, input.InputSize)
		if err != nil {
			// fail
			logger.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	if input.PartitionRate > 0 {
		// update PartitionRate
		err = adapter.logic.UpdatePartitionRate(jobID, input.PartitionRate)
		if err != nil {
			// fail
			logger.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	if !input.Completed {
		// unset Completed as false
		err = adapter.logic.UpdateJobCompleted(jobID, input.Completed)
		if err != nil {
			// fail
			logger.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	if input.DeviceStartIndex > 0 {
		// update DeviceStartIndex
		err = adapter.logic.UpdateDeviceStartIndex(jobID, input.DeviceStartIndex)
		if err != nil {
			// fail
			logger.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	if input.DeviceEndIndex > 0 {
		// update DockerImage
		err = adapter.logic.UpdateDeviceEndIndex(jobID, input.DeviceEndIndex)
		if err != nil {
			// fail
			logger.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	if input.PodStartIndex > 0 {
		// update PodStartIndex
		err = adapter.logic.UpdatePodStartIndex(jobID, input.PodStartIndex)
		if err != nil {
			// fail
			logger.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	if input.PodEndIndex > 0 {
		// update PodEndIndex
		err = adapter.logic.UpdatePodEndIndex(jobID, input.PodEndIndex)
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
