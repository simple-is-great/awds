package main

import (
	"awds/commons"
	"awds/db"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	configPath   string
	envConfig    bool
	version      bool
	help         bool
	debug        bool
	clearDB      bool
)

func setCommonFlags(command *cobra.Command) {
	command.Flags().StringVarP(&configPath, "config", "c", "", "Set config file (yaml or json)")
	command.Flags().BoolVarP(&envConfig, "envconfig", "e", false, "Read config from environmental variables")
	command.Flags().BoolVarP(&version, "version", "v", false, "Print version")
	command.Flags().BoolVarP(&help, "help", "h", false, "Print help")
	command.Flags().BoolVarP(&debug, "debug", "d", false, "Enable debug mode")
	command.Flags().BoolVar(&clearDB, "clear_db", false, "Clear DB data")
}

func processFlags(command *cobra.Command) (bool, error) {
	logger := log.WithFields(log.Fields{
		"package":  "main",
		"function": "processFlags",
	})

	if debug {
		log.SetLevel(log.DebugLevel)
	}

	if help {
		printHelp(command)
		return false, nil // stop here
	}

	if version {
		printVersion(command)
		return false, nil // stop here
	}

	if len(configPath) > 0 {
		loadedConfig, err := commons.LoadConfigFile(configPath)
		if err != nil {
			logger.Error(err)
			return false, err // stop here
		}

		// overwrite config
		config = loadedConfig
	}

	if envConfig {
		loadedConfig, err := commons.LoadConfigEnv()
		if err != nil {
			logger.Error(err)
			return false, err // stop here
		}

		// overwrite config
		config = loadedConfig
	}

	log.SetLevel(config.GetLogLevel())

	// prioritize command-line flag over config files
	if debug {
		log.SetLevel(log.DebugLevel)
	}

	if clearDB {
		// clear db
		err := db.RemoveDBFile(config)
		if err != nil {
			logger.Error(err)
			return false, err
		}
	}

	return true, nil // contiue
}

func printVersion(command *cobra.Command) error {
	info, err := commons.GetVersionJSON()
	if err != nil {
		return err
	}

	fmt.Println(info)
	return nil
}

func printHelp(command *cobra.Command) error {
	return command.Usage()
}
