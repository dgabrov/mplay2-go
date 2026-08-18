package endpoint

import (
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /", wrapHandler(&RootEndpoint{}))
}

func wrapHandler(h Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Handle(w, r)
	}
}
