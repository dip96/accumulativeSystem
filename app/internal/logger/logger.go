package logger

import (
	"accumulativeSystem/internal/config"
	"github.com/sirupsen/logrus"
)

type Logger interface {
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
}

type logger struct {
	log *logrus.Logger
}

func Init(env string) Logger {
	l := logrus.New()
	l.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	switch env {
	case config.EnvLocal:
		l.SetLevel(logrus.DebugLevel)
	case config.EnvTest:
		l.SetLevel(logrus.DebugLevel)
	case config.EnvProd:
		l.SetLevel(logrus.InfoLevel)
	default:
		panic("no stand specified")
	}

	return &logger{log: l}
}

func (l *logger) Debug(args ...interface{}) {
	l.log.Debug(args...)
}

func (l *logger) Debugf(format string, args ...interface{}) {
	l.log.Debugf(format, args...)
}

func (l *logger) Info(args ...interface{}) {
	l.log.Info(args...)
}

func (l *logger) Infof(format string, args ...interface{}) {
	l.log.Infof(format, args...)
}

func (l *logger) Error(args ...interface{}) {
	l.log.Error(args...)
}

func (l *logger) Errorf(format string, args ...interface{}) {
	l.log.Errorf(format, args...)
}
