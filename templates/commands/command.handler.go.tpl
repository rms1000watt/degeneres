package {{.ServiceName.Camel}}

import (
	"encoding/json"
	"fmt"
	"net/http"

	"{{.ImportPath}}/data"
	"{{.ImportPath}}/helpers"

	log "github.com/sirupsen/logrus"
)

func {{.TitleCamel}}Handler(w http.ResponseWriter, r *http.Request) {
	log.Debug("Starting {{.TitleCamel}}Handler")
	defer log.Debug("Finished {{.TitleCamel}}Handler")

	if r.Method == http.MethodOptions {
		// Process headers and return
		return
	}

	{{if .Methods}}
	switch r.Method {
		{{range $method := .Methods}}case http.Method{{$method.TitleCamel}}:
			helpers.HandleMiddlewares({{$.TitleCamel}}Handler{{$method.UpperCamel}}, {{$.MiddlewareNames}})(w, r)
		{{end}}default:
			log.Debug("Method not allowed: ", r.Method)
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
	{{else}}
	{{template "handler-logic.tpl" .}}
	{{end}}
}

{{range $method := .Methods}}
func {{$.TitleCamel}}Handler{{$method.UpperCamel}}(w http.ResponseWriter, r *http.Request) {
	log.Debug("Starting {{$.TitleCamel}}Handler{{$method.UpperCamel}}")
	defer log.Debug("Finished {{$.TitleCamel}}Handler{{$method.UpperCamel}}")

	{{template "handler-logic.tpl" $}}	
}
{{end}}