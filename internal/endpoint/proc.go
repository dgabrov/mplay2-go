package endpoint

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
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

func createRandomToken() string {
	u7 := uuid.Must(uuid.NewV7())
	randomBytes := make([]byte, 14)
	_, _ = rand.Read(randomBytes)
	return fmt.Sprintf("%s%s", u7.String(), hex.EncodeToString(randomBytes))[:64]
}

func getTokenFromRequest(r *http.Request) (string, error) {
	cookie, err := r.Cookie("jtoken12")
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}
