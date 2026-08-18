package endpoint

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/amanagement24/mplay2-go/internal/service"
)

type DeleteMediaRequest struct {
	Ids []string `json:"ids"`
}

type DeleteMediaResponse struct {
	Deleted int `json:"deleted"`
}

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

func (e *PostDeleteMediaEndpoint) process(ctx context.Context, r *http.Request) (DeleteMediaResponse, error) {
	token, err := e.getToken(r)
	if err != nil {
		return DeleteMediaResponse{}, err
	}

	userID, err := e.servr.ValidateToken(ctx, token)
	if err != nil {
		return DeleteMediaResponse{}, err
	}

	var req DeleteMediaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return DeleteMediaResponse{}, fmt.Errorf("invalid request body: %w", err)
	}

	if len(req.Ids) == 0 {
		return DeleteMediaResponse{}, fmt.Errorf("no media ids provided")
	}

	err = e.servr.VerifyMediaOwnership(ctx, userID, req.Ids)
	if err != nil {
		return DeleteMediaResponse{}, err
	}

	count, err := e.servr.DeleteMedia(ctx, req.Ids)
	if err != nil {
		return DeleteMediaResponse{}, err
	}

	e.deleteMediaFiles(req.Ids)

	return DeleteMediaResponse{Deleted: count}, nil
}

func (e *PostDeleteMediaEndpoint) getToken(r *http.Request) (string, error) {
	cookie, err := r.Cookie("jtoken12")
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

func (e *PostDeleteMediaEndpoint) deleteMediaFiles(ids []string) {
	for _, id := range ids {
		filePath := e.uploadsFolder + "/" + id + ".dat"
		_ = os.Remove(filePath)
	}
}
