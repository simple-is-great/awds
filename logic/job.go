package logic

import (
	"awds/types"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
)

var (
	// batchSize int = 50
	batchInit int = 50 // batchSize used for first batch
	averageInputSize = 1000 // 10e5B -> 100 KB
	// failedJob map[string][]string
)

// map to record device info for calculating batch size
type deviceRecord map[string][]float64

func (logic *Logic) ListJobs() ([]types.Job, error) {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "ListJobs",
	})

	logger.Debug("received ListJobs()")

	return logic.dbAdapter.ListJobs()
}

func (logic *Logic) GetJob(jobID string) (types.Job, error) {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "GetJob",
	})

	logger.Debug("received GetJob()")

	return logic.dbAdapter.GetJob(jobID)
}

func (logic *Logic) InsertJob(job *types.Job) error {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "InsertJob",
	})

	logger.Debug("received InsertJob()")

	return logic.dbAdapter.InsertJob(job)
}

func (logic *Logic) UpdateDeviceIDList(jobID string, deviceIDList []string) error {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "UpdateDeviceIDList",
	})

	deviceIDTemp := []string{}
	for _, deviceID := range deviceIDList{
		deviceIDTemp = append(deviceIDTemp, deviceID)
	}

	deviceIDListCSV := strings.Join(deviceIDTemp, ",") // make deviceID list into string

	logger.Debug("received UpdateJobDevice()")

	return logic.dbAdapter.UpdateDeviceIDList(jobID, deviceIDListCSV)
}

func (logic *Logic) UpdateEndIndex(jobID string, endIndex int) error {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "UpdateEndIndex",
	})

	logger.Debug("received UpdateEndIndex()")

	return logic.dbAdapter.UpdateEndIndex(jobID, endIndex)
}

func (logic *Logic) DeleteJob(jobID string) error {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "DeleteJob",
	})

	logger.Debug("received DeleteJob()")

	return logic.dbAdapter.DeleteJob(jobID)
}

func (logic *Logic) AdjustBatchSize(devIdQ *Queue, devRcdMap *deviceRecord, startIdx int, endIdx int) (int, error) {
	logger := log.WithFields(log.Fields{
		"package": "logic",
		"struct" : "Logic",
		"function" : "AdjustBatchSize",
	})
	logger.Debug("received AdjustBatchSize()")

	
	deviceNum := len(*devIdQ)
	// precomputation: 1% / deviceNum, serve as endIdx for precomputation
	adjustBatchSize := int(0.01 * float64(endIdx - startIdx - 1) / float64(deviceNum))
	// set minimum size to 1
	if adjustBatchSize < 1 {
		adjustBatchSize = 1 
	}
	
	var wg sync.WaitGroup
	wg.Add(len(*devIdQ))
	errChan := make(chan error, 1)

	for idx, deviceID := range *devIdQ{
		device, err := logic.dbAdapter.GetDevice(deviceID)
		if err != nil {
			return 0, err
		}

		func(i int)() {
			defer wg.Done()
			elapsedTime, err := logic.Compute(&device, startIdx, startIdx + adjustBatchSize)
			if err != nil {
				errChan <- err
				return 
			}
			device, err := logic.GetDeviceResourceMetrics(&device)
			if err != nil {
				errChan <- err
				return
			}
			// SetNextBatchSize predictTime float64, elapsedTime float64, batchSize float64, adjustBatchSize int, batchNum int
			// Predict elapsedTime float64, batchSize float64, batchSize float64, batchNum int
			// predictTime, elapsedTime, batchSize, batchSize, batchNum
			nextBatchSize := logic.SetNextBatchSize(float64(0), elapsedTime, device.Memory, float64(adjustBatchSize), adjustBatchSize, 1)
			predictTime := logic.Predict(elapsedTime, float64(adjustBatchSize), float64(nextBatchSize), 1)
			fmt.Println("AdjustBatchSize(deviceID, elapsedTime, adjustBatchSize, nextBatchSize, predictTime): ", device.ID, elapsedTime, adjustBatchSize, nextBatchSize, predictTime)
			(*devRcdMap)[deviceID][0] = predictTime // current PredictTime
			(*devRcdMap)[deviceID][1] = elapsedTime // current elapsedTime
			(*devRcdMap)[deviceID][2] = float64(adjustBatchSize) // current BatchSize
			(*devRcdMap)[deviceID][3] = float64(nextBatchSize)// nextBatchSize
			(*devRcdMap)[deviceID][4] += 1 // used to count batchNumber
			return
			
		}(idx)
		startIdx += adjustBatchSize // update StartIdx
	}

	wg.Wait()

	return adjustBatchSize, nil
}

