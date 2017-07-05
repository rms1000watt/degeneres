package cmd

import (
	"{{.ImportPath}}/{{.Camel}}"
	"github.com/spf13/cobra"
)

var {{.Camel}}Cfg {{.Camel}}.Config

// {{.Camel}}Cmd represents the {{.Camel}} command
var {{.Camel}}Cmd = &cobra.Command{
	Use:   "{{.Dash}}",
	Short: "{{.ShortDescription}}",
	Long: `{{.LongDescription}}`,
	Run: run{{.TitleCamel}},
}

func init() {
	RootCmd.AddCommand({{.Camel}}Cmd)

	{{.Camel}}Cmd.Flags().StringVar(&{{.Camel}}Cfg.Host, "host", "0.0.0.0", "Host address for server")
	{{.Camel}}Cmd.Flags().IntVar(&{{.Camel}}Cfg.Port, "port", 8080, "Port for server")

	{{.Camel}}Cmd.Flags().StringVar(&{{.Camel}}Cfg.CertsPath, "certs-path", "", "Path for certs")
	{{.Camel}}Cmd.Flags().StringVar(&{{.Camel}}Cfg.KeyName, "key-name", "", "Private key name in certs path")
	{{.Camel}}Cmd.Flags().StringVar(&{{.Camel}}Cfg.CertName, "cert-name", "", "Public key name in certs path")

	SetFlagsFromEnv({{.Camel}}Cmd)
}

func run{{.TitleCamel}}(cmd *cobra.Command, args []string) {
	configureLogging()
	
	{{.Camel}}.{{.TitleCamel}}({{.Camel}}Cfg)
}
