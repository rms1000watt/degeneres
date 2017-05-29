{{range $k, $v := .CommandLine.GlobalArgs}}RootCmd.PersistentFlags().{{$v.Type | Title}}Var(&{{$k | ToLower}}, "{{$k | ToLower}}", {{HandleQuotes $v.Default $v.Type}} ,"{{$v.Description}}")
{{end}}