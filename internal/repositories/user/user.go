package user

import (
	userModel "accumulativeSystem/internal/models/user"
	storage "accumulativeSystem/internal/storage/postgres"
	"context"
	"github.com/jackc/pgx/v5"
	"time"
)

type UserRepository interface {
	CreateUser(ctx context.Context, login string, password []byte) (*userModel.User, error)
	GetUser(ctx context.Context, login string) (*userModel.User, error)
	GetUserWithPassword(ctx context.Context, login string) (*userModel.User, error)
	Begin(ctx context.Context) (pgx.Tx, error)
}

type userRepository struct {
	db storage.Storage
}

func NewUserRepository(storage storage.Storage) UserRepository {
	return &userRepository{db: storage}
}

func (r *userRepository) CreateUser(ctx context.Context, login string, password []byte) (*userModel.User, error) {
	sqlQuery := "INSERT INTO users (login, password) VALUES ($1,$2)"
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.db.Exec(ctx, sqlQuery,
		login,
		password,
	)

	if err != nil {
		return nil, err
	}

	user, err := r.GetUser(ctx, login)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepository) GetUser(ctx context.Context, login string) (*userModel.User, error) {
	sqlSelect := "SELECT id, login, created_at FROM users WHERE login = $1"
	var user userModel.User
	err := r.db.QueryRow(ctx, sqlSelect, login).Scan(&user.Id, &user.Login, &user.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetUserWithPassword(ctx context.Context, login string) (*userModel.User, error) {
	sqlSelect := "SELECT id, login, password FROM users WHERE login = $1"
	var user userModel.User
	err := r.db.QueryRow(ctx, sqlSelect, login).Scan(&user.Id, &user.Login, &user.HashPassword)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) Begin(ctx context.Context) (pgx.Tx, error) {
	return r.db.Begin(ctx)
}
