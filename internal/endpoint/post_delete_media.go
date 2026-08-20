package endpoint

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/amanagement24/mplay2-go/internal/data"
	"github.com/amanagement24/mplay2-go/internal/service"
)

type PostDeleteMediaEndpoint struct {
	servr         *service.Servr
	uploadsFolder string
}

func NewPostDeleteMediaEndpoint(servr *service.Servr, uploadsFolder string) *PostDeleteMediaEndpoint {
	return &PostDeleteMediaEndpoint{
		servr:         servr,
		uploadsFolder: uploadsFolder,
	}
}

func (e *PostDeleteMediaEndpoint) Handle(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	payload, err := e.process(ctx, r)
	writeJsonResponse(w, payload, err)

	return nil
}

func (e *PostDeleteMediaEndpoint) process(ctx context.Context, r *http.Request) (data.DeleteMediaResponse, error) {
	token, err := getTokenFromRequest(r)
	if err != nil {
		return data.DeleteMediaResponse{}, err
	}

	userID, err := e.servr.ValidateToken(ctx, token)
	if err != nil {
		return data.DeleteMediaResponse{}, err
	}

	var req data.DeleteMediaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return data.DeleteMediaResponse{}, fmt.Errorf("invalid request body: %w", err)
	}

	if len(req.Ids) == 0 {
		return data.DeleteMediaResponse{}, fmt.Errorf("no media ids provided")
	}

	err = e.servr.VerifyMediaOwnership(ctx, userID, req.Ids)
	if err != nil {
		return data.DeleteMediaResponse{}, err
	}

	count, err := e.servr.DeleteMedia(ctx, req.Ids)
	if err != nil {
		return data.DeleteMediaResponse{}, err
	}

	e.deleteMediaFiles(req.Ids)

	return data.DeleteMediaResponse{Deleted: count}, nil
}

func (e *PostDeleteMediaEndpoint) deleteMediaFiles(ids []string) {
	for _, id := range ids {
		filePath := e.uploadsFolder + "/" + id + ".dat"
		_ = os.Remove(filePath)
	}
}
