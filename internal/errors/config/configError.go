package config

type ConfigError struct {
	msg   string
	error error
}

func New(text string, err error) error {
	return &ConfigError{text, err}
}

func (e *ConfigError) Error() string {
	return e.msg
}