func (logic *Logic) SetNextBatchSize(predictTime float64, elapsedTime float64, availableMemory float64, batchSize float64, adjustBatchSize int, batchNum int) int {
	// set nextBatchSize based on predictTime, elapsedTime, batchSize of current batchSize
	switch (batchNum){
	case 0:
		return adjustBatchSize
	case 1:
		return batchInit
	default:
		nextBatch := int(predictTime / elapsedTime * batchSize)
		if nextBatch < 1{
			nextBatch = 1
		} else if int(averageInputSize * nextBatch) > int(0.8 * availableMemory){
			nextBatch = int(0.8 * availableMemory / float64(averageInputSize))
		}
		return nextBatch
	}
}


func (logic *Logic) Predict(elapsedTime float64, batchSize float64, nextBatchSize float64, batchNum int) float64{
	// predict based on elapsedTime, batchSize, nextBatchSize
	if batchNum == 0 {
		return 0
	} else {
		return (elapsedTime/ batchSize) * nextBatchSize
	}
}

func (logic *Logic) Compute(device *types.Device, batchStartIdx int, batchEndIdx int) (float64, error) {
	logger := log.WithFields(log.Fields{
		"package": "logic",
		"struct" : "Logic",
		"function" : "Compute",
	})

	logger.Debugf("received Compute()")

	type Response struct {
		Result	float64	`json: "result"`
	}
	var response Response
	fullEndpoint := logic.GetFullEndpoint(device.IP, device.Port,device.Endpoint, batchStartIdx, batchEndIdx)
	
	client := resty.New()
	_, err := client.R().SetResult(&response).Get(fullEndpoint)
	if err != nil {
		return -1, err
	}
	
	elapsedTime := response.Result

	return elapsedTime, nil
}


// TODO: implement map for each Job Object and 
// func (logic *Logic) SaveFailedWorkload(jobID string, startIndex int, batchSize int) error {
// 	// ensure failedJob map is initialized
// 	if failedJob == nil {
// 		failedJob = make(map[string][]string)
// 	}

// 	var failedIndex string = fmt.Sprintf("%d-%d", startIndex, startIndex+batchSize)
// 	failedJob[jobID] = append(failedJob[jobID], failedIndex)
	
// 	return nil
// }

// func(logic *Logic) ComputeWithRetry(jobID string, deviceID string, batchSize int, maxRetries int) error {
// 	for i := 0; i < maxRetries; i++ {
// 		device, err := logic.dbAdapter.GetDevice(deviceID)
// 		if err != nil {
// 			return err
// 		}
		
// 		job, err := logic.dbAdapter.GetJob(jobID)
// 		if err != nil {
// 			return err
// 		}

// 		err = logic.Compute(&device, &job, batchSize)
// 		if err == nil {
// 			// success, no need to retry
// 			return nil 
// 		}	
// 	}
// 	return fmt.Errorf("compute failed, saved to failedJob and will be computed later")
// }


