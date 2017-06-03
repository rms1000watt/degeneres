package {{.ServiceName.Camel}}

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func {{.TitleCamel}}Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Starting {{.TitleCamel}}Handler...")

}
