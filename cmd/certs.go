package cmd

import (
	"github.com/rms1000watt/degeneres/generate/certs"

	"github.com/spf13/cobra"
)

var certsCmd = &cobra.Command{
	Use:   "certs",
	Short: "Generate Certs for HTTPS communication",
	Long:  `Generate Self-Signed, Untrusted Certs for HTTPS communication.`,
	Run:   runCerts,
}

var certsCfg certs.Config

func init() {
	generateCmd.AddCommand(certsCmd)
	certsCmd.Flags().StringVarP(&certsCfg.OutputPath, "output-path", "o", "./certs", "Output path for newly generated certs")
	certsCmd.Flags().StringVarP(&certsCfg.OpensslConfig, "openssl-config", "f", "./certs/openssl.cnf", "Openssl config location")
}

func runCerts(cmd *cobra.Command, args []string) {
	configureLogging()

	certs.Certs(certsCfg)
}
