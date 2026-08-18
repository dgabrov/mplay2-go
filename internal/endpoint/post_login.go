package endpoint

import (
	"encoding/json"
	"net/http"

	"github.com/amanagement24/mplay2-go/internal/data"
)

type PostLoginEndpoint struct{}

func (e *PostLoginEndpoint) Handle(w http.ResponseWriter, r *http.Request) error {
	payload, err := e.process(r)
	writeJsonResponse(w, payload, err)

	return nil
}

func (e *PostLoginEndpoint) process(r *http.Request) (data.LoginResponse, error) {
	var loginData data.LoginData
	if err := json.NewDecoder(r.Body).Decode(&loginData); err != nil {
		return data.LoginResponse{}, err
	}

	return data.LoginResponse{}, nil
}
