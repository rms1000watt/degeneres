package cmd

import (
	"github.com/rms1000watt/degeneres/generate"

	"github.com/spf13/cobra"
)

var (
	genCfg      = generate.Config{}
	generateCmd = &cobra.Command{
		Use:   "generate",
		Short: "Generates code from a `proto` file",
		Long:  `Generates code from a "proto" file`,
		Run:   runGenerate,
	}
)

func init() {
	RootCmd.AddCommand(generateCmd)

	generateCmd.Flags().StringVarP(&genCfg.ProtoFilePath, "proto-file", "f", "./pb/test.proto", "Protobuf filepath used for generation")
	generateCmd.Flags().BoolVarP(&genCfg.Verbose, "verbose", "v", false, "Generator verbosity")
}

func runGenerate(cmd *cobra.Command, args []string) {
	generate.Generate(genCfg)
}
