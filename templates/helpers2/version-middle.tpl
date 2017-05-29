{{if eq .CommandLine.Command.Name "version"}}
fmt.Println("Version: {{.Version}}"){{end}}