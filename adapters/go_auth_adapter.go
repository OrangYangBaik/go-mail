package adapters

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type GoAuthAdapter interface {
	RefreshToken(refreshToken, googleId string) error
}

type goAuthAdapter struct{}

func NewGoAuthAdapter() GoAuthAdapter {
	return &goAuthAdapter{}
}

func (a *goAuthAdapter) RefreshToken(refreshToken, googleId string) error {
	url := "http://localhost:8080/auth/google/refresh"
	payload := map[string]string{"refresh_token": refreshToken, "google_id": googleId}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("failed to refresh token")
	}

	return nil
}
