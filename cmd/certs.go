package cmd

import (
	"github.com/rms1000watt/degeneres/generate"

	"github.com/spf13/cobra"
)

var certsCmd = &cobra.Command{
	Use:   "certs",
	Short: "Generate Certs for HTTPS communication",
	Long:  `Generate Self-Signed, Untrusted Certs for HTTPS communication.`,
	Run:   runCerts,
}

var certsCfg generate.CertsConfig

func init() {
	generateCmd.AddCommand(certsCmd)
	certsCmd.Flags().StringVar(&certsCfg.CertsPath, "certs-path", "./certs", "Output path for newly generated certs")
	certsCmd.Flags().StringVar(&certsCfg.OpensslConfig, "openssl-config", "./certs/openssl.cnf", "Openssl config location")
}

func runCerts(cmd *cobra.Command, args []string) {
	generate.Certs(certsCfg)
}
