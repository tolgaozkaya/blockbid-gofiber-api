package logger

import (
	"os"

	l "github.com/sirupsen/logrus"
)

func Init() {
	Log = l.New()
	Log.SetReportCaller(true)
	Log.SetFormatter(&Formatter{MinLevel: l.DebugLevel})
	file := file()
	Log.SetLevel(GetLogLevel())
	Log.Out = file
}

func GetLogLevel() l.Level {
	var level l.Level = l.InfoLevel
	var err error
	if env := os.Getenv("LOG_LEVEL"); env != "" {
		level, err = l.ParseLevel(env)
		if err != nil {
			Log.Fatalln(err)
		}
	}
	return level
}
