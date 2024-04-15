package postgres

import (
	"accumulativeSystem/internal/model/user"
	"context"
	"time"
)

func (s *Postgres) CreateUser(user *user.User) error {
	sqlQuery := "INSERT INTO users (login, password) VALUES ($1,$2)"
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	//TODO не смог разобраться почему получаю ошибку в данном подходе
	//sqlQuery := "INSERT INTO users (login, password) VALUES (@login,@password)"
	//_, err := s.Pool.Exec(ctx, sqlQuery,
	//	sql.Named("login", user.Login),
	//	sql.Named("password", user.HashPassword),
	//)

	_, err := s.Pool.Exec(ctx, sqlQuery,
		user.Login,
		user.HashPassword,
	)

	if err != nil {
		return err
	}

	return nil
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
