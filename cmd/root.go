package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "degeneres",
	Short: "Degeneres: Golang Microservice Generator",
	Long:  `Degeneres: Golang Microservice Generator using Protobuf`,
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
