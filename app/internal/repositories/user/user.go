package user

import (
	userModel "accumulativeSystem/internal/models/user"
	"accumulativeSystem/internal/storage"
	"context"
	"github.com/jackc/pgx/v5"
)

type UserRepository interface {
	CreateUser(ctx context.Context, tx pgx.Tx, login string, password []byte) error
	GetUser(ctx context.Context, tx pgx.Tx, login string) (*userModel.User, error)
	GetUserWithPassword(ctx context.Context, tx pgx.Tx, login string) (*userModel.User, error)
	Begin(ctx context.Context) (pgx.Tx, error)
}

// TODO добавить интерфейся для работы только с транзакциями??? Например
//type UserTransactionalRepository interface {
//	CreateUserTx(ctx context.Context, tx pgx.Tx, login string, password []byte) error
//	GetUserTx(ctx context.Context, tx pgx.Tx, login string) (*userModel.User, error)
//	GetUserWithPasswordTx(ctx context.Context, tx pgx.Tx, login string) (*userModel.User, error)
//}

type userRepository struct {
	db storage.Storage
}

func NewUserRepository(storage storage.Storage) UserRepository {
	return &userRepository{db: storage}
}

func (r *userRepository) CreateUser(ctx context.Context, tx pgx.Tx, login string, password []byte) error {
	sqlQuery := "INSERT INTO users (login, password) VALUES ($1,$2)"

	if tx == nil {
		_, err := r.db.Exec(ctx, sqlQuery, login, password)
		return err
	}

	_, err := tx.Exec(ctx, sqlQuery, login, password)

	return err
}

func (r *userRepository) GetUser(ctx context.Context, tx pgx.Tx, login string) (*userModel.User, error) {
	sqlSelect := "SELECT id, login, created_at FROM users WHERE login = $1"
	var user userModel.User

	if tx == nil {
		err := r.db.QueryRow(ctx, sqlSelect, login).Scan(&user.ID, &user.Login, &user.CreatedAt)
		if err != nil {
			return nil, err
		}
		return &user, nil
	}

	err := tx.QueryRow(ctx, sqlSelect, login).Scan(&user.ID, &user.Login, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetUserWithPassword(ctx context.Context, tx pgx.Tx, login string) (*userModel.User, error) {
	sqlSelect := "SELECT id, login, password FROM users WHERE login = $1"
	var user userModel.User
	err := r.db.QueryRow(ctx, sqlSelect, login).Scan(&user.ID, &user.Login, &user.HashPassword)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) Begin(ctx context.Context) (pgx.Tx, error) {
	return r.db.Begin(ctx)
}
