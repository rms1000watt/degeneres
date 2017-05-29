    {{range $k, $v := .Args}}{{$k | Title}}: {{$k | ToLower}},
    {{end}}