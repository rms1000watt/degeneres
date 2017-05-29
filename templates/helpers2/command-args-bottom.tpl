{{range $k, $v := .Args}}{{$.Name}}Cmd.Flags().{{$v.Type | Title}}Var(&{{$k | ToLower}}, "{{$k | ToLower}}", {{HandleQuotes $v.Default $v.Type}} ,"{{$v.Description}}")
{{end}}