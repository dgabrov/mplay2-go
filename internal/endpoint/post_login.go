package endpoint

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"slices"

	"github.com/amanagement24/mplay2-go/internal/data"
	"github.com/amanagement24/mplay2-go/internal/service"
)

type PostLoginEndpoint struct {
	config *data.ConfigData
	servr  *service.Servr
}

func NewPostLoginEndpoint(config *data.ConfigData, servr *service.Servr) *PostLoginEndpoint {
	return &PostLoginEndpoint{
		config: config,
		servr:  servr,
	}
}

func (e *PostLoginEndpoint) Handle(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	result, err := e.process(ctx, r)

	if err == nil && result != nil {
		cookie := &http.Cookie{
			Name:     data.TokenCookieName,
			Value:    result.token,
			HttpOnly: true,
			Path:     "/",
		}
		http.SetCookie(w, cookie)
		writeJsonResponse(w, result.response, err)
	} else {
		writeJsonResponse(w, nil, err)
	}

	return nil
}

type loginResult struct {
	response *data.LoginApiResponse
	token    string
}

func (e *PostLoginEndpoint) process(ctx context.Context, r *http.Request) (*loginResult, error) {
	loginData, err := e.getLoginData(r)
	if err != nil {
		return nil, err
	}

	authResponse, err := e.authenticateWithProvider(ctx, loginData)
	if err != nil {
		return nil, err
	}

	err = e.checkUserRights(authResponse)
	if err != nil {
		return nil, err
	}

	user, err := e.getOrCreateUser(ctx, authResponse)
	if err != nil {
		return nil, err
	}

	token := createRandomToken()
	_, err = e.servr.CreateSession(ctx, user.UserID, token)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return &loginResult{
		response: &data.LoginApiResponse{
			Id:    user.UserID,
			Login: user.Login,
			Name:  user.Name,
		},
		token: token,
	}, nil
}

func (e *PostLoginEndpoint) getLoginData(r *http.Request) (*data.LoginData, error) {
	var loginData data.LoginData
	if err := json.NewDecoder(r.Body).Decode(&loginData); err != nil {
		return nil, fmt.Errorf("invalid login data: %w", err)
	}
	return &loginData, nil
}

func (e *PostLoginEndpoint) authenticateWithProvider(ctx context.Context, loginData *data.LoginData) (*data.LoginResponse, error) {
	body, _ := json.Marshal(loginData)
	req, err := http.NewRequestWithContext(ctx, "POST", e.config.Auth.URL, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to authenticate with provider: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("authentication failed: provider returned status %d", resp.StatusCode)
	}

	var authResponse data.LoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResponse); err != nil {
		return nil, fmt.Errorf("failed to parse auth response: %w", err)
	}

	return &authResponse, nil
}

func (e *PostLoginEndpoint) checkUserRights(authResponse *data.LoginResponse) error {
	if slices.Contains(authResponse.Rights, e.config.Auth.Right) {
		return nil
	}
	return fmt.Errorf("you don't have the right to access this application")
}

func (e *PostLoginEndpoint) getOrCreateUser(ctx context.Context, authResponse *data.LoginResponse) (*data.User, error) {
	user, err := e.servr.GetUserByProvidedId(ctx, authResponse.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to check user: %w", err)
	}

	if user != nil {
		slog.Info("User found", "user_id", user.UserID)
		return user, nil
	}

	newUser := &data.User{
		ProvidedUserID: authResponse.Id,
		Login:          authResponse.Login,
		Name:           authResponse.Name,
	}

	err = e.servr.CreateUser(ctx, newUser)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	slog.Info("User created", "user_id", newUser.UserID)
	return newUser, nil
}
