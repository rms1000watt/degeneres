package {{.CommandLine.Command.Name}}

{{if .CommandLine.Command.API}}
import (
	"encoding/json"
	"fmt"
	"net/http"
)

{{range $path := .API.Paths}}
func {{$path.Name | Title}}Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Starting {{$path.Name | Title}}Handler...")


	// TODO: Add `HandleMiddlewares()()` to handlerFuncs below
	switch r.Method {
	{{range $method := $path.Methods}}case http.{{GetHTTPMethod $method.Name}}:
		HandleMiddlewares({{$path.Name | Title}}Handler{{$method.Name | ToUpper}}{{GetMethodMiddlewares $method.Name $}})(w, r)
	{{end}}case http.MethodOptions:
		HandleMiddlewares({{$path.Name | Title}}HandlerOPTIONS{{GetMethodMiddlewares "options" $}})(w, r)
	default:
		fmt.Println("Method not allowed:", r.Method)
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}

	fmt.Println("Finished {{$path.Name | Title}}Handler!")
}
{{end}}

{{range $path := .API.Paths}}
func {{$path.Name | Title}}HandlerOPTIONS(w http.ResponseWriter, r *http.Request) {
	// Handle Options...
}

{{range $method := $path.Methods}}
func {{$path.Name | Title}}Handler{{$method.Name | ToUpper}}(w http.ResponseWriter, r *http.Request) {
	// Assume JSON Serialization for now
	input := &{{$path.Name | Title}}Input{{$method.Name | ToUpper}}{}
	if err := Unmarshal(r, input); err != nil {
		fmt.Println("Failed decoding input:", err)
		http.Error(w, ErrorJSON("Input Error"), http.StatusInternalServerError)
		return
	}

	msg, err := Validate(input)
	if err != nil {
		fmt.Println("Failed validating input:", err)
		http.Error(w, ErrorJSON("Input Error"), http.StatusInternalServerError)
		return
	}

	if msg != "" {
		fmt.Println("Failed validation:", msg)
		http.Error(w, ErrorJSON("Invalid Input"), http.StatusBadRequest)
		return
	}	

	if err := Transform(input); err != nil {
		fmt.Println("Failed transforming input:", err)
		http.Error(w, ErrorJSON("Transform Error"), http.StatusInternalServerError)
		return
	}

	// Developer make updates here...

	output := get{{$path.Name | Title}}Output{{$method.Name | ToUpper}}(input)

	jsonBytes, err := json.Marshal(output)
	if err != nil {
		fmt.Println("Failed marshalling to JSON:", err)
		http.Error(w, ErrorJSON("JSON Marshal Error"), http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(jsonBytes); err != nil {
		fmt.Println("Failed writing to response writer:", err)
		http.Error(w, ErrorJSON("Failed writing to output"), http.StatusInternalServerError)
		return
	}
}
{{end}}
{{end}}
{{end}}