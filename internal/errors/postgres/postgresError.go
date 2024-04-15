package postgres

type PostgresError struct {
	msg   string
	error error
}

func New(text string, err error) error {
	return &PostgresError{text, err}
}

func (e *PostgresError) Error() string {
	return e.msg
}
