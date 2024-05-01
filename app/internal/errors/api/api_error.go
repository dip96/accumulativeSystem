package api

import "fmt"

type APIError struct {
	Code    int // Код ошибки
	Message string
	Err     error
}

func (e *APIError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s (Code: %d, Error: %v)", e.Message, e.Code, e.Err)
	}
	return fmt.Sprintf("%s (Code: %d)", e.Message, e.Code)
}

func NewError(code int, message string, err error) error {
	return &APIError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}
