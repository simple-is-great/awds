// server.go: program running on edge server
// todo: make 3 modules(communication, scheduling, DB access)
// later: integrate with volume-service

package main

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	CPULowerBound = float64(2) // in Cores
	MemoryLowerBound = float64(2) // in GB
	
	SQLiteDBFileName = "awds.db"
)

type DeviceStatus struct {
	CPU float64 `json: "cpu"`
	Memory float64 `json: "memory"`
	TaskNum int `json: "taskNum"`
}

type AppRunStatus struct {
	CPU float64 `json: "cpu"`
	Memory float64 `json: "memory"`
	TaskNum int `json: "taskNum"`
}

type Workload struct { // may change later
	StartNum int 
	CurrentNum int
	EndNum	int
}

// for REST
func startRESTServer() {
	logger := log.WithFields(log.Fields{
		"package": "main",
		"function": "startRESTServer",
	})

	router := gin.Default()
	router.Use(cors.New(
		cors.Config{
			AllowOrigins: []string{"http://localhost:5173", "http://155.230.36.126"},
			AllowMethods: []string{"POST", "GET"},
			AllowHeaders: []string{"Content-Type", "application/json"},
			MaxAge: 24 * time.Hour,
		},
	))

	// set up router, change later
	router.POST("/awdp", receiveAppRunStatus)
	
	router.Run("155.230.36.27:3015")

	logger.Infoln("communication Module Started...")

}

func receiveDevicesStatus(c *gin.Context) error {
	logger := log.WithFields(log.Fields{
		"package":  "main",
		"function": "sendDeviceStatus",
	})

	// receive `{"cpu": %f, "mem", %f, "taskNum": %d}` as body
	// call scheduler
	// return "input": "output": in JSON format

	var deviceStatus DeviceStatus
	var 
	err := c.BindJSON(&status)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		logger.Error(err)
		return err
	}
	
	reschedule(status)

	return nil
}


// for scheduler
func reschedule(podStatus AppRunStatus, deviceStatus DeviceStatus)  {
	logger := log.WithFields(log.Fields{
		"package": "main",
		"function": "reschedule",
	})

	logger.Infoln("Scheduling Module Started...")
	
	// TODO: make two functions for scheduling
	// need to specify condition
	// check condition and call offloading function
	if (deviceStatus.CPU < CPULowerBound) && (deviceStatus.Memory < MemoryLowerBound) {
		// fully offload
	}
	
	// add more job to server

	// other condition
	// cannot offload

	return 
}

func fullyOffload(status DeviceStatus) {
	// stop job 
	// fully offload to edge server
	
}

func partiallyOffload(status DeviceStatus) {
	// add more job to the pod

}

// for DB
func startAppRunStatusDB() (*gorm.DB, error) {
	logger := log.WithFields(log.Fields{
		"package": "main",
		"function": "startAppRunStatusDB",
	})

	logger.Infoln("DBAdapter Module Started...")

	db, err := gorm.Open(sqlite.Open(SQLiteDBFileName), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(Workload{})
	if err != nil {
		return nil, err
	}

	return db, nil
	
}


// TODO: make two functions(INSERT and SELECT)
// INSERT and SELECT * FROM DB
// match DB 

func (db *gorm.DB) getWorkload(deviceID string, appRunID string) (Workload, error) {
	var workload Workload
	result := db.Where("apprun_id = ?", appRunID).First(&workload)
	if result.Error != nil {
		return Workload{}, result.Error
	}
	

	return workload, nil
}

func (db *gorm.DB) InsertAppRunStatus(status AppRunStatus) error {
	result := db.Create(&status)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected != 1 {
		return xerrors.Errorf("failed to insert an appRunStatus")
	}

	return nil
}


func main() {
	receiveAppRunStatus()
	dbAdapter, err := startAppRunStatusDB()


}

// func waitForCtrlC() {
// 	var endWaiter sync.WaitGroup

// 	endWaiter.Add(1)
// 	signalChannel := make(chan os.Signal, 1)

// 	signal.Notify(signalChannel, os.Interrupt)

// 	go func() {
// 		<-signalChannel
// 		endWaiter.Done()
// 	}()

// 	endWaiter.Wait()
// }