package main

import (
	"awds/commons"
	"awds/db"
	"awds/logic"
	"awds/rest"
	"fmt"
	"os"
	"os/signal"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var config *commons.Config

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "awds [args..]",
	Short: "Distributes workload adaptively",
	Long:  `AWDS distributes workload adaptively.`,
	RunE:  processCommand,
}

func Execute() error {
	return rootCmd.Execute()
}

func processCommand(command *cobra.Command, args []string) error {
	logger := log.WithFields(log.Fields{
		"package":  "main",
		"function": "processCommand",
	})

	cont, err := processFlags(command)
	if err != nil {
		logger.Error(err)
	}

	if !cont {
		return err
	}

	// start service
	logger.Info("Starting DB Adapter...")
	dbAdapter, err := db.Start(config)
	if err != nil {
		logger.Fatal(err)
	}
	defer dbAdapter.Stop()
	logger.Info("DB Adapter Started")

	// logger.Info("Starting Scheduler...")
	// scheduler, err := schedule.Start(config)
	// if err != nil {
	// 	logger.Fatal(err)
	// }
	// defer scheduler.Stop()
	// logger.Info("Scheduler Started")
	
	logik, err := logic.Start(config, dbAdapter)
	if err != nil {
		logger.Fatal(err)
	}
	defer logik.Stop()

	logger.Info("Starting REST Adapter...")
	restAdapter, err := rest.Start(config, logik)
	if err != nil {
		logger.Fatal(err)
	}
	defer restAdapter.Stop()
	logger.Info("REST Adapter Started")

	// wait
	fmt.Println("Press Ctrl+C to stop...")
	waitForCtrlC()

	return nil
}

func main() {
	log.SetFormatter(&log.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05.000",
		FullTimestamp:   true,
	})

	config = commons.GetDefaultConfig()
	log.SetLevel(config.GetLogLevel())

	logger := log.WithFields(log.Fields{
		"package":  "main",
		"function": "main",
	})

	// attach common flags
	setCommonFlags(rootCmd)

	err := Execute()
	if err != nil {
		logger.Fatal(err)
		os.Exit(1)
	}
}

func waitForCtrlC() {
	var endWaiter sync.WaitGroup

	endWaiter.Add(1)
	signalChannel := make(chan os.Signal, 1)

	signal.Notify(signalChannel, os.Interrupt)

	go func() {
		<-signalChannel
		endWaiter.Done()
	}()

	endWaiter.Wait()
}
