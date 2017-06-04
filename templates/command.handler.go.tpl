package {{.ServiceName.Camel}}

import (
	"encoding/json"
	"fmt"
	"net/http"

	"{{.ImportPath}}/data"
	"{{.ImportPath}}/helpers"
)

func {{.TitleCamel}}Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Starting {{.TitleCamel}}Handler")

	{{if .Methods}}
	switch r.Method {
		{{range $method := .Methods}}case http.Method{{$method.TitleCamel}}:
			helpers.HandleMiddlewares({{$.TitleCamel}}Handler{{$method.UpperCamel}}, {{$.MiddlewareNames}})(w, r)
		{{end}}default:
			fmt.Println("Method not allowed:", r.Method)
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
	{{else}}
	{{template "handler-logic.tpl" .}}
	{{end}}

	fmt.Println("Finished {{.TitleCamel}}Handler")
}

{{range $method := .Methods}}
func {{$.TitleCamel}}Handler{{$method.UpperCamel}}(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Starting {{$.TitleCamel}}Handler{{$method.UpperCamel}}")

	{{template "handler-logic.tpl" $}}

	fmt.Println("Finished {{$.TitleCamel}}Handler{{$method.UpperCamel}}")
}
{{end}}