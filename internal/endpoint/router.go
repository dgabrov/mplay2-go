package endpoint

import (
	"fmt"
	"github.com/amanagement24/mplay2-go/internal/data"
	"github.com/amanagement24/mplay2-go/internal/service"
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux, context string, config *data.ConfigData, servr *service.Servr) {
	mux.HandleFunc(fmt.Sprintf("GET %s/", context), wrapHandler(&RootEndpoint{}))
	mux.HandleFunc(fmt.Sprintf("POST %s/login", context), wrapHandler(NewPostLoginEndpoint(config, servr)))
	mux.HandleFunc(fmt.Sprintf("GET %s/searchMedia", context), wrapHandler(NewGetSearchMediaEndpoint(servr)))
	mux.HandleFunc(fmt.Sprintf("GET %s/searchPlaylist", context), wrapHandler(NewGetSearchPlaylistEndpoint(servr)))
	mux.HandleFunc(fmt.Sprintf("POST %s/deleteMedia", context), wrapHandler(NewPostDeleteMediaEndpoint(servr, config.UploadsFolder)))
	mux.HandleFunc(fmt.Sprintf("POST %s/deletePlaylist", context), wrapHandler(NewPostDeletePlaylistEndpoint(servr)))
}

func wrapHandler(h Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Handle(w, r)
	}
}
