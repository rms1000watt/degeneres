package cmd

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var logLevel string

var RootCmd = &cobra.Command{
	Use:   "degeneres",
	Short: "Degeneres: Golang Microservice Generator",
	Long:  `Degeneres: Golang Microservice Generator using Protobuf`,
}

func init() {
	RootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "debug", "Set log level (debug, info, error, fatal)")
}

func Execute() {
	if level, err := log.ParseLevel(logLevel); err != nil {
		log.Error("log-level argument malformed: ", logLevel)
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(level)
	}

	if err := RootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
