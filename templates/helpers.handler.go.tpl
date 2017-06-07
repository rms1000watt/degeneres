package helpers

import (
    "fmt"
    "net/http"
    "encoding/json"
)

func WriteJSON(out interface{}, w http.ResponseWriter) {
    jsonBytes, err := json.Marshal(out)
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

func ErrorJSON(msg string) (out string) {
	return `{"error":"` + msg + `"}`
}