package cmd

import (
	"github.com/rms1000watt/degeneres/generate"

	"github.com/spf13/cobra"
)

var (
	certsPath   string
	commonName  string
	letsEncrypt bool
)

var certsCmd = &cobra.Command{
	Use:   "certs",
	Short: "Generate Certs for HTTPS communication",
	Long:  `Generate Self-Signed, Untrusted Certs for HTTPS communication`,
	Run:   runCerts,
}

func init() {
	generateCmd.AddCommand(certsCmd)
	// TODO: Separate this out better. openssl-config, out-path. Update code to reflect
	certsCmd.Flags().StringVar(&certsPath, "certs-path", "./certs", "Certs path that contains openssl.cnf")
	certsCmd.Flags().StringVar(&commonName, "common-name", "localhost", "Common Name for the cert, ie. localhost")
}

func runCerts(cmd *cobra.Command, args []string) {
	generate.Certs(certsPath, commonName)
}
