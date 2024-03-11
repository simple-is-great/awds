package db

import (
	"awds/types"

	"golang.org/x/xerrors"
)

func (adapter *DBAdapter) ListJobs() ([]types.Job, error) {
	sqliteJobs := []types.JobSQLiteObj{}
	result := adapter.db.Find(&sqliteJobs)
	if result.Error != nil {
		return nil, result.Error
	}

	// convert to Job
	jobs := []types.Job{}
	for _, sqliteJob := range sqliteJobs {
		jobs = append(jobs, sqliteJob.ToJobObj())
	}

	return jobs, nil
}

func (adapter *DBAdapter) GetJob(jobID string) (types.Job, error) {
	var sqliteJob types.JobSQLiteObj
	var job types.Job
	result := adapter.db.Where("id = ?", jobID).First(&sqliteJob)
	if result.Error != nil {
		return job, result.Error
	}

	// convert to Job
	job = sqliteJob.ToJobObj()

	return job, nil
}

func (adapter *DBAdapter) InsertJob(job *types.Job) error {
	// convert to JobSQLiteObj
	sqliteJob := job.ToJobSQLiteObj()

	result := adapter.db.Create(&sqliteJob)
	if result.Error != nil {
		return result.Error
	}
	
	if result.RowsAffected != 1 {
		return xerrors.Errorf("failed to insert a job")
	}

	return nil
}

func (adapter *DBAdapter) UpdateDeviceIDList(jobID string, deviceIDList string) error {
	var record types.JobSQLiteObj
	result := adapter.db.Where("id = ?", jobID).Find(&record)
	if result.Error != nil {
		return result.Error
	}

	record.DeviceIDList = deviceIDList

	adapter.db.Save(&record)

	return nil
}

func (adapter *DBAdapter) UpdateStartIndex(jobID string, startIndex int) error {
	var record types.JobSQLiteObj
	result := adapter.db.Where("id = ?", jobID).Find(&record)
	if result.Error != nil {
		return result.Error
	}

	record.StartIndex = startIndex

	adapter.db.Save(&record)

	return nil
}

func (adapter *DBAdapter) UpdateEndIndex(jobID string, endIndex int) error {
	var record types.JobSQLiteObj
	result := adapter.db.Where("id = ?", jobID).Find(&record)
	if result.Error != nil {
		return result.Error
	}

	record.EndIndex = endIndex

	adapter.db.Save(&record)

	return nil
}

func (adapter *DBAdapter) UpdateJobCompleted(jobID string, completed bool) error {
	var record types.JobSQLiteObj
	result := adapter.db.Where("id = ?", jobID).Find(&record)
	if result.Error != nil {
		return result.Error
	}

	record.Completed = completed

	adapter.db.Save(&record)

	return nil
}

func (adapter *DBAdapter) DeleteJob(jobID string) error{
	var job types.JobSQLiteObj
	result := adapter.db.Where("id = ?", jobID).Delete(&job)
	
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected != 1 {
		return xerrors.Errorf("failed to delete a job")
	}
	return nil
}

// func (adapter *DBAdapter) PrepareToCompute(jobID string, batchSize int) (types.Job, int, error) {
// 	var sqliteJob types.JobSQLiteObj
// 	var startIndex int
// 	var newStartIndex int
// 	var job types.Job
// 	// transaction
// 	err := adapter.db.Transaction(func(tx *gorm.DB) error{
// 		if err := tx.Where("id = ?", jobID).First(&sqliteJob).Error; err != nil {
// 			// failed to get sqliteJob
// 			return err
// 		}
// 		// claim workload(update startIndex)
// 		newStartIndex = startIndex + batchSize
// 		sqliteJob.StartIndex = newStartIndex

// 		adapter.db.Save(&sqliteJob)
// 		return nil}, 
// 		// sqlite only supports SERIALIZABLE(DEFAULT), SNAPSHOT ISOLATION, READ UNCOMMITTED
// 		// may change to 
// 		&sql.TxOptions{Isolation: sql.LevelReadUncommitted}) 

// 	// err occurred during transaction
// 	if err != nil {
// 		return job, 0, err
// 	}

// 	// change back to Job object
// 	job = (&sqliteJob).ToJobObj()

// 	return job, startIndex, err
// }