package endpoint

import (
	"fmt"
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux, context string) {
	mux.HandleFunc(fmt.Sprintf("GET %s/", context), wrapHandler(&RootEndpoint{}))
	mux.HandleFunc(fmt.Sprintf("POST %s/login", context), wrapHandler(&PostLoginEndpoint{}))
}

func wrapHandler(h Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Handle(w, r)
	}
}
