
func ServerHandler() http.Handler {
	mux := http.NewServeMux()

	{{range $path := .Endpoints}}mux.HandleFunc("{{$path.Pattern}}", HandleMiddlewares({{$path.Name}}Handler{{GetPathMiddlewares $}}))
	{{end}}

	return mux
}
