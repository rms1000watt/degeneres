package helpers

import (
    "net/http"
    "encoding/json"
    
    log "github.com/sirupsen/logrus"
)

func WriteJSON(out interface{}, w http.ResponseWriter) {
    jsonBytes, err := json.Marshal(out)
    if err != nil {
        log.Error("Failed marshalling to JSON:", err)
        http.Error(w, ErrorJSON("JSON Marshal Error"), http.StatusInternalServerError)
        return
    }

    if _, err := w.Write(jsonBytes); err != nil {
        log.Error("Failed writing to response writer:", err)
        http.Error(w, ErrorJSON("Failed writing to output"), http.StatusInternalServerError)
        return
    }
}

func ErrorJSON(msg string) (out string) {
	return `{"error":"` + msg + `"}`
}