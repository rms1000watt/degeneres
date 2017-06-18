package {{.Camel}}

import (
    "fmt"
    "log"
	"net/http"
	"path/filepath"
    
    "{{.ImportPath}}/helpers"
)

func {{.TitleCamel}}(cfg Config) {
    fmt.Println("{{.TitleCamel}} Config:", cfg)
    addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

    if cfg.CertsPath != "" && cfg.CertName != "" && cfg.KeyName != "" {
        fmt.Println("Starting HTTPS server at:", addr)
        log.Fatal(http.ListenAndServeTLS(addr, filepath.Join(cfg.CertsPath, cfg.CertName), filepath.Join(cfg.CertsPath, cfg.KeyName), ServerHandler()))
    } else {
        fmt.Println("Starting HTTP server at:", addr)
        log.Fatal(http.ListenAndServe(addr, ServerHandler()))
    }
}

func ServerHandler() http.Handler {
	mux := http.NewServeMux()

	{{range $endpoint := .Endpoints}}mux.HandleFunc("{{$endpoint.Pattern}}", helpers.HandleMiddlewares({{$endpoint.TitleCamel}}Handler, {{$endpoint.MiddlewareNames}}))
	{{end}}

	return mux
}
