package server

import (
    "fmt"
	"net/http"
	"path/filepath"
    
    "{{.ImportPath}}/{{.Camel}}"
    "{{.ImportPath}}/helpers"
    log "github.com/sirupsen/logrus"
)

func {{.TitleCamel}}(cfg Config, {{.Camel}}Cfg {{.Camel}}.Config) {
    log.Debug("{{.TitleCamel}} Config: ", cfg)
    addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	if err := {{.Camel}}.PreServe({{.Camel}}Cfg); err != nil {
		log.Error("Failed running preserve function: ", err)
		return
	}

    srv := http.Server{
		Addr:              addr,
		ReadTimeout:       30 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       30 * time.Second,
		Handler:           {{.TitleCamel}}ServerHandler(),
	}

    if cfg.CertsPath != "" && cfg.CertName != "" && cfg.KeyName != "" {
        log.Info("Starting HTTPS server at: ", addr)
        log.Fatal(srv.ListenAndServeTLS(filepath.Join(cfg.CertsPath, cfg.CertName), filepath.Join(cfg.CertsPath, cfg.KeyName)))
    } else {
        log.Info("Starting HTTP server at: ", addr)
        log.Fatal(srv.ListenAndServe())
    }
}

func {{.TitleCamel}}ServerHandler() http.Handler {
	mux := http.NewServeMux()

	{{range $endpoint := .Endpoints}}mux.HandleFunc("{{$endpoint.Pattern}}", helpers.HandleMiddlewares({{$endpoint.Camel}}Handler, {{$.MiddlewareNames}}))
	{{end}}mux.HandleFunc("/", helpers.HandleMiddlewares(RootHandler, helpers.MiddlewareLogger))

	return mux
}

{{range $endpoint := .Endpoints}}
func {{$endpoint.Camel}}Handler(w http.ResponseWriter, r *http.Request) {
	log.Debug("Starting {{$endpoint.Camel}}Handler")
	defer log.Debug("Finished {{$endpoint.Camel}}Handler")

	if r.Method == http.MethodOptions {
		// Process headers and return
		return
	}

	{{if $endpoint.Methods}}
	switch r.Method {
		{{range $method := $endpoint.Methods}}case http.Method{{$method.TitleCamel}}:
			helpers.HandleMiddlewares({{$.Camel}}.{{$endpoint.TitleCamel}}Handler{{$method.UpperCamel}}, {{$endpoint.MiddlewareNames}})(w, r)
		{{end}}default:
			log.Debug("Method not allowed: ", r.Method)
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
	{{else}}
    helpers.HandleMiddlewares({{$.Camel}}.{{$endpoint.TitleCamel}}Handler, {{$endpoint.MiddlewareNames}})(w, r)
	{{end}}
}
{{end}}

func RootHandler(w http.ResponseWriter, r *http.Request) {
    log.Debug("Path not found: ", r.URL.Path)
    http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}
