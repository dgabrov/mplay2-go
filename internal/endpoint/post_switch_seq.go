package endpoint

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/amanagement24/mplay2-go/internal/data"
	"github.com/amanagement24/mplay2-go/internal/service"
)

type PostSwitchSeqEndpoint struct {
	servr *service.Servr
}

func NewPostSwitchSeqEndpoint(servr *service.Servr) *PostSwitchSeqEndpoint {
	return &PostSwitchSeqEndpoint{
		servr: servr,
	}
}

func (e *PostSwitchSeqEndpoint) Handle(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	payload, err := e.process(ctx, r)
	writeJsonResponse(w, payload, err)

	return nil
}

func (e *PostSwitchSeqEndpoint) process(ctx context.Context, r *http.Request) (data.SuccessResponse, error) {
	token, err := getTokenFromRequest(r)
	if err != nil {
		return data.SuccessResponse{}, fmt.Errorf("user not logged in")
	}

	userID, err := e.servr.ValidateToken(ctx, token)
	if err != nil {
		return data.SuccessResponse{}, err
	}

	var req data.SwitchSeqRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		return data.SuccessResponse{}, fmt.Errorf("invalid request body: %w", err)
	}

	if err = e.servr.SwitchMediaSequence(ctx, userID, req.PlaylistId, req.Media1, req.Media2); err != nil {
		return data.SuccessResponse{}, err
	}

	return data.SuccessResponse{Success: true}, nil
}
