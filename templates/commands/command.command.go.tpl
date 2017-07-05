package {{.Camel}}

import (
    "fmt"
	"net/http"
	"path/filepath"
    
    "{{.ImportPath}}/helpers"
    log "github.com/sirupsen/logrus"
)

func {{.TitleCamel}}(cfg Config) {
    log.Debug("{{.TitleCamel}} Config: ", cfg)
    addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

    srv := http.Server{
		Addr:              addr,
		ReadTimeout:       30 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       30 * time.Second,
		Handler:           ServerHandler(),
	}

    if cfg.CertsPath != "" && cfg.CertName != "" && cfg.KeyName != "" {
        log.Info("Starting HTTPS server at: ", addr)
        log.Fatal(srv.ListenAndServeTLS(filepath.Join(cfg.CertsPath, cfg.CertName), filepath.Join(cfg.CertsPath, cfg.KeyName)))
    } else {
        log.Info("Starting HTTP server at: ", addr)
        log.Fatal(srv.ListenAndServe())
    }
}

func ServerHandler() http.Handler {
	mux := http.NewServeMux()

	{{range $endpoint := .Endpoints}}mux.HandleFunc("{{$endpoint.Pattern}}", helpers.HandleMiddlewares({{$endpoint.TitleCamel}}Handler, {{$.MiddlewareNames}}))
	{{end}}
    mux.HandleFunc("/", helpers.HandleMiddlewares(RootHandler, helpers.MiddlewareLogger))

	return mux
}

func RootHandler(w http.ResponseWriter, r *http.Request) {
    log.Debug("Path not found: ", r.URL.Path)
    http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}
