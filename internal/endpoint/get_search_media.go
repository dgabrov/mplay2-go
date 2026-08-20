package endpoint

import (
	"context"
	"net/http"
	"strings"

	"github.com/amanagement24/mplay2-go/internal/data"
	"github.com/amanagement24/mplay2-go/internal/service"
)

type GetSearchMediaEndpoint struct {
	servr *service.Servr
}

func NewGetSearchMediaEndpoint(servr *service.Servr) *GetSearchMediaEndpoint {
	return &GetSearchMediaEndpoint{
		servr: servr,
	}
}

func (e *GetSearchMediaEndpoint) Handle(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	payload, err := e.process(ctx, r)
	writeJsonResponse(w, payload, err)

	return nil
}

func (e *GetSearchMediaEndpoint) process(ctx context.Context, r *http.Request) ([]*data.Media, error) {
	token, err := getTokenFromRequest(r)
	if err != nil {
		return nil, err
	}

	userID, err := e.servr.ValidateToken(ctx, token)
	if err != nil {
		return nil, err
	}

	searchParam := r.URL.Query().Get("searchMedia")
	searchTerms := e.transformSearchParam(searchParam)

	media, err := e.servr.SearchMedia(ctx, userID, searchTerms)
	if err != nil {
		return nil, err
	}

	return media, nil
}

func (e *GetSearchMediaEndpoint) transformSearchParam(param string) []string {
	// Replace ? with _ and * with %
	transformed := strings.ReplaceAll(param, "?", "_")
	transformed = strings.ReplaceAll(transformed, "*", "%")

	// Split by whitespace
	words := strings.Fields(transformed)

	// Add % at beginning and end if not already there
	for i, word := range words {
		if !strings.HasPrefix(word, "%") {
			word = "%" + word
		}
		if !strings.HasSuffix(word, "%") {
			word = word + "%"
		}
		words[i] = word
	}

	if len(words) == 0 {
		words = append(words, "%")
	}

	return words
}
