package postgres

import (
	userModel "accumulativeSystem/internal/models/user"
	"context"
	"github.com/jackc/pgx/v5"
	"time"
)

func (s *Postgres) CreateUser(login string, password []byte) (*userModel.User, error) {
	//начинаем транзакцию
	tx, err := s.Pool.Begin(context.Background())
	if err != nil {
		return nil, err
	}

	defer func(tx pgx.Tx, ctx context.Context) {
		err := tx.Rollback(ctx)
		if err != nil {
			//TODO добавить логи
		}
	}(tx, context.Background())

	sqlQuery := "INSERT INTO users (login, password) VALUES ($1,$2)"
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	//TODO не смог разобраться почему получаю ошибку в данном подходе
	//sqlQuery := "INSERT INTO users (login, password) VALUES (@login,@password)"
	//_, err := s.Pool.Exec(ctx, sqlQuery,
	//	sql.Named("login", user.Login),
	//	sql.Named("password", user.HashPassword),
	//)

	_, err = s.Pool.Exec(ctx, sqlQuery,
		login,
		password,
	)

	if err != nil {
		return nil, err
	}

	user, err := s.GetUser(login)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Postgres) GetUser(login string) (*userModel.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	sqlSelect := "SELECT id, login, created_at FROM users WHERE login = $1"
	var user userModel.User
	err := s.Pool.QueryRow(ctx, sqlSelect, login).Scan(&user.Id, &user.Login, &user.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *Postgres) GetUserPassword(login string) (*userModel.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	sqlSelect := "SELECT id, login, password FROM users WHERE login = $1"
	var user userModel.User
	err := s.Pool.QueryRow(ctx, sqlSelect, login).Scan(&user.Id, &user.Login, &user.HashPassword)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

//func (s *Postgres) checkUniqueLogin(login string) error {
//	sqlQuery := "SELECT login FROM users WHERE login = $1"
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//
//	_, err := s.Pool.Exec(ctx, sqlQuery, login)
//
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
