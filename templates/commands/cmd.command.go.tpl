package cmd

import (
	"{{.ImportPath}}/{{.Camel}}"
	"{{.ImportPath}}/server"
	"github.com/spf13/cobra"
)

var (
	{{.Camel}}Cfg {{.Camel}}.Config
	serverCfg server.Config
)

// {{.Camel}}Cmd represents the {{.Camel}} command
var {{.Camel}}Cmd = &cobra.Command{
	Use:   "{{.Dash}}",
	Short: "{{.ShortDescription}}",
	Long: `{{.LongDescription}}`,
	Run: run{{.TitleCamel}},
}

func init() {
	RootCmd.AddCommand({{.Camel}}Cmd)

	{{.Camel}}Cmd.Flags().StringVar(&serverCfg.Host, "host", "0.0.0.0", "Host address for server")
	{{.Camel}}Cmd.Flags().IntVar(&serverCfg.Port, "port", 8080, "Port for server")

	{{.Camel}}Cmd.Flags().StringVar(&serverCfg.CertsPath, "certs-path", "", "Path for certs")
	{{.Camel}}Cmd.Flags().StringVar(&serverCfg.KeyName, "key-name", "", "Private key name in certs path")
	{{.Camel}}Cmd.Flags().StringVar(&serverCfg.CertName, "cert-name", "", "Public key name in certs path")

	SetFlagsFromEnv({{.Camel}}Cmd)
}

func run{{.TitleCamel}}(cmd *cobra.Command, args []string) {
	configureLogging()
	
	server.{{.TitleCamel}}(serverCfg, {{.Camel}}Cfg)
}
