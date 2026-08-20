package endpoint

import (
	"net/http"
)

type RootEndpoint struct{}

func (e *RootEndpoint) Handle(w http.ResponseWriter, r *http.Request) error {
	payload, err := e.process(r)
	writeJsonResponse(w, payload, err)

	return nil
}

func (e *RootEndpoint) process(_ *http.Request) (string, error) {
	return `{"message":"root endpoint"}`, nil
}
