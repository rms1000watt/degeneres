package {{.ServiceName.Camel}}

import (
	"encoding/json"
	"fmt"
	"net/http"

	"{{.ImportPath}}/data"
	"{{.ImportPath}}/helpers"

	log "github.com/sirupsen/logrus"
)

{{if not .Methods}}
func {{.TitleCamel}}Handler(w http.ResponseWriter, r *http.Request) {
	log.Debug("Starting {{.TitleCamel}}Handler")
	defer log.Debug("Finished {{.TitleCamel}}Handler")

	{{template "handler-logic.tpl" .}}
}
{{end}}

{{range $method := .Methods}}
func {{$.TitleCamel}}Handler{{$method.UpperCamel}}(w http.ResponseWriter, r *http.Request) {
	log.Debug("Starting {{$.TitleCamel}}Handler{{$method.UpperCamel}}")
	defer log.Debug("Finished {{$.TitleCamel}}Handler{{$method.UpperCamel}}")

	{{template "handler-logic.tpl" $}}	
}
{{end}}