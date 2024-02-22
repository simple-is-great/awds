package db

import (
	"awds/types"

	"golang.org/x/xerrors"
)

func (adapter *DBAdapter) ListJobs() ([]types.Job, error) {
	jobs := []types.Job{}
	result := adapter.db.Find(&jobs)
	if result.Error != nil {
		return nil, result.Error
	}

	return jobs, nil
}

func (adapter *DBAdapter) GetJob(jobID string) (types.Job, error) {
	var job types.Job
	result := adapter.db.Where("id = ?", jobID).First(&job)
	if result.Error != nil {
		return job, result.Error
	}

	return job, nil
}

func (adapter *DBAdapter) InsertJob(job *types.Job) error {
	result := adapter.db.Create(&job)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected != 1 {
		return xerrors.Errorf("failed to insert an job")
	}

	return nil
}

func (adapter *DBAdapter) UpdateJobDevice(jobID string, deviceID string) error {
	var record types.Job
	result := adapter.db.Where("id = ?", jobID).Find(&record)
	if result.Error != nil {
		return result.Error
	}

	record.DeviceID = deviceID

	adapter.db.Save(&record)

	return nil
}

func (adapter *DBAdapter) UpdateJobPod(jobID string, podID string) error {
	var record types.Job
	result := adapter.db.Where("id = ?", jobID).Find(&record)
	if result.Error != nil {
		return result.Error
	}

	record.PodID = podID

	adapter.db.Save(&record)

	return nil
}

func (adapter *DBAdapter) UpdateInputSize(jobID string, inputSize int) error {
	var record types.Job
	result := adapter.db.Where("id = ?", jobID).Find(&record)
	if result.Error != nil {
		return result.Error
	}

	record.InputSize = inputSize

	adapter.db.Save(&record)

	return nil
}

func (adapter *DBAdapter) UpdatePartitionRate(jobID string, partitionRate float64) error {
	var record types.Job
	result := adapter.db.Where("id = ?", jobID).Find(&record)
	if result.Error != nil {
		return result.Error
	}

	record.PartitionRate = partitionRate

	adapter.db.Save(&record)

	return nil
}

func (adapter *DBAdapter) UpdateDeviceStartIndex(jobID string, deviceStartIndex int) error {
	var record types.Job
	result := adapter.db.Where("id = ?", jobID).Find(&record)
	if result.Error != nil {
		return result.Error
	}

	record.DeviceStartIndex = deviceStartIndex

	adapter.db.Save(&record)

	return nil
}

func (adapter *DBAdapter) UpdateDeviceEndIndex(jobID string, deviceEndIndex int) error {
	var record types.Job
	result := adapter.db.Where("id = ?", jobID).Find(&record)
	if result.Error != nil {
		return result.Error
	}

	record.DeviceEndIndex = deviceEndIndex

	adapter.db.Save(&record)

	return nil
}

func (adapter *DBAdapter) UpdatePodStartIndex(jobID string, podStartIndex int) error {
	var record types.Job
	result := adapter.db.Where("id = ?", jobID).Find(&record)
	if result.Error != nil {
		return result.Error
	}

	record.PodStartIndex = podStartIndex

	adapter.db.Save(&record)

	return nil
}

func (adapter *DBAdapter) UpdatePodEndIndex(jobID string, podEndIndex int) error {
	var record types.Job
	result := adapter.db.Where("id = ?", jobID).Find(&record)
	if result.Error != nil {
		return result.Error
	}

	record.PodEndIndex = podEndIndex

	adapter.db.Save(&record)

	return nil
}

func (adapter *DBAdapter) UpdateJobCompleted(jobID string, completed bool) error {
	var record types.Job
	result := adapter.db.Where("id = ?", jobID).Find(&record)
	if result.Error != nil {
		return result.Error
	}

	record.Completed = completed

	adapter.db.Save(&record)

	return nil
}

func (adapter *DBAdapter) DeleteJob(jobID string) error{
	var job types.Job
	result := adapter.db.Where("id = ?", jobID).Delete(&job)
	
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected != 1 {
		return xerrors.Errorf("failed to delete a job")
	}
	return nil
}