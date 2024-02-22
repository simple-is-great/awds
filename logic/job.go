package logic

import (
	"awds/types"
	"fmt"
	"math"

	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
)

var (
	// starts small, grow later
	batchSize int = 10
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

func (logic *Logic) CreateJob(job *types.Job) error {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "CreateJob",
	})

	logger.Debug("received CreateJob()")

	return logic.dbAdapter.InsertJob(job)
}

func (logic *Logic) UpdateJobDevice (jobID string, deviceID string) error {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "UpdateJobDevice",
	})

	logger.Debug("received UpdateJobDevice()")

	return logic.dbAdapter.UpdateJobDevice(jobID, deviceID)
}

func (logic *Logic) UpdateJobPod(jobID string, podID string) error {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "UpdateJobPod",
	})

	logger.Debug("received UpdateJobPod()")

	return logic.dbAdapter.UpdateJobPod(jobID, podID)
}

func (logic *Logic) UpdateInputSize(jobID string, inputSize int) error {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "UpdateInputSize",
	})

	logger.Debug("received UpdateInputSize()")

	return logic.dbAdapter.UpdateInputSize(jobID, inputSize)
}

func (logic *Logic) UpdatePartitionRate(jobID string, partitionRate float64) error {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "UpdatePartitionRate",
	})

	logger.Debug("received UpdatePartitionRate()")

	return logic.dbAdapter.UpdatePartitionRate(jobID, partitionRate)
}

func (logic *Logic) UpdateJobCompleted(jobID string, completed bool) error {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "UpdateJobCompleted",
	})

	logger.Debug("received UpdateJobCompleted()")

	return logic.dbAdapter.UpdateJobCompleted(jobID, completed)
}

func (logic *Logic) UpdateDeviceStartIndex(jobID string, deviceStartIndex int) error {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "UpdateDeviceStartIndex",
	})

	logger.Debug("received UpdateDeviceStartIndex()")

	return logic.dbAdapter.UpdateDeviceStartIndex(jobID, deviceStartIndex)
}


func (logic *Logic) UpdateDeviceEndIndex(jobID string, deviceEndIndex int) error {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "UpdateDeviceEndIndex",
	})

	logger.Debug("received UpdateDeviceEndIndex()")

	return logic.dbAdapter.UpdateDeviceEndIndex(jobID, deviceEndIndex)
}


func (logic *Logic) UpdatePodStartIndex(jobID string, podStartIndex int) error {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "UpdatePodStartIndex",
	})

	logger.Debug("received UpdatePodStartIndex()")

	return logic.dbAdapter.UpdatePodStartIndex(jobID, podStartIndex)
}



