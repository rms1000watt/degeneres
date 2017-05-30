package {{.Camel}}

import (
    "fmt"
    "log"
	"net/http"
	"path/filepath"
)

func {{.TitleCamel}}(cfg Config) {
    fmt.Println("{{.TitleCamel}} Config:", cfg)
    addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

    fmt.Println("Starting server at:", addr)
    {{if .CertsPath}}log.Fatal(http.ListenAndServeTLS(addr, filepath.Join(cfg.CertsPath, cfg.CertName), filepath.Join(cfg.CertsPath, cfg.KeyName), ServerHandler()))
    {{else}}log.Fatal(http.ListenAndServe(addr, ServerHandler()))
    {{end}}
}

func ServerHandler() http.Handler {
	mux := http.NewServeMux()

	{{range $path := .Endpoints}}mux.HandleFunc("{{$path.Pattern}}", HandleMiddlewares({{$path.Name}}Handler{{GetPathMiddlewares $}}))
	{{end}}

	return mux
}
