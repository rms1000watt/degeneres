package helpers

import (
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

var CORSOrigins = []string{ 
    {{.Origins}}
}

// Usage: HandleMiddlewares(PersonHandlerGET, MiddlewareNoCache, MiddlewareCORS)(w, r)
func HandleMiddlewares(handlerFunc http.HandlerFunc, middlewares ...func(http.HandlerFunc) http.HandlerFunc) (h http.HandlerFunc) {
	for _, mw := range middlewares {
		handlerFunc = mw(handlerFunc)
	}
	return handlerFunc
}

func MiddlewareNoCache(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
		fn(w, r)
	}
}

func MiddlewareLogger(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
		fn(w, r)
		log.Debugf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	}
}

func MiddlewareAllows(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		fn(w, r)
	}
}

func MiddlewareCORS(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Debug("CORS: Request Origin:", r.Header.Get("Origin"))

		if len(CORSOrigins) == 0 {
			log.Debug("No CORS Origins defined, but CORS middleware called. No header write.")
			fn(w, r)
			return
		}

		if origin := r.Header.Get("Origin"); origin != "" && valInArr(origin, CORSOrigins) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			fn(w, r)
			return
		}

		w.Header().Set("Access-Control-Allow-Origin", CORSOrigins[0])
		fn(w, r)
	}
}

func valInArr(val string, arr []string) bool {
    for _, a := range arr {
        if val == a {
            return true
        }
    }
    return false
}
