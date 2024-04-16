package user

//import (
//	handRegistration "accumulativeSystem/internal/http-server/handlers/user/registration"
//	"accumulativeSystem/internal/lib/hash"
//	"accumulativeSystem/internal/storage/postgres"
//	"net/http"
//)
//
//type Service struct {
//	postgres *postgres.Postgres
//}
//
//func (s Service) CreateUser(req handRegistration.Request) {
//	// Проверка на уникальность логина
//	if err := s.checkUniqueLogin(login); err != nil {
//		return nil, err
//	}
//
//	hashPassword, err := hash.HashPassword(req.Password)
//
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusBadRequest)
//		return
//	}
//
//	// Создание пользователя
//	user, err := s.postgres.CreateUser(login, password)
//	if err != nil {
//		return nil, err
//	}
//
//	return user, nil
//}
