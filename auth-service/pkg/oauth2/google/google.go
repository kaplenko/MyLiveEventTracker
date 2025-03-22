package google

import (
	"context"
	"encoding/json"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"log/slog"
)

type Service struct {
	oauth2Cfg *oauth2.Config
	log       *slog.Logger
}

type User struct {
	ID           string `json:"id"`
	Login        string `json:"name"`
	Email        string `json:"email"`
	Provider     string `json:"provider"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresAt    int64  `json:"expiresAt"`
}

func NewGoogleService(ClientID, ClientSecret, RedirectURL string, log *slog.Logger) *Service {
	return &Service{
		oauth2Cfg: &oauth2.Config{
			ClientID:     ClientID,
			ClientSecret: ClientSecret,
			RedirectURL:  RedirectURL,
			Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
			Endpoint:     google.Endpoint,
		},
		log: log,
	}
}

func (s *Service) AuthURL(state string) string {
	return s.oauth2Cfg.AuthCodeURL(state)
}

func (s *Service) UserInfo(ctx context.Context, code string) (*User, error) {
	token, err := s.oauth2Cfg.Exchange(ctx, code)
	if err != nil {
		s.log.ErrorContext(ctx, "Error exchanging code for token: %v", err)
		return nil, err
	}

	s.log.Info("Token obtained successfully, fetching user info")

	client := s.oauth2Cfg.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		s.log.ErrorContext(ctx, "Error getting user info: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	var user User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		s.log.ErrorContext(ctx, "Error decoding user info: %v", err)
		return nil, err
	}

	user.Provider = "google"
	user.AccessToken = token.AccessToken
	user.RefreshToken = token.RefreshToken
	if token.Expiry.IsZero() {
		user.ExpiresAt = 0
	} else {
		user.ExpiresAt = token.Expiry.Unix()
	}
	s.log.Info("Successfully fetched user info: %s", user.Email)
	return &user, nil
}
