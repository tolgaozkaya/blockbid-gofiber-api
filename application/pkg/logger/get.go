package logger

import (
	l "github.com/sirupsen/logrus"
)

var Log *l.Logger

type Formatter struct {
	Timestamp bool
	MinLevel  l.Level
}

func init() {
	if Log == nil {
		Init()
	}
}
