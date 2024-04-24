package user

import (
	apiError "accumulativeSystem/internal/errors/api"
	"accumulativeSystem/internal/lib/hash"
	balanceModel "accumulativeSystem/internal/models/balance"
	userModel "accumulativeSystem/internal/models/user"
	"accumulativeSystem/internal/repositories/balance"
	userRepository "accumulativeSystem/internal/repositories/user"
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
	"net/http"
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

func NewUserService(repo userRepository.UserRepository, balance balance.BalanceRepository) UserService {
	return &userService{repo: repo, repoBalance: balance}
}

func (s *userService) CreateUser(login, password string) (*userModel.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Хэшируем пароль перед созданием пользователя
	hashPassword, err := hash.HashPassword(password)
	if err != nil {
		return nil, err
	}

	// Начинаем транзакцию
	tx, err := s.repo.Begin(ctx)
	if err != nil {
		return nil, apiError.NewError(http.StatusInternalServerError, "Internal Server Error", err)
	}

	defer func() {
		//TODO интересный момент в случаи паники, err == nil
		if err != nil {
			tx.Rollback(ctx)
		} else {
			err = tx.Commit(ctx)
			if err != nil {
				return
			}
		}
	}()

	//TODO добавить транзакцию Создаем пользователя в рамках транзакции
	err = s.repo.CreateUser(ctx, tx, login, hashPassword)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return nil, apiError.NewError(http.StatusConflict, "duplicate login", pgErr)
			}
		}
		return nil, apiError.NewError(http.StatusInternalServerError, "Internal Server Error", err)
	}

	user, err := s.repo.GetUser(ctx, tx, login)

	if err != nil {
		return nil, apiError.NewError(http.StatusInternalServerError, "Internal Server Error", err)
	}

	var uBalance balanceModel.UserBalance
	uBalance.UserID = user.Id
	uBalance.Balance = 0
	uBalance.WithdrawnBalance = 0

	err = s.repoBalance.CreateBalance(ctx, tx, &uBalance)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) GetUser(login string) (*userModel.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.repo.GetUser(ctx, nil, login)
}

func (s *userService) GetUserWithPassword(login string) (*userModel.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	//TODO возможно ошибка свянная с бд, а не с отсутствующим логином
	//Ошибки - no rows in result set
	//Ошибки - context deadline exceeded

	return s.repo.GetUserWithPassword(ctx, nil, login)
}
