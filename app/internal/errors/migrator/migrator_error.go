package migrator

type MigratorError struct {
	Msg string
	err error
}

func New(text string, err error) error {
	return &MigratorError{text, err}
}

func (e MigratorError) Error() string {
	return e.Msg + ": " + e.err.Error()
}
