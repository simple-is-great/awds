package logic

import (
	"awds/types"
	"fmt"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
)

var (
	// starts small, grow later
	batchSize int = 10
	// failedJob map[string][]string
)

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

func (logic *Logic) Precompute(devIdQ *Queue, devRcdMap deviceRecord, startIdx int, endIdx int) (int, error) {
	logger := log.WithFields(log.Fields{
		"package": "logic",
		"struct" : "Logic",
		"function" : "Precompute",
	})
	logger.Debug("received Precompute()")

	deviceNum := len(*devIdQ)
	// deviceLatencyList := make([]float64, deviceNum)
	// precomputation: 1% / deviceNum, serve as endIdx for precomputation
	precomputeSize := int(0.01 * float64(endIdx - startIdx - 1) / float64(deviceNum))
	// set minimum size to 1
	if precomputeSize < 1 {
		precomputeSize = 1 
	}

	errChan := make(chan error, 1)

	for idx, deviceID := range *devIdQ{
		device, err := logic.dbAdapter.GetDevice(deviceID)
		if err != nil {
			return -1, err
		}

		go func(i int)() {
			precomputeLatency, err := logic.Compute(&device, startIdx, startIdx + precomputeSize)
			if err != nil {
				errChan <- err
				return 
			}
			device, err := logic.GetDeviceResourceMetrics(&device)
			if err != nil {
				errChan <- err
				return
			}

			predictLatency := float64(batchSize / precomputeSize) * (precomputeLatency - float64(precomputeSize)/device.NetworkLatency)

			devRcdMap[deviceID][0] = float64(precomputeSize) // precomputeSize in float64
			devRcdMap[deviceID][1] = device.NetworkLatency // networkLatency
			devRcdMap[deviceID][2] = precomputeLatency // precomputeLatency
			devRcdMap[deviceID][3] = predictLatency
		}(idx)
		startIdx += precomputeSize // update StartIdx


	}
	
	return precomputeSize * deviceNum, nil
}

func (logic *Logic) Compute(device *types.Device, batchStartIdx int, batchEndIdx int) (float64, error) {
	logger := log.WithFields(log.Fields{
		"package": "logic",
		"struct" : "Logic",
		"function" : "Compute",
	})

	logger.Debugf("received Compute()")

	var response map[string]interface{}
	fullEndpoint := logic.GetFullEndpoint(device.IP, device.Port,device.Endpoint, batchStartIdx, batchEndIdx)
	
	client := resty.New()
	_, err := client.R().SetResult(&response).Get(fullEndpoint)
	if err != nil {
		return -1, err
	}
	
	computeLatency := response["result"].(float64)

	return computeLatency, nil
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
	
	startTime := time.Now()
	job, err := logic.dbAdapter.GetJob(jobID)
	if err != nil {
		return err
	}
	// start index for batch, increases from job.StartIndex to job.EndIndex	
	batchStartIdx := job.StartIndex
	batchEndIdx := 0 // end index for batch
	jobStartIdx := job.StartIndex // start index for job
	jobEndIdx := job.EndIndex // end index for job

	// err channel for goroutines
	errChan := make(chan error, 1)
	// map to hold previous and current latency
	// queue to hold available deviceID
	var deviceIDQueue Queue
	var deviceRecordMap deviceRecord
	for _, deviceID := range job.DeviceIDList{
		deviceIDQueue.Enqueue(deviceID)
		deviceRecordMap[deviceID] = make([]float64, 4)
	}
	
	// 전체 device 동시에 수행
	precomputeSize, err := logic.Precompute(&deviceIDQueue, deviceRecordMap, jobStartIdx, jobEndIdx)
		if err != nil {
			return err
	}

	jobStartIdx += precomputeSize // update precompute results

	for {
		// TODO: we need to calculate batchSize and pass into threads
		// get deviceID from queue
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

		// not first time -> modify startIdx then move on
		if batchStartIdx != jobStartIdx{
			batchStartIdx += batchSize
			batchEndIdx = batchStartIdx + batchSize // TODO: need to get batchSize from precomputation results

			if batchStartIdx + batchSize > jobEndIdx{
				batchEndIdx = jobEndIdx
			}
		}

		// batchStartIdx  < jobEndIndex -> create thread
		if batchStartIdx < jobEndIdx{
			// func inside thread
			go func(){
				// send compute request to device
				currentLatency, err := logic.Compute(&device, batchStartIdx, batchEndIdx)
				if err != nil {
					errChan <- err
					return 
				}
				// set current latency for computation
				fmt.Println(dID, "current latency", currentLatency, "s")
				// // update resource metric of device when batch ends
				fmt.Println("before getting metrics ", device.Memory, device.NetworkLatency)
				device, err := logic.GetDeviceResourceMetrics(&device)
				if err != nil {
					errChan <- err
					return
				}

				predictLatency := float64(float64(precomputeSize)/deviceRecordMap[dID][3]) * float64(batchSize)

				deviceRecordMap[dID][0] = float64(batchSize) // precomputeSize in float64
				deviceRecordMap[dID][1] = device.NetworkLatency // networkLatency
				deviceRecordMap[dID][2] = currentLatency // precomputeLatency
				deviceRecordMap[dID][3] = predictLatency // predict next batch's latency
				
				// err = logic.dbAdapter.UpdateDeviceResourceMetrics(dID, device.Memory, device.NetworkLatency)
				// if err != nil {
				// 	errChan <- err
				// 	return
				// }
				fmt.Println("after getting metrics", device.Memory, device.NetworkLatency)
				// compute succeed -> enqueue deviceID to get another batch
				deviceIDQueue.Enqueue(dID)
			}()
		}

		// first time -> modify startIdx after compute
		// if batchStartIdx == job.StartIndex{
		// 	batchStartIdx += batchSize
		// 	batchEndIdx = batchStartIdx + batchSize // TODO: need to get batchSize from precomputation results
		
		// 	if batchStartIdx + batchSize > jobEndIdx{
		// 		batchEndIdx = jobEndIdx 
		// 	}
		// }
		
		fmt.Println("batchStartIdx & batchEndIdx", batchStartIdx, batchEndIdx)
		// job ended
		if (batchEndIdx == jobEndIdx) || (batchStartIdx >= jobEndIdx) {
			fmt.Println("Job Done!")
			break
		}
	}

	// TODO: handle resources after goroutine finishes
	// TODO: handle errors if goroutine fails
	
	fmt.Println("startIdx after scheduling", batchStartIdx)
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

func (logic *Logic) CalcBatchSize(deviceID string, devRcdMap deviceRecord)(int, error){
	// 1. 예측
	// 2. 다음 배치 계산




	return batchSize, nil
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