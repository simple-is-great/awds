package logic

import (
	"awds/types"
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
)

var (
	// starts small, grow later
	batchSize int = 10
	// hashmap to store failedJob
	// expected key: jobID, expected slice: "startidx-startidx+(batchSize)"
	// should be included in the scheduler module later
	failedJob map[string][]string
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

func (logic *Logic) Compute(device *types.Device, index int, batchSize int) (error) {
	logger := log.WithFields(log.Fields{
		"package": "logic",
		"struct" : "Logic",
		"function" : "Compute",
	})

	logger.Debugf("received Compute()")

	// assign job first
	// err := logic.dbAdapter.UpdateStartIndex(job.ID, job.StartIndex + batchSize)
	// if err != nil {
	// 	return err
	// }

	var response map[string]interface{}
	fullEndpoint := logic.GetFullEndpoint(device.IP, device.Port,device.Endpoint, job.StartIndex, job.StartIndex + batchSize)
	// fmt.Println("compute full endpoint", fullEndpoint)
	
	client := resty.New()
	_, err := client.R().SetResult(&response).Get(fullEndpoint)
	if err != nil {
		return err
	}
	
	result, ok := response["result"].(float64)
	if !ok {
		// failed to compute -> save index
		// later change this if each schedule object holds map
		// return logic.SaveFailedWorkload(job.ID, job.StartIndex, batchSize)
	}

	fmt.Println(device.ID, result)

	return nil
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

    errChan := make(chan error, 1) // Buffer of 1 to avoid blocking
	
	job, err := logic.dbAdapter.GetJob(jobID)
	if err != nil {
		return err
	}

	// initialize start index as 0
	// err = logic.dbAdapter.UpdateStartIndex(jobID, 0)
	// if err != nil {
	// 	return err
	// }

	// precompute
	// change this later
	// err = logic.Precompute(jobID, device.Endpoint, pod.Endpoint, job.InputSize)
	idx := job.StartIndex
	// 사용 가능한 디바이스 리스트 / 맵 구조 구현 필요
	// 1. scheduler 만들어서 scheduler에 queue를 저장 -> 여러 schedule들에 대해서 queue를 관리하기 어렵다
	// 2. init 할 때 queue를 생성
	
	var q Queue
	for _, deviceID := range job.DeviceIDList{
		q.Enqueue(deviceID)
	}
	
	for {
		flag := 0 // flag to exit infinite for loop
		// TODO: we need to calculate index and pass into threads
		// 
        go func(){
            for {
				// Check job status before continuing
				currentJob, err := logic.dbAdapter.GetJob(jobID)
				if err != nil {
					errChan <- err
					return
				}
				
				if currentJob.Completed {
					flag = 1
					return
				}

				dID := q.Dequeue().(string)
				
				device, err := logic.dbAdapter.GetDevice(dID)
				if err != nil {
					return
				}
				// device별 batchSize 계산 -> idx 사용할 것
				err = logic.Compute(&device, idx, batchSize)
				// TODO: change to compute
				if err != nil {
					errChan <- err
					return
				}
				q.Enqueue(dID)
			}
			// sleep for 0.01s, to prevent race condition
			// TODO: find more fancy solution
			// time.Sleep(10 * time.Millisecond)
		}()
		if flag == 1 {
			break
		}
	}
	
	// Wait for all goroutines to finish
	go func() {
        wg.Wait()
        close(errChan) // Close the channel to signal completion
    }()

    // Handle errors from goroutines
    for err := range errChan {
        cancel() // Cancel all goroutines on error
        return err // Return the first error encountered
    }
	
	// Update Completed to true outside of the goroutines to ensure it's only done once
    err = logic.dbAdapter.UpdateJobCompleted(jobID, true)
    if err != nil {
        return err
    }

	// // check for error
	// if len(failedJob[jobID]) > 0 {
	// 	for _, deviceID := range job.DeviceIDList {
	// 		go func(dID string) {
	// 			for !stop {
	// 				// need to get new batchSize here
	// 				// TODO: need to modify ComputeWithRetry to accomodate failedJob
	// 				err := logic.ComputeWithRetry(jobID, dID, batchSize, 3) // set MaxRetry 3
	// 				if err != nil {
	// 					errChan <- err
	// 					return
	// 				}
	// 				// check job status before continuing
	// 				currentJob, err := logic.dbAdapter.GetJob(jobID)
	// 				if err != nil || currentJob.Completed {
	// 					stop = true
	// 					return
	// 				}
	// 			}
	// 		}(deviceID)
	// 		// sleep for 0.01s, to prevent race condition
	// 		// TODO: find more fancy solution
	// 		time.Sleep(10 * time.Millisecond)
	// 	}
	// }
	
	return nil
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