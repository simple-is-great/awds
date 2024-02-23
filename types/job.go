package types

import (
	"fmt"
	"strings"
	"time"

	"github.com/rs/xid"
	"golang.org/x/xerrors"
)

const (
	jobIDPrefix string = "job"
)

// Job represents an job, holding all necessary info. about job
type Job struct {
	ID          	string    			`json:"id" gorm:"primaryKey"`
	DeviceIDList 	[]string			`json:"device_id_list"`
	StartIndex		int					`json:"start_index"`
	EndIndex 		int 		  		`json:"end_index"` // equals to the inputSize
	Scheduled		bool				`json:"scheduled"`
	Completed		bool 				`json:"completed"`
	CreatedAt   	time.Time 			`json:"created_at,omitempty"`
	UpdatedAt   	time.Time 			`json:"updated_at,omitempty"`
}

type JobSQLiteObj struct {
	ID          	string    			`json:"id" gorm:"primaryKey"`
	DeviceIDList 	string				`json:"device_id_list"` // store it as comma-separated string
	StartIndex		int					`json:"start_index"`
	EndIndex 		int 		  		`json:"end_index"` // equals to the inputSize
	Scheduled		bool				`json:"scheduled"`
	Completed		bool 				`json:"completed"`
	CreatedAt   	time.Time 			`json:"created_at,omitempty"`
	UpdatedAt   	time.Time 			`json:"updated_at,omitempty"`
}

func (job *Job) ToJobSQLiteObj() JobSQLiteObj {
	deviceIDList := []string{}
	for _, deviceID := range job.DeviceIDList{
		deviceIDList = append(deviceIDList, deviceID)
	}

	deviceIDListCSV := strings.Join(deviceIDList, ",")

	return JobSQLiteObj{
		ID: 			job.ID,
		DeviceIDList: 	deviceIDListCSV,
		StartIndex:  	job.StartIndex,
		EndIndex: 		job.EndIndex,
		Scheduled: 		job.Scheduled,
		Completed: 		job.Completed,
		CreatedAt: 		job.CreatedAt,
		UpdatedAt: 		job.UpdatedAt,
	}
}

func (job *JobSQLiteObj) ToJobObj() Job {
	deviceIDList := strings.Split(job.DeviceIDList, ",")
	
	return Job {
		ID: 			job.ID,
		DeviceIDList: 	deviceIDList,
		StartIndex:  	job.StartIndex,
		EndIndex: 		job.EndIndex,
		Scheduled: 		job.Scheduled,
		Completed: 		job.Completed,
		CreatedAt: 		job.CreatedAt,
		UpdatedAt: 		job.UpdatedAt,
	}
}

func ValidateJobID(id string) error {
	if len(id) == 0 {
		return xerrors.Errorf("empty job id")
	}

	prefix := fmt.Sprintf("%s_", jobIDPrefix)

	if !strings.HasPrefix(id, prefix) {
		return xerrors.Errorf("invalid job id - %s", id)
	}
	return nil
}

// NewAppID creates a new App ID
func NewJobID() string {
	return fmt.Sprintf("%s_%s", jobIDPrefix, xid.New().String())
}

