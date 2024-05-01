package user

import (
	APIError "accumulativeSystem/internal/errors/api"
	"accumulativeSystem/internal/lib/hash"
	"accumulativeSystem/internal/logger"
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
	logger      logger.Logger
}

func NewUserService(repo userRepository.UserRepository, balance balance.BalanceRepository, logger logger.Logger) UserService {
	return &userService{repo: repo, repoBalance: balance, logger: logger}
}

func (s *userService) CreateUser(login, password string) (*userModel.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Хэшируем пароль перед созданием пользователя
	hashPassword, err := hash.HashPassword(password)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}

	// Начинаем транзакцию
	tx, err := s.repo.Begin(ctx)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, APIError.NewError(http.StatusInternalServerError, "Internal Server Error", err)
	}

	defer func() {
		//TODO интересный момент в случаи паники, err == nil
		if err != nil {
			s.logger.Error(err.Error())
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
				s.logger.Error(err.Error())
				return nil, APIError.NewError(http.StatusConflict, "duplicate login", pgErr)
			}
		}
		s.logger.Error(err.Error())
		return nil, APIError.NewError(http.StatusInternalServerError, "Internal Server Error", err)
	}

	user, err := s.repo.GetUser(ctx, tx, login)

	if err != nil {
		s.logger.Error(err.Error())
		return nil, APIError.NewError(http.StatusInternalServerError, "Internal Server Error", err)
	}

	var uBalance balanceModel.UserBalance
	uBalance.UserID = user.Id
	uBalance.Balance = 0
	uBalance.WithdrawnBalance = 0

	err = s.repoBalance.CreateBalance(ctx, tx, &uBalance)

	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}

	return user, nil
}

func (s *userService) GetUser(login string) (*userModel.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user, err := s.repo.GetUser(ctx, nil, login)

	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}

	return user, nil
}

func (s *userService) GetUserWithPassword(login string) (*userModel.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user, err := s.repo.GetUserWithPassword(ctx, nil, login)

	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}

	return user, nil
}
