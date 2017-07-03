package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	log "github.com/sirupsen/logrus"
)

var logLevel string

var RootCmd = &cobra.Command{
	Use:   "{{.ProjectNameCommander}}",
	Short: "{{.ShortDescription}}",
	Long:  `{{.LongDescription}}`,
}

func init() {
	RootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "info", "Set log level (debug, info, warn, error, fatal)")

	SetPFlagsFromEnv(RootCmd)
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func SetPFlagsFromEnv(cmd *cobra.Command) {
	// Courtesy of https://github.com/coreos/pkg/blob/master/flagutil/env.go
	cmd.PersistentFlags().VisitAll(func(f *pflag.Flag) {
		key := strings.ToUpper(strings.Replace(f.Name, "-", "_", -1))
		if val := os.Getenv(key); val != "" {
			if err := cmd.PersistentFlags().Set(f.Name, val); err != nil {
				fmt.Println("Failed setting flag from env:", err)
			}
		}
	})
}

func SetFlagsFromEnv(cmd *cobra.Command) {
	// Courtesy of https://github.com/coreos/pkg/blob/master/flagutil/env.go
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		key := strings.ToUpper(strings.Replace(f.Name, "-", "_", -1))
		if val := os.Getenv(key); val != "" {
			if err := cmd.Flags().Set(f.Name, val); err != nil {
				fmt.Println("Failed setting flag from env:", err)
			}
		}
	})
}

func configureLogging() {
	if level, err := log.ParseLevel(logLevel); err != nil {
		log.Error("log-level argument malformed: ", logLevel, ": ", err)
		log.SetLevel(log.InfoLevel)
	} else {
		log.SetLevel(level)
	}
}
