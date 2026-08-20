package endpoint

import (
	"context"
	"fmt"
	"net/http"

	"github.com/amanagement24/mplay2-go/internal/data"
	"github.com/amanagement24/mplay2-go/internal/service"
)

type GetMediaForPlaylistEndpoint struct {
	servr *service.Servr
}

func NewGetMediaForPlaylistEndpoint(servr *service.Servr) *GetMediaForPlaylistEndpoint {
	return &GetMediaForPlaylistEndpoint{
		servr: servr,
	}
}

func (e *GetMediaForPlaylistEndpoint) Handle(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	payload, err := e.process(ctx, r)
	writeJsonResponse(w, payload, err)

	return nil
}

func (e *GetMediaForPlaylistEndpoint) process(ctx context.Context, r *http.Request) ([]*data.Media, error) {
	token, err := getTokenFromRequest(r)
	if err != nil {
		return nil, err
	}

	userID, err := e.servr.ValidateToken(ctx, token)
	if err != nil {
		return nil, err
	}

	playlistID := r.URL.Query().Get("playlistId")
	if playlistID == "" {
		return nil, fmt.Errorf("playlistId query parameter is required")
	}

	if err := e.servr.VerifyPlaylistOwnership(ctx, userID, []string{playlistID}); err != nil {
		return nil, err
	}

	media, err := e.servr.GetMediaForPlaylist(ctx, userID, playlistID)
	if err != nil {
		return nil, err
	}

	if media == nil {
		media = make([]*data.Media, 0)
	}

	return media, nil
}
