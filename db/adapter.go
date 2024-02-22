package db

import (
	"awds/commons"
	"awds/types"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// we use GORM here to connect DataBase
// check out GORM examples below
// https://gorm.io/docs/create.html

const (
	SQLiteDBFileName string = "awds.db"
)

type DBAdapter struct {
	config *commons.Config
	db     *gorm.DB
}

func RemoveDBFile(config *commons.Config) error {
	logger := log.WithFields(log.Fields{
		"package":  "db",
		"function": "RemoveDBFile",
	})

	absPath, err := filepath.Abs(SQLiteDBFileName)
	if err != nil {
		return err
	}

	fi, err := os.Stat(absPath)
	if err == nil && !fi.IsDir() {
		// exist
		logger.Debugf("Removing db file %s", absPath)
		return os.RemoveAll(SQLiteDBFileName)
	}

	return nil
}

// Start starts DBAdapter
func Start(config *commons.Config) (*DBAdapter, error) {
	db, err := gorm.Open(sqlite.Open(SQLiteDBFileName), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(types.Device{}, types.Pod{}, types.Job{})
	if err != nil {
		return nil, err
	}

	adapter := &DBAdapter{
		config: config,
		db:     db,
	}

	return adapter, nil
}

// Stop stops DBAdapter
func (adapter *DBAdapter) Stop() error {
	return nil
}
