package user

import (
	"accumulativeSystem/internal/lib/hash"
	balanceModel "accumulativeSystem/internal/models/balance"
	userModel "accumulativeSystem/internal/models/user"
	"accumulativeSystem/internal/repositories/balance"
	userRepository "accumulativeSystem/internal/repositories/user"
	"context"
	"time"
)

type UserService interface {
	CreateUser(login, password string) (*userModel.User, error)
	GetUser(login string) (*userModel.User, error)
	GetUserWithPassword(login string) (*userModel.User, error)
}
type userService struct {
	repo        userRepository.UserRepository
	repoBalance balance.BalanceRepository
}

func NewUserService(repo userRepository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) CreateUser(login, password string) (*userModel.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Хэшируем пароль перед созданием пользователя
	hashPassword, err := hash.HashPassword(password)
	if err != nil {
		return nil, err
	}

	//TODO добавить транзакцию
	// Создаем пользователя в рамках транзакции
	user, err := s.repo.CreateUser(ctx, login, hashPassword)
	if err != nil {
		return nil, err
	}

	//TODO добавить корректные ошибки с возможность их идентификации в хендлере
	//if err != nil {
	//	var postgresErr *errPostgres.PostgresError
	//	if errors.As(err, &postgresErr) {
	//		http.Error(w, postgresErr.Error(), http.StatusConflict)
	//	} else {
	//		http.Error(w, err.Error(), http.StatusInternalServerError)
	//	}
	//
	//	return
	//}
	//END TODO

	var uBalance balanceModel.UserBalance
	uBalance.UserID = user.Id
	_, err = s.repoBalance.CreateBalance(ctx, &uBalance)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) GetUser(login string) (*userModel.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.repo.GetUser(ctx, login)
}

func (s *userService) GetUserWithPassword(login string) (*userModel.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	//TODO возможно ошибка свянная с бд, а не с отсутствующим логином
	//Ошибки - no rows in result set
	//Ошибки - context deadline exceeded

	return s.repo.GetUserWithPassword(ctx, login)
}
