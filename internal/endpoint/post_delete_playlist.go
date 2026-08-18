package endpoint

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/amanagement24/mplay2-go/internal/service"
)

type DeletePlaylistRequest struct {
	Ids []string `json:"ids"`
}

type DeletePlaylistResponse struct {
	Deleted int `json:"deleted"`
}

type PostDeletePlaylistEndpoint struct {
	servr *service.Servr
}

func NewPostDeletePlaylistEndpoint(servr *service.Servr) *PostDeletePlaylistEndpoint {
	return &PostDeletePlaylistEndpoint{
		servr: servr,
	}
}

func (e *PostDeletePlaylistEndpoint) Handle(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	payload, err := e.process(ctx, r)
	writeJsonResponse(w, payload, err)

	return nil
}

func (e *PostDeletePlaylistEndpoint) process(ctx context.Context, r *http.Request) (DeletePlaylistResponse, error) {
	token, err := getTokenFromRequest(r)
	if err != nil {
		return DeletePlaylistResponse{}, err
	}

	userID, err := e.servr.ValidateToken(ctx, token)
	if err != nil {
		return DeletePlaylistResponse{}, err
	}

	var req DeletePlaylistRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return DeletePlaylistResponse{}, fmt.Errorf("invalid request body: %w", err)
	}

	if len(req.Ids) == 0 {
		return DeletePlaylistResponse{}, fmt.Errorf("no playlist ids provided")
	}

	err = e.servr.VerifyPlaylistOwnership(ctx, userID, req.Ids)
	if err != nil {
		return DeletePlaylistResponse{}, err
	}

	count, err := e.servr.DeletePlaylist(ctx, req.Ids)
	if err != nil {
		return DeletePlaylistResponse{}, err
	}

	return DeletePlaylistResponse{Deleted: count}, nil
}
