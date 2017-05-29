{{if .CommandLine.Command.API}}
func ServerHandler() http.Handler {
	mux := http.NewServeMux()

	{{range $path := .API.Paths}}mux.HandleFunc("{{$path.Pattern}}", HandleMiddlewares({{$path.Name | Title}}Handler{{GetPathMiddlewares $}}))
	{{end}}

	return mux
}
{{end}}


