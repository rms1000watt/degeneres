package cmd

import (
    "fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version of this project",
	Long:  `Print version of this project`,
	Run:   runVersion,
}

func init() {
	RootCmd.AddCommand(versionCmd)
}

func runVersion(cmd *cobra.Command, args []string) {
    fmt.Println("{{.ProjectName}}: {{.Version}}")
    fmt.Println("Degeneres: {{.GeneratorVersion}}")
    fmt.Println("")
    fmt.Println("Built with love using: https://github.com/rms1000watt/degeneres")
}
