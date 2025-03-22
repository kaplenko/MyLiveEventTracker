package github

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"log/slog"
)

type Service struct {
	oauth2Cfg *oauth2.Config
	log       *slog.Logger
}

type User struct {
	ID           int64  `json:"id"`
	Login        string `json:"login"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	Provider     string `json:"provider"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresAt    int64  `json:"expiresAt"`
}

func NewGithubService(ClientID, ClientSecret, RedirectURL string, log *slog.Logger) *Service {
	return &Service{
		oauth2Cfg: &oauth2.Config{
			ClientID:     ClientID,
			ClientSecret: ClientSecret,
			RedirectURL:  RedirectURL,
			Scopes:       []string{"read:user", "user:email"},
			Endpoint:     github.Endpoint,
		},
		log: log,
	}
}
func (g *Service) AuthURL(state string) string {
	return g.oauth2Cfg.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

func (g *Service) UserInfo(ctx context.Context, code string) (*User, error) {
	g.log.Info("Code is: ", code)
	token, err := g.oauth2Cfg.Exchange(ctx, code)
	if err != nil {
		g.log.ErrorContext(ctx, "Error exchanging code for token: %v", err)
		return nil, err
	}

	g.log.Info("Token obtained successfully, fetching user info")
	client := g.oauth2Cfg.Client(ctx, token)

	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		g.log.ErrorContext(ctx, "Error fetching user info: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	var user User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		g.log.ErrorContext(ctx, "Error decoding user info: %v", err)
		return nil, err
	}

	user.Provider = "github"
	user.AccessToken = token.AccessToken
	user.RefreshToken = token.RefreshToken
	if token.Expiry.IsZero() {
		user.ExpiresAt = 0
	} else {
		user.ExpiresAt = token.Expiry.Unix()
	}

	g.log.Info("Successfully fetched user info: %s", user.Login)
	return &user, nil
}

func GenerateState() (string, error) {
	state := make([]byte, 32)
	if _, err := rand.Read(state); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(state), nil
}
