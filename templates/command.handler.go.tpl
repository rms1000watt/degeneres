package {{.ServiceName.Camel}}

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func {{.TitleCamel}}Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Starting {{.TitleCamel}}Handler...")

	{{if .MethodCheck}}if !({{.MethodCheck}}) {
		fmt.Println(http.StatusText(http.StatusMethodNotAllowed))
		http.Error(w, ErrorJSON(http.StatusText(http.StatusMethodNotAllowed)), http.StatusMethodNotAllowed)
		return
	}{{end}}

}
