package endpoint

import (
	"context"
	"net/http"
	"strings"

	"github.com/amanagement24/mplay2-go/internal/data"
	"github.com/amanagement24/mplay2-go/internal/service"
)

type GetSearchPlaylistEndpoint struct {
	servr *service.Servr
}

func NewGetSearchPlaylistEndpoint(servr *service.Servr) *GetSearchPlaylistEndpoint {
	return &GetSearchPlaylistEndpoint{
		servr: servr,
	}
}

func (e *GetSearchPlaylistEndpoint) Handle(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	payload, err := e.process(ctx, r)
	writeJsonResponse(w, payload, err)

	return nil
}

func (e *GetSearchPlaylistEndpoint) process(ctx context.Context, r *http.Request) ([]*data.PlayList, error) {
	token, err := getTokenFromRequest(r)
	if err != nil {
		return []*data.PlayList{}, err
	}

	userID, err := e.servr.ValidateToken(ctx, token)
	if err != nil {
		return []*data.PlayList{}, err
	}

	searchParam := r.URL.Query().Get("searchPlaylist")
	searchTerms := e.transformSearchParam(searchParam)

	playlists, err := e.servr.SearchPlaylist(ctx, userID, searchTerms)
	if err != nil {
		return []*data.PlayList{}, err
	}

	if playlists == nil {
		return []*data.PlayList{}, nil
	}

	return playlists, nil
}

func (e *GetSearchPlaylistEndpoint) transformSearchParam(param string) []string {
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
