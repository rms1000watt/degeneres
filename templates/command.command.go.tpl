package {{.CommandLine.Command.Name}}

import (
    "fmt"
)

func {{.CommandLine.Command.Name | Title}}(cfg Config) {
    fmt.Println("{{.CommandLine.Command.Name | Title}} Config:", cfg)
    {{template "api-middle.tpl" .}}
    {{template "version-middle.tpl" .}}
}

{{template "api-bottom.tpl" .}}
