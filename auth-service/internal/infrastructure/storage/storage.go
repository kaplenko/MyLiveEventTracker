package storage

import (
	"auth-service/internal/entity"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	pool *pgxpool.Pool
}

func New(connStr string) (*Storage, error) {
	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return nil, err
	}
	if err = pool.Ping(context.Background()); err != nil {
		return nil, err
	}
	return &Storage{pool: pool}, nil
}

func (s Storage) SaveUser(ctx context.Context, user entity.User, passwordHash []byte) (int64, error) {
	query := `INSERT INTO users (username, email, pass_hash) 
			  VALUES ($1, $2, $3) 
			  RETURNING id`
	var userId int64
	if err := s.pool.QueryRow(ctx, query, user.Username, user.Email, passwordHash).Scan(&userId); err != nil {
		return 0, err
	}
	return userId, nil
}

func (s Storage) GetUserByID(ctx context.Context, id int64) (entity.User, error) {
	query := `SELECT id, username, email, pass_hash 
			  FROM users 
			  WHERE id = $1`
	var user entity.User
	if err := s.pool.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PassHash,
	); err != nil {
		return entity.User{}, err
	}
	return user, nil
}

func (s Storage) GetUserByEmail(ctx context.Context, email string) (entity.User, error) {
	query := `SELECT id, username, email, pass_hash
			  FROM users
			  WHERE email = $1`
	var user entity.User
	if err := s.pool.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PassHash,
	); err != nil {
		return entity.User{}, err
	}
	return user, nil
}
