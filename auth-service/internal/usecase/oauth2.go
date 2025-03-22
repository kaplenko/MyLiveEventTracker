package usecase

import (
	"auth-service/internal/entity"
	"auth-service/pkg/jwt"
	"auth-service/pkg/oauth2/github"
	"context"
)

type Oauth2Storage interface {
	UserByGithubID(context.Context, string, int64) (*entity.User, error)
	SaveOauthConnection(context.Context, int64, *github.User) error
}

func (uc *UseCase) AuthenticateOAuthUser(ctx context.Context, oauthUser *github.User) (string, error) {
	existingUser, err := uc.oauthStorage.UserByGithubID(ctx, oauthUser.Provider, oauthUser.ID)
	if err != nil {
		uc.log.ErrorContext(ctx, "failed to find OAuth user", "error", err)
		return "", err
	}
	var userID int64
	if existingUser != nil {
		userID = existingUser.ID
	} else {
		//TODO: исправить логику создания oauth2 user без пароля
		pass, err := github.GenerateState()
		if err != nil {
			uc.log.ErrorContext(ctx, "failed to generate state", "error", err)
			return "", err
		}
		user := &entity.User{
			Username: oauthUser.Login,
			Email:    oauthUser.Email,
			PassHash: []byte(pass),
		}

		userID, err = uc.userSaver.SaveUser(ctx, *user, []byte(pass))
		if err != nil {
			uc.log.ErrorContext(ctx, "failed to register OAuth user", "error", err)
			return "", err
		}
		uc.log.InfoContext(ctx, "registered new OAuth user", "user_id", userID)
	}

	err = uc.oauthStorage.SaveOauthConnection(ctx, userID, oauthUser)
	if err != nil {
		uc.log.ErrorContext(ctx, "failed to save OAuth connection", "error", err)
		return "", err
	}
	token, err := jwt.NewToken(entity.User{ID: userID, Username: oauthUser.Login, Email: oauthUser.Email}, uc.tokenTTL)
	if err != nil {
		uc.log.ErrorContext(ctx, "failed to create token")
		return "", err
	}
	return token, nil
}
