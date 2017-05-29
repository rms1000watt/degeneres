package {{.CommandLine.Command.Name}}

type Config struct {
    {{range $k, $v := .CommandLine.Command.Args}}{{$k | Title}} {{$v.Type | ToLower}}
    {{end}}
}
