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
	batchSize int = 5
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

	deviceIDListCSV := strings.Join(deviceIDTemp, ",")

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

// func (logic *Logic) UpdateJobCompleted(jobID string, completed bool) error {
// 	logger := log.WithFields(log.Fields{
// 		"package":  "logic",
// 		"struct":   "Logic",
// 		"function": "UpdateJobCompleted",
// 	})

// 	logger.Debug("received UpdateJobCompleted()")

// 	return logic.dbAdapter.UpdateJobCompleted(jobID, completed)
// }

func (logic *Logic) DeleteJob(jobID string) error {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "DeleteJob",
	})

	logger.Debug("received DeleteJob()")

	return logic.dbAdapter.DeleteJob(jobID)
}

// func (logic *Logic) Precompute(jobID string, deviceEndpoint string, podEndpoint string, inputSize int) error {
// 	logger := log.WithFields(log.Fields{
// 		"package": "logic",
// 		"struct" : "Logic",
// 		"function" : "Precompute",
// 	})

// 	logger.Debug("received Precompute()")

// 	var deviceResult, podResult float64
// 	startIdx := 0
// 	// precomputation: 0.1%, serve as endIdx for precomputation
// 	precomputeSize := int(0.001 * float64(inputSize)) 
// 	if precomputeSize < 1 {
// 		precomputeSize = 1 
// 	}

// 	deviceFullEndpoint := logic.GetFullEndpoint(deviceEndpoint, startIdx, precomputeSize)
// 	podFullEndpoint := logic.GetFullEndpoint(podEndpoint, startIdx, precomputeSize)
// 	fmt.Println("precompute device endpoint", deviceFullEndpoint)
// 	fmt.Println("precompute pod endpoint", podFullEndpoint)
	
// 	// set client to pull results from device
// 	deviceResultChan := make(chan float64, 1)
// 	podResultChan := make(chan float64, 1)
// 	errChan := make(chan error, 2)

// 	// precompute device
// 	go func() {
// 		deviceResult, err := logic.ComputeDevice(jobID, deviceEndpoint, startIdx, precomputeSize)
// 			if err != nil {
// 				errChan <- err
// 				return
// 			}
// 		deviceResultChan <- deviceResult
// 	}()
	
// 	// precompute pod
// 	go func() {
// 		podResult, err := logic.ComputePod(jobID, podEndpoint, startIdx, precomputeSize)
// 			if err != nil {
// 				errChan <- err
// 				return
// 			}
// 		podResultChan <- podResult
// 	}()
	
// 	// wait for both results
// 	for i := 0; i < 2; i++ {
// 		select {
// 		case deviceResult = <- deviceResultChan:
// 		case podResult = <- podResultChan:
// 		case err := <- errChan:
// 			return err
// 		}
// 	}

// 	deviceUnitResult := deviceResult / float64(precomputeSize - startIdx)
// 	podUnitResult := podResult / float64(precomputeSize - startIdx)

// 	// set partitionRate based on precomputation result
// 	partitionRate := math.Round(podUnitResult / (deviceUnitResult + podUnitResult) * 100) / 100

// 	deviceEndIdx := int( partitionRate * float64(batchSize) )
	
// 	if deviceEndIdx < 1 {
// 		deviceEndIdx = 1
// 	}
// 	fmt.Println("Before save, deviceEndIdx", deviceEndIdx)

// 	// update device start, end index
// 	err := logic.dbAdapter.UpdateDeviceStartIndex(jobID, 0)
// 	if err != nil {
// 		return err
// 	}

// 	err = logic.dbAdapter.UpdateDeviceEndIndex(jobID, deviceEndIdx)
// 	if err != nil {
// 		return err
// 	}
	
// 	// update pod start, end index
// 	err = logic.dbAdapter.UpdatePodStartIndex(jobID, deviceEndIdx)
// 	if err != nil {
// 		return err
// 	}
	