func (logic *Logic) ScheduleJob(jobID string) error {
	logger := log.WithFields(log.Fields{
		"package": "logic",
		"struct" : "Logic",
		"function" : "ScheduleJob",
	})
	
	logger.Debug("received ScheduleJob()")
	
	// filename := time.Now().Format("01-02-15:04:05")
	// file, err := os.Create(filename + ".txt")
	// if err != nil {
	// 	fmt.Println("Error create file:", err)
	// 	os.Exit(1)
	// }
	// defer file.Close()

	// os.Stdout = file

	startTime := time.Now()
	job, err := logic.dbAdapter.GetJob(jobID)
	if err != nil {
		return err
	}
	// start index for batch, increases from job.StartIndex to job.EndIndex	
	jobStartIdx := job.StartIndex // start index for job
	jobEndIdx := job.EndIndex // end index for job
	batchStartIdx := jobStartIdx
	var batchEndIdx int
	

	errChan := make(chan error, 1)
	// queue to hold available deviceID
	var deviceIDQueue Queue
	// map to hold previous and current latency
	deviceRecordMap := deviceRecord{}
	for _, deviceID := range job.DeviceIDList{
		deviceIDQueue.Enqueue(deviceID)
		deviceRecordMap[deviceID] = make([]float64, 5)
	}
	
	// adjustBatchSize for entire time
	adjustBatchSize, err := logic.AdjustBatchSize(&deviceIDQueue, &deviceRecordMap, jobStartIdx, jobEndIdx)
		if err != nil {
			return err
	}

	// for _, deviceID := range job.DeviceIDList{
	// 	fmt.Printf("after adjustment: ")
	// 	fmt.Println(deviceID, deviceRecordMap[deviceID])
	// }

	jobStartIdx += adjustBatchSize * len(job.DeviceIDList) // update precompute results
	batchEndIdx = jobStartIdx

	for {
		// get deviceID from queue
		// TODO: busy waiting -> change later
		if len(deviceIDQueue) == 0 {
			continue
		}
		
		dID, err := deviceIDQueue.Dequeue()
		if err != nil {
			fmt.Println(err)
		}

		device, err := logic.dbAdapter.GetDevice(dID)
		if err != nil {
			deviceIDQueue.Enqueue(dID)
			return err
		}

		if (batchStartIdx > jobEndIdx) || (batchEndIdx == jobEndIdx) {
			break // job done
		}

		batchSize := int(deviceRecordMap[dID][3]) // call batchSize from deviceRecord

		// swap batchEndIdx and batchStartIdx
		temp := batchEndIdx
		batchStartIdx := batchEndIdx
		batchEndIdx = temp + batchSize

		// batchEndIdx cannot exceed jobEndIdx
		if batchEndIdx > jobEndIdx {
			batchEndIdx = jobEndIdx
		}

		deviceRecordMap[dID][2] = deviceRecordMap[dID][3] // update nextbatchSize to currentBatchSize
		deviceRecordMap[dID][4] += 1 // update batchNum

		// batchStartIdx  < jobEndIndex -> create thread
		if batchStartIdx < jobEndIdx{
			// func inside thread
			go func(){
				// send compute request to device
				elapsedTime, err := logic.Compute(&device, batchStartIdx, batchEndIdx)
				if err != nil {
					errChan <- err
					return 
				}

				// update resource metric of device when batch ends
				device, err := logic.GetDeviceResourceMetrics(&device)
				if err != nil {
					errChan <- err
					return
				}
				
				nextBatchSize := logic.SetNextBatchSize(deviceRecordMap[dID][0], deviceRecordMap[dID][1], device.Memory, deviceRecordMap[dID][2], 0, int(deviceRecordMap[dID][4]))
				predictTime := logic.Predict(deviceRecordMap[dID][1], deviceRecordMap[dID][2], float64(nextBatchSize), int(deviceRecordMap[dID][4]))
				fmt.Println("In Schedule loop(deviceID, elapsedTime, batchSize, nextBatchSize, predictTime): ", dID, elapsedTime, batchSize, nextBatchSize, predictTime)
				deviceRecordMap[dID][0] = predictTime // predictTime
				deviceRecordMap[dID][1] = elapsedTime // elapsedTime
				deviceRecordMap[dID][2] = float64(batchSize) // current batch Size
				deviceRecordMap[dID][3] = float64(nextBatchSize) // next batch Size
				
				// err = logic.dbAdapter.UpdateDeviceResourceMetrics(dID, device.Memory, device.NetworkLatency)
				// if err != nil {
				// 	errChan <- err
				// 	return
				// }

				// compute succeed -> enqueue deviceID to get another batch
				deviceIDQueue.Enqueue(dID)
			}()
		}

		// fmt.Println("batchStartIdx & batchEndIdx", batchStartIdx, batchEndIdx)
		// job ended
		// if (batchEndIdx == jobEndIdx) || (batchStartIdx >= jobEndIdx) {
		// 	fmt.Println("Job Done!")
		// 	break
		// }

	}
	// TODO: handle resources after goroutine finishes
	// fmt.Println("startIdx after scheduling", batchStartIdx)
	// update startIdx
	err = logic.dbAdapter.UpdateStartIndex(jobID, batchStartIdx)
	if err != nil {
        return err
    }

	// update Completed to true
    err = logic.dbAdapter.UpdateJobCompleted(jobID, true)
    if err != nil {
        return err
    }
	timeTaken := time.Since(startTime)
	fmt.Println("time taken:", timeTaken)
	return nil
}
		
// func (logic *Logic) UnscheduleJob(jobID string) error {
// 	logger := log.WithFields(log.Fields{
// 		"package": "logic",
// 		"struct" : "Logic",
// 		"function" : "UnscheduleJob",
// 	})

// 	logger.Debug("received UnscheduleJob()")

// 	job, err := logic.dbAdapter.GetJob(jobID)
// 	if err != nil {
// 		return err
// 	}

// 	pod, err := logic.dbAdapter.GetPod(job.PodID)
// 	if err != nil {
// 		return err
// 	}

// 	for true {
// 		job, _ = logic.dbAdapter.GetJob(jobID)
// 		Compute()
		
// 		// break when job completes
// 		if job.Completed /* &&  device Job and pod Job completes  */ {
// 			break
// 		}
// 	}
// }