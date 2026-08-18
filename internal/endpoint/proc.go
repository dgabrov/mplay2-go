package endpoint

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type ErrorResponse struct {
	Message string   `json:"message"`
	Items   []string `json:"items"`
}

func writeJsonResponse(writer http.ResponseWriter, payload any, err error) {
	status := http.StatusOK
	processPayload := payload

	if err != nil {
		processPayload = ErrorResponse{
			Message: err.Error(),
			Items:   make([]string, 0),
		}

		status = http.StatusBadRequest
	}

	header := writer.Header()
	header.Set("Content-Type", "application/json")
	header.Set("Cache-Control", "no-cache")

	writer.WriteHeader(status)

	// write the actual object here
	bts, err := json.Marshal(processPayload)
	if err != nil {
		slog.Error("error marshalling payload")
	}

	_, err = writer.Write(bts)
	if err != nil {
		slog.Error("error writing payload")
	}
}
