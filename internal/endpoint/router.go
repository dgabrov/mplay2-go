package endpoint

import (
	"fmt"
	"github.com/amanagement24/mplay2-go/internal/data"
	"github.com/amanagement24/mplay2-go/internal/service"
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux, context string, config *data.ConfigData, servr *service.Servr) {
	mux.HandleFunc(fmt.Sprintf("GET %s/", context), wrapHandler(&RootEndpoint{}))
	mux.HandleFunc(fmt.Sprintf("GET %s/playMedia", context), wrapHandler(NewGetPlayMediaEndpoint(servr, config.UploadsFolder, config.MediaSlice)))
	mux.HandleFunc(fmt.Sprintf("POST %s/login", context), wrapHandler(NewPostLoginEndpoint(config, servr)))
	mux.HandleFunc(fmt.Sprintf("GET %s/searchMedia", context), wrapHandler(NewGetSearchMediaEndpoint(servr)))
	mux.HandleFunc(fmt.Sprintf("GET %s/searchPlaylist", context), wrapHandler(NewGetSearchPlaylistEndpoint(servr)))
	mux.HandleFunc(fmt.Sprintf("GET %s/getMediaForPlaylist", context), wrapHandler(NewGetMediaForPlaylistEndpoint(servr)))
	mux.HandleFunc(fmt.Sprintf("POST %s/deleteMedia", context), wrapHandler(NewPostDeleteMediaEndpoint(servr, config.UploadsFolder)))
	mux.HandleFunc(fmt.Sprintf("POST %s/updateMedia", context), wrapHandler(NewPostUpdateMediaEndpoint(servr, config.UploadsFolder)))
	mux.HandleFunc(fmt.Sprintf("POST %s/deletePlaylist", context), wrapHandler(NewPostDeletePlaylistEndpoint(servr)))
	mux.HandleFunc(fmt.Sprintf("POST %s/addPlaylist", context), wrapHandler(NewPostAddPlaylistEndpoint(servr)))
	mux.HandleFunc(fmt.Sprintf("POST %s/updatePlaylist", context), wrapHandler(NewPostUpdatePlaylistEndpoint(servr)))
	mux.HandleFunc(fmt.Sprintf("POST %s/removeMediaFromPlaylist", context), wrapHandler(NewPostRemoveMediaFromPlaylistEndpoint(servr)))
	mux.HandleFunc(fmt.Sprintf("POST %s/addMediaToPlaylist", context), wrapHandler(NewPostAddMediaToPlaylistEndpoint(servr)))
}

func wrapHandler(h Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Handle(w, r)
	}
}