// 	err = logic.dbAdapter.UpdatePodEndIndex(jobID, batchSize)
// 	if err != nil {
// 		return err
// 	}

// 	// update partition rate, used when device fully offloads and need to reassign works
// 	err = logic.dbAdapter.UpdatePartitionRate(jobID, partitionRate)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

func (logic *Logic) Compute(jobID string, deviceID string) error {
	logger := log.WithFields(log.Fields{
		"package": "logic",
		"struct" : "Logic",
		"function" : "Compute",
	})

	logger.Debugf("received Compute()")

	job, err := logic.GetJob(jobID)
	if err != nil {
		return err
	}

	// assign job first
	err = logic.dbAdapter.UpdateStartIndex(jobID, job.StartIndex + batchSize)
	if err != nil {
		return err
	}

	device, err := logic.GetDevice(deviceID)
	if err != nil {
		return err
	}

	var response map[string]interface{}
	fullEndpoint := logic.GetFullEndpoint(device.Endpoint, job.StartIndex, job.StartIndex + batchSize)
	fmt.Println("compute full endpoint", fullEndpoint)
	
	client := resty.New()
	
	_, err = client.R().SetResult(&response).Get(fullEndpoint)
	if err != nil {
		return err
	}
	
	result, err := logic.HandleResponse(response) // TODO: change return type of handleResponse
	if err != nil {
		return err
	}

	fmt.Println(deviceID, result)
	// later change this line to test whether compute succeed
	
	return nil
}

