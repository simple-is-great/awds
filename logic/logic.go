package logic

import (
	"awds/commons"
	"awds/db"
)

type Logic struct {
	config *commons.Config
	// scheduler *schedule.Scheduler
	dbAdapter  *db.DBAdapter
}

// Start starts Logic
func Start(config *commons.Config, dbAdapter *db.DBAdapter) (*Logic, error) {
	logic := &Logic{
		config:     config,
		dbAdapter:  dbAdapter,
		// scheduler: 	scheduler,
	}

	return logic, nil
}

// Stop stops Logic
func (logic *Logic) Stop() error {
	return nil
}
