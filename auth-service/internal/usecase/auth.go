package usecase

import (
	"auth-service/internal/entity"
	"auth-service/pkg/jwt"
	"context"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"time"
)

type Repository interface {
	GetByID(ctx context.Context, id int64) (entity.User, error)
	SaveUser(ctx context.Context, user entity.User, passwordHash string) (int64, error)
	GetByEmail(ctx context.Context, email string) (entity.User, error)
}

type UseCase struct {
	repo     Repository
	lg       *slog.Logger
	tokenTTL time.Duration
}

func New(r Repository) *UseCase {
	return &UseCase{repo: r}
}

func (uc *UseCase) GetUser(ctx context.Context, id int64) (entity.User, error) {
	user, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		uc.lg.ErrorContext(ctx, "failed to get user", "id", id, "error", err)
		return user, err
	}
	uc.lg.InfoContext(ctx, "got user", "id", id, "user", user)
	return user, nil
}

func (uc *UseCase) Registre(ctx context.Context, user entity.User) (int64, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		uc.lg.ErrorContext(ctx, "failed to hash password")
		return 0, err
	}

	id, err := uc.repo.SaveUser(ctx, user, string(hashedPassword))
	if err != nil {
		uc.lg.ErrorContext(ctx, "failed to create user", "id", user.ID, "error", err)
		return 0, err
	}

	uc.lg.InfoContext(ctx, "created user", "id", user.ID, "user", user)
	return id, nil
}

func (uc *UseCase) Login(ctx context.Context, email, password string) (string, error) {
	user, err := uc.repo.GetByEmail(ctx, email)
	if err != nil {
		uc.lg.ErrorContext(ctx, "failed to get user", "email", email, "error", err)
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		uc.lg.ErrorContext(ctx, "failed to compare password", "email", email, "error", err)
		return "", err
	}
	uc.lg.InfoContext(ctx, "user authenticated", "email", email, "user", user)

	token, err := jwt.NewToken(user, uc.tokenTTL)
	if err != nil {
		uc.lg.ErrorContext(ctx, "failed to create token")
		return "", err
	}
	return token, nil
}
