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
	ID          string    				`json:"id" gorm:"primaryKey"`
	DeviceID 	string 	  				`json:"device_id"`
	PodID 		string 		  			`json:"pod_id"`
	InputSize 	int 		  			`json:"input_size"` // input size: total number of input for a job
	PartitionRate float64				`json:"partition_rate`
	DeviceStartIndex 	int				`json:"device_start_index"` 
	DeviceEndIndex 	int					`json:"device_end_index"` 
	PodStartIndex 	int					`json:"pod_start_index"` 
	PodEndIndex 	int					`json:"pod_end_index"` 
	Completed	bool 					`json:"completed"`
	CreatedAt   time.Time 				`json:"created_at,omitempty"`
	UpdatedAt   time.Time 				`json:"updated_at,omitempty"`
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

