package usecase

import (
	"auth-service/internal/entity"
	"auth-service/pkg/jwt"
	"context"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"time"
)

type UserSaver interface {
	SaveUser(ctx context.Context, user entity.User, passwordHash []byte) (int64, error)
}

type UseProvider interface {
	GetUserByID(ctx context.Context, id int64) (entity.User, error)
	GetUserByEmail(ctx context.Context, email string) (entity.User, error)
}

type UseCase struct {
	userSaver    UserSaver
	userProvider UseProvider
	log          *slog.Logger
	tokenTTL     time.Duration
}

func New(us UserSaver, up UseProvider, log *slog.Logger, tokenTTL time.Duration) *UseCase {
	return &UseCase{
		userSaver:    us,
		userProvider: up,
		log:          log,
		tokenTTL:     tokenTTL,
	}
}

func (uc *UseCase) GetUser(ctx context.Context, id int64) (entity.User, error) {
	user, err := uc.userProvider.GetUserByID(ctx, id)
	if err != nil {
		uc.log.ErrorContext(ctx, "failed to get user", "id", id, "error", err)
		return entity.User{}, err
	}
	uc.log.InfoContext(ctx, "got user", "id", id, "user", user)
	return user, nil
}

func (uc *UseCase) Registre(ctx context.Context, user entity.User) (int64, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword(user.PassHash, bcrypt.DefaultCost)
	if err != nil {
		uc.log.ErrorContext(ctx, "failed to hash password")
		return 0, err
	}

	id, err := uc.userSaver.SaveUser(ctx, user, hashedPassword)
	if err != nil {
		uc.log.ErrorContext(ctx, "failed to create user", "id", id, "error", err)
		return 0, err
	}

	uc.log.InfoContext(ctx, "created user", "id", id, "user", user)
	return id, nil
}

func (uc *UseCase) Login(ctx context.Context, email string, password []byte) (string, error) {
	user, err := uc.userProvider.GetUserByEmail(ctx, email)
	if err != nil {
		uc.log.ErrorContext(ctx, "failed to get user", "error", err)
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, password); err != nil {
		uc.log.ErrorContext(ctx, "failed to compare password", "error", err)
		return "", err
	}
	uc.log.InfoContext(ctx, "user authenticated", "user", user)

	token, err := jwt.NewToken(user, uc.tokenTTL)
	if err != nil {
		uc.log.ErrorContext(ctx, "failed to create token")
		return "", err
	}
	return token, nil
}