func (logic *Logic) UpdatePodEndIndex(jobID string, podEndIndex int) error {
	logger := log.WithFields(log.Fields{
		"package":  "logic",
		"struct":   "Logic",
		"function": "UpdatePodEndIndex",
	})

	logger.Debug("received UpdatePodEndIndex()")

	return logic.dbAdapter.UpdatePodEndIndex(jobID, podEndIndex)
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

func (logic *Logic) Precompute(jobID string, deviceEndpoint string, podEndpoint string, inputSize int) error {
	logger := log.WithFields(log.Fields{
		"package": "logic",
		"struct" : "Logic",
		"function" : "Precompute",
	})

	logger.Debug("received Precompute()")

	var deviceResult, podResult float64
	startIdx := 0
	// precomputation: 0.1%, serve as endIdx for precomputation
	precomputeSize := int(0.001 * float64(inputSize)) 
	if precomputeSize < 1 {
		precomputeSize = 1 
	}

	deviceFullEndpoint := logic.GetFullEndpoint(deviceEndpoint, startIdx, precomputeSize)
	podFullEndpoint := logic.GetFullEndpoint(podEndpoint, startIdx, precomputeSize)
	fmt.Println("precompute device endpoint", deviceFullEndpoint)
	fmt.Println("precompute pod endpoint", podFullEndpoint)
	
	// set client to pull results from device
	deviceResultChan := make(chan float64, 1)
	podResultChan := make(chan float64, 1)
	errChan := make(chan error, 2)

	// precompute device
	go func() {
		deviceResult, err := logic.ComputeDevice(jobID, deviceEndpoint, startIdx, precomputeSize)
			if err != nil {
				errChan <- err
				return
			}
		deviceResultChan <- deviceResult
	}()
	
	// precompute pod
	go func() {
		podResult, err := logic.ComputePod(jobID, podEndpoint, startIdx, precomputeSize)
			if err != nil {
				errChan <- err
				return
			}
		podResultChan <- podResult
	}()
	
	// wait for both results
	for i := 0; i < 2; i++ {
		select {
		case deviceResult = <- deviceResultChan:
		case podResult = <- podResultChan:
		case err := <- errChan:
			return err
		}
	}

	deviceUnitResult := deviceResult / float64(precomputeSize - startIdx)
	podUnitResult := podResult / float64(precomputeSize - startIdx)

	// set partitionRate based on precomputation result
	partitionRate := math.Round(podUnitResult / (deviceUnitResult + podUnitResult) * 100) / 100

	// if deviceResult - podResult > 0.99 {
	// 	batchSize *= int(1 / partitionRate)
	// 	if batchSize > inputSize {
	// 		batchSize = inputSize
	// 	}
	// }

	deviceEndIdx := int( partitionRate * float64(batchSize) )
	
	if deviceEndIdx < 1 {
		deviceEndIdx = 1
	}
	fmt.Println("Before save, deviceEndIdx", deviceEndIdx)

	// update device start, end index
	err := logic.dbAdapter.UpdateDeviceStartIndex(jobID, 0)
	if err != nil {
		return err
	}

	err = logic.dbAdapter.UpdateDeviceEndIndex(jobID, deviceEndIdx)
	if err != nil {
		return err
	}
	
	// update pod start, end index
	err = logic.dbAdapter.UpdatePodStartIndex(jobID, deviceEndIdx)
	if err != nil {
		return err
	}
	
	err = logic.dbAdapter.UpdatePodEndIndex(jobID, batchSize)
	if err != nil {
		return err
	}

	// update partition rate, used when device fully offloads and need to reassign works
	err = logic.dbAdapter.UpdatePartitionRate(jobID, partitionRate)
	if err != nil {
		return err
	}
	return nil
}

func (logic *Logic) ComputeDevice(jobID string, deviceEndpoint string, startIdx int, endIdx int) (float64, error) {
	logger := log.WithFields(log.Fields{
		"package": "logic",
		"struct" : "Logic",
		"function" : "ComputeDevice",
	})

	logger.Debugf("received ComputeDevice()")

	var deviceResponse map[string]interface{}
	client := resty.New()
	// request to device	
	deviceFullEndpoint := logic.GetFullEndpoint(deviceEndpoint, startIdx, endIdx)
	fmt.Println("compute device full endpoint", deviceFullEndpoint)

	_, err := client.R().SetResult(&deviceResponse).Get(deviceFullEndpoint)
	if err != nil {
		return -1, err
	}
	// HandleResponse transforms deviceResponse object to float64
	deviceResult, err := logic.HandleResponse(deviceResponse)
	if err != nil {
		return -1, err
	}

	fmt.Println("device", deviceResult)
	
	return deviceResult, nil
}

func (logic *Logic) ComputePod(jobID string, podEndpoint string, startIdx int, endIdx int) (float64, error) {
	logger := log.WithFields(log.Fields{
		"package": "logic",
		"struct" : "Logic",
		"function" : "ComputePod",
	})

	logger.Debugf("received ComputePod()")

	var podResponse map[string]interface{}
	client := resty.New()	
	podFullEndpoint := logic.GetFullEndpoint(podEndpoint, startIdx, endIdx)
	fmt.Println("compute pod full endpoint", podFullEndpoint)

	// request to pod	
	_, err := client.R().SetResult(&podResponse).Get(podFullEndpoint)
	if err != nil {
		return -1, err
	}

	podResult , err := logic.HandleResponse(podResponse)
	if err != nil {
		return -1, err
	}
	
	fmt.Println("pod", podResult)

	return podResult, nil

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

	device, err := logic.dbAdapter.GetDevice(job.DeviceID)
	if err != nil {
		return err
	}

	pod, err := logic.dbAdapter.GetPod(job.PodID)
	if err != nil {
		return err
	}

	// if job finishes in single batch, set batchSize to inputSize
	if batchSize >= job.InputSize {
		batchSize = job.InputSize		
	}

	// precompute
	err = logic.Precompute(jobID, device.Endpoint, pod.Endpoint, job.InputSize)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("precomputation ended...")

	// compute til end
	for {
		job, err = logic.dbAdapter.GetJob(jobID)
		if err != nil {
			return err
		}
		
		// break when job completes
		if job.Completed {
			break
		}
		
		//set channels to get device, pod results
		deviceResultsChan := make(chan float64, 1)
		podResultsChan := make(chan float64, 1)
		errChan := make(chan error, 2)
		var deviceResult, podResult float64

		// compute device
		go func() {
			result, err := logic.ComputeDevice(jobID, device.Endpoint, job.DeviceStartIndex, job.DeviceEndIndex)
			if err != nil {
				errChan <- err
				return
			}
			deviceResultsChan <- result
		}()

		// compute pod
		go func() {
			result, err := logic.ComputePod(jobID, pod.Endpoint, job.PodStartIndex, job.PodEndIndex)
			if err != nil {
				errChan <- err
				return
			}
			podResultsChan <- result
		}()
		
		// TODO
		// if pod or device ends but the other doesn't end
		// run next job

		// wait for the results
		for i :=0; i < 2; i++{
			select {
			case result := <-deviceResultsChan:
				deviceResult = result
			case result := <-podResultsChan:
				podResult = result
			case err := <-errChan:
				return err
			}
		}
		
		podUnitResult := podResult / (float64(job.PodEndIndex) - float64(job.PodStartIndex))
		deviceUnitResult := deviceResult / (float64(job.DeviceEndIndex) - float64(job.DeviceStartIndex))

		// calculate partitionRate
		// modified to get unit result
		partitionRate := math.Round(podUnitResult / (deviceUnitResult + podUnitResult) * 100) / 100
		fmt.Println("partitionRate", partitionRate)
		
		// dynamically change batchSize
		// batchSize *= int(1 / partitionRate)
		// if batchSize > job.InputSize{
		// 	batchSize = job.InputSize
		// }

		// calculate workload for device and pod
		deviceWork := int(float64(batchSize) * partitionRate)
		fmt.Println("before setting", deviceWork)
		if deviceWork < 1{
			// if device is 100 times slower than pod
			// allocate minimum work(1 input)
			// pod didn't do work -> get previousPartitionRate
			deviceWork = 1
		}
		fmt.Println("after setting", deviceWork)
		podWork := batchSize - deviceWork 
		
		fmt.Println("deviceWork", deviceWork)
		fmt.Println("podWork", podWork)

		newDeviceStartIndex := job.PodEndIndex
		newDeviceEndIndex := job.PodEndIndex + deviceWork

		// update device start, end index
		if job.PodEndIndex >= job.InputSize {
			newDeviceStartIndex = job.InputSize
		}

		if newDeviceEndIndex >= job.InputSize {
			newDeviceEndIndex = job.InputSize 
		}

		err = logic.dbAdapter.UpdateDeviceStartIndex(jobID, newDeviceStartIndex)
		if err != nil {
			return err
		}
		
		err = logic.dbAdapter.UpdateDeviceEndIndex(jobID, newDeviceEndIndex)
		if err != nil {
			return err
		}

		newPodStartIndex := job.PodEndIndex + deviceWork
		newPodEndIndex := newPodStartIndex + podWork

		// update pod start index and end index
		if newPodStartIndex > job.InputSize {
			newPodStartIndex = job.InputSize 
		}
		if newPodEndIndex > job.InputSize {
			newPodEndIndex = job.InputSize 
		}

		err = logic.dbAdapter.UpdatePodStartIndex(jobID, newPodStartIndex)
		if err != nil {
			return err
		}

		err = logic.dbAdapter.UpdatePodEndIndex(jobID, newPodEndIndex)
		if err != nil {
			return err
		}
		
		// job completed
		if ( job.DeviceEndIndex == job.InputSize ) || ( job.PodEndIndex == job.InputSize ){
			err = logic.dbAdapter.UpdateJobCompleted(jobID, true)
			if err != nil {
				return err
			}
		}		
	}

	return nil
}

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