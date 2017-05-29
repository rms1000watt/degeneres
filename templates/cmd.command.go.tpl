package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	
)

{{template "command-args-top.tpl" .}}

// {{.Name}}Cmd represents the {{.Name}} command
var {{.Name}}Cmd = &cobra.Command{
	Use:   "{{.Name}}",
	Short: "{{.ShortDescription}}",
	Long: `{{.LongDescription}}`,
	Run: Run{{.Name | Title}},
}

func init() {
	RootCmd.AddCommand({{.Name}}Cmd)

	{{template "command-args-bottom.tpl" .}}

	SetFlagsFromEnv({{.Name}}Cmd)
}

func Run{{.Name | Title}}(cmd *cobra.Command, args []string) {
	// Get config arguments and pass it to the function itself
	{{.Name}}Cfg := {{.Name}}.Config{
		{{template "command-args-middle.tpl" .}}
	}

	{{.Name}}.{{.Name | Title}}({{.Name}}Cfg)
}