func (logic *Logic) ScheduleJob(jobID string) error {
	logger := log.WithFields(log.Fields{
		"package": "logic",
		"struct" : "Logic",
		"function" : "ScheduleJob",
	})

	logger.Debug("received ScheduleJob()")

	job, err := logic.dbAdapter.GetJob(jobID)
	if err != nil {
		return err
	}

	// initialize start index as 0
	err = logic.dbAdapter.UpdateStartIndex(jobID, 0)
	if err != nil {
		return err
	}

	computeWithRetry := func(jobID, deviceID string, maxRetries int) error {
        for i := 0; i < maxRetries; i++ {
            err := logic.Compute(jobID, deviceID)
            if err == nil {
                return nil // success, no need to retry
            }
			// Optionally, implement some backoff strategy here
			
        }
		// return failed range(start to end)
        return fmt.Errorf("compute failed after %d attempts", maxRetries)
    }

	// precompute
	// change this later
	// err = logic.Precompute(jobID, device.Endpoint, pod.Endpoint, job.InputSize)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return err
	// }
	// fmt.Println("precomputation ended...")

	var stop bool
	errChan := make(chan error, 1) // buffer of 1 to avoid blocking
	for _, deviceID := range job.DeviceIDList {
        go func(dID string) {
            for !stop {
                err := computeWithRetry(jobID, dID, 3) // set MaxRetry 3
                if err != nil {
                    errChan <- err
                    return
                }
                // check job status before continuing
                currentJob, err := logic.dbAdapter.GetJob(jobID)
                if err != nil || currentJob.Completed {
                    stop = true
                    return
                }
            }
        }(deviceID)
		// sleep for 0.01s, to prevent race condition
		// TODO: find more fancy solution
		time.Sleep(10 * time.Millisecond)
    }
	
	select {
    case err := <-errChan:
        stop = true // signal other goroutines to stop
        return err
    default:
		fmt.Println("Job completed")
	}
	// compute til end
	// for {
	// 	job, err = logic.dbAdapter.GetJob(jobID)
	// 	if err != nil {
	// 		return err
	// 	}
		
	// 	// break when job completes
	// 	if job.Completed {
	// 		break
	// 	}
		
	// 	deviceCount := len(job.DeviceIDList)
	// 	// set channels to get results
	// 	// deviceResultsChan := make(chan float64, device_num)
	// 	errChan := make(chan error, deviceCount)
	// 	wg := sync.WaitGroup{}
	// 	wg.Add(deviceCount)
		
	// 	for _, deviceID := range job.DeviceIDList{
	// 		// compute
	// 		// if err != nil -> do again
	// 		go func(dID string) {
	// 			defer wg.Done()
	// 			err := logic.Compute(jobID, deviceID)
	// 			if err != nil {
	// 				errChan <- err
	// 			}
	// 			// deviceResultsChan <- result
	// 		}(deviceID)

	// 	}
	// 	wg.Wait()
	// 	close(errChan)

	// 	// Check if there were any errors
	// 	if len(errChan) > 0 {
	// 		return <-errChan // returns the first error encountered
	// 	}

		
	// 	// job completed
	// 	if job.StartIndex == job.EndIndex{
	// 		err = logic.dbAdapter.UpdateJobCompleted(jobID, true)
	// 		if err != nil {
	// 			return err
	// 		}
	// 	}		
	// }

	return nil
}

		// batchSize 결정할 때 참고
		// podUnitResult := podResult / (float64(job.PodEndIndex) - float64(job.PodStartIndex))
		// deviceUnitResult := deviceResult / (float64(job.DeviceEndIndex) - float64(job.DeviceStartIndex))

		// calculate partitionRate
		// modified to get unit result
		// partitionRate := math.Round(podUnitResult / (deviceUnitResult + podUnitResult) * 100) / 100
		// fmt.Println("partitionRate", partitionRate)
		// dynamically change batchSize
		// batchSize *= int(1 / partitionRate)
		// if batchSize > job.InputSize{
		// 	batchSize = job.InputSize
		// }

		// calculate workload for device and pod
		// deviceWork := int(float64(batchSize) * partitionRate)
		// fmt.Println("before setting", deviceWork)
		// if deviceWork < 1{
		// 	// if device is 100 times slower than pod
		// 	// allocate minimum work(1 input)
		// 	// pod didn't do work -> get previousPartitionRate
		// 	deviceWork = 1
		// }
		// podWork := batchSize - deviceWork 
		

	// 먼저 끝나면, 끝난 장치에 요청 다시 보내야
	// 다음 배치 가져와서 실행
	// 파드 먼저 끝남, 디바이스 예측 시간이랑 비슷하게 끝나게 설정
	// 파드에서 끝난 시간 알 수 있음
	// 파드에서 * 2로 처리 -> 나중에 변경
	// 디바이스에서는 한 파티션 끝날 때 파드가 두 번째 게 돌고 있을 것 같으므로
	// 파드에서 2번째 것 끝날 시간이랑 거의 겹치게 다음 배치 데이터를 넣어줌
	// 마지막으로 끝난 애 기준으로, 늦은 애가 끝났을 때 빠른 애가 계속 돌고 있으므로
	// 거기에 맞춰서 빠른 애가 몇 번 맞췄는지, 
	// 마지막 세트: 20 = 8 + 8 + 4, 마지막 배치 비율로


	// type Scheduler struct {
		// 	config *commons.Config
		// 	Job_list []*types.Job
		// }
		
		// type ScheduleJob interface {
		// 	getJob(string) (*types.Job)
		// 	Precompute(*types.Job, string, string) (float64, error)
		// 	Compute(*types.Job, string, string) (error)
		// 	// add method if necessary
		// }
		
		// // getJob returns Job from jobID
		// func (scheduler *Scheduler) GetJob (jobID string) (types.Job) {
		// 	var idx int
		// 	jobList := scheduler.Job_list
		
		// 	// using generic, supported after Go 1.21
		// 	// idx = slices.IndexFunc(scheduler.Job_list, func(j_ptr *types.Job) bool { return (*j_ptr).ID == jobID })
		
		// 	// using for loop, slower than Generic, use this for Go version lower than 1.21
		// 	for i := range jobList {
		// 		if jobList[i].ID == jobID {
		// 			idx = i
		// 			break
		// 		}
		// 	}
		
		// 	return *jobList[idx]
		// }

		
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