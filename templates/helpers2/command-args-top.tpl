var (
    {{range $k, $v := .Args}}{{$k | ToLower}} {{$v.Type | ToLower}}
    {{end}}
)