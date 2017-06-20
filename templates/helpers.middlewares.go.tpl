package helpers

import (
	"fmt"
	"net/http"
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

func MiddlewareLogging(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
        // NOP for now...
		fn(w, r)
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
		fmt.Println("Origin:", r.Header.Get("Origin"))

		if len(CORSOrigins) == 0 {
			fmt.Println("No CORS Origins defined, but CORS middleware called. No header write.")
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
