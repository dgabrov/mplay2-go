package endpoint

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/amanagement24/mplay2-go/internal/data"
	"github.com/amanagement24/mplay2-go/internal/service"
)

type PostUpdatePlaylistEndpoint struct {
	servr *service.Servr
}

func NewPostUpdatePlaylistEndpoint(servr *service.Servr) *PostUpdatePlaylistEndpoint {
	return &PostUpdatePlaylistEndpoint{
		servr: servr,
	}
}

func (e *PostUpdatePlaylistEndpoint) Handle(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	payload, err := e.process(ctx, r)
	writeJsonResponse(w, payload, err)

	return nil
}

func (e *PostUpdatePlaylistEndpoint) process(ctx context.Context, r *http.Request) (SuccessResponse, error) {
	token, err := e.getToken(r)
	if err != nil {
		return SuccessResponse{}, err
	}

	userID, err := e.servr.ValidateToken(ctx, token)
	if err != nil {
		return SuccessResponse{}, err
	}

	var update data.DescriptionUpdate
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		return SuccessResponse{}, fmt.Errorf("invalid request body: %w", err)
	}

	err = e.servr.VerifyPlaylistOwnership(ctx, userID, []string{update.Id})
	if err != nil {
		return SuccessResponse{}, err
	}

	err = e.servr.UpdatePlaylist(ctx, userID, update.Id, update.Description)
	if err != nil {
		return SuccessResponse{}, err
	}

	return SuccessResponse{Success: true}, nil
}

func (e *PostUpdatePlaylistEndpoint) getToken(r *http.Request) (string, error) {
	cookie, err := r.Cookie("jtoken12")
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}
