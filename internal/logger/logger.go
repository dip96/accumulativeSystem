package logger

import (
	"accumulativeSystem/internal/config"
	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func Init(env string) {
	Log = logrus.New()
	Log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	switch env {
	case config.EnvLocal:
		Log.SetLevel(logrus.DebugLevel)
	case config.EnvTest:
		Log.SetLevel(logrus.DebugLevel)
	case config.EnvProd:
		Log.SetLevel(logrus.InfoLevel)
	default:
		panic("no stand specified")
	}

	//logEntry := Log.WithField("test1", "test")
	//Log = logEntry.Logger
}
