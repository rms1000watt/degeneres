var (
    {{range $k, $v := .CommandLine.GlobalArgs}}{{$k | ToLower}} {{$v.Type | ToLower}}
    {{end}}
)