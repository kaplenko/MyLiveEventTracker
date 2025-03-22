package storage

import (
	"auth-service/internal/entity"
	"auth-service/pkg/oauth2/github"
	"context"
	"github.com/jackc/pgx/v5"
)

func (s *Storage) UserByGithubID(ctx context.Context, provider string, providerID int64) (*entity.User, error) {
	var user entity.User
	query := `SELECT u.id, u.username, u.email
			  FROM users u
			  JOIN oauth_connections o ON u.id = o.user_id
			  WHERE o.provider = $1 AND o.provider_id = $2`
	err := s.pool.QueryRow(ctx, query, provider, providerID).Scan(&user.ID, &user.Username, &user.Email)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (s *Storage) SaveOauthConnection(ctx context.Context, userID int64, user *github.User) error {
	query := `INSERT INTO oauth_connections (user_id, provider, provider_id, access_token, refresh_token, expires_at)
			  VALUES ($1, $2, $3, $4, $5, to_timestamp($6))
			  ON CONFLICT (provider, provider_id) DO UPDATE
			  SET access_token = EXCLUDED.access_token,
				refresh_token = EXCLUDED.refresh_token,
				expires_at = EXCLUDED.expires_at`
	_, err := s.pool.Exec(ctx, query, userID, user.Provider, user.ID, user.AccessToken, user.RefreshToken, user.ExpiresAt)
	return err
}
