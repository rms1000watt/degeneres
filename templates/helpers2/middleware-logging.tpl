func MiddlewareLogging(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
        // NOP for now...
		fn(w, r)
	}
}
