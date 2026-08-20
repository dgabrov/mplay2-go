package endpoint

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/amanagement24/mplay2-go/internal/data"
	"github.com/amanagement24/mplay2-go/internal/service"
)

type PostAddMediaToPlaylistEndpoint struct {
	servr *service.Servr
}

func NewPostAddMediaToPlaylistEndpoint(servr *service.Servr) *PostAddMediaToPlaylistEndpoint {
	return &PostAddMediaToPlaylistEndpoint{
		servr: servr,
	}
}

func (e *PostAddMediaToPlaylistEndpoint) Handle(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	payload, err := e.process(ctx, r)
	writeJsonResponse(w, payload, err)

	return nil
}

func (e *PostAddMediaToPlaylistEndpoint) process(ctx context.Context, r *http.Request) (data.SuccessResponse, error) {
	token, err := getTokenFromRequest(r)
	if err != nil {
		return data.SuccessResponse{}, err
	}

	userID, err := e.servr.ValidateToken(ctx, token)
	if err != nil {
		return data.SuccessResponse{}, err
	}

	var req data.RemoveMediaFromPlaylistRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return data.SuccessResponse{}, fmt.Errorf("invalid request body: %w", err)
	}

	if req.PlaylistId == "" {
		return data.SuccessResponse{}, fmt.Errorf("playlistId is required")
	}

	if len(req.Ids) == 0 {
		return data.SuccessResponse{}, fmt.Errorf("no media ids provided")
	}

	_, err = e.servr.AddMediaToPlaylist(ctx, userID, req.PlaylistId, req.Ids)
	if err != nil {
		return data.SuccessResponse{}, err
	}

	return data.SuccessResponse{Success: true}, nil
}
