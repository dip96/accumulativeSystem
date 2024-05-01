package user

import (
	"github.com/jackc/pgx/v5/pgtype"
	"net/http"
)

type User struct {
	ID           int    `json:"id,omitempty"`
	Login        string `json:"login"`
	Password     string `json:"password,omitempty"`
	HashPassword []byte
	CreatedAt    pgtype.Timestamp `json:"created_at,omitempty"`
}

func (u *User) Bind(r *http.Request) error {
	// доп валидация
	return nil
}
