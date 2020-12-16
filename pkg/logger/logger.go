package logger

import (
	"io"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	*logrus.Logger
}

func NewLogger(level logrus.Level, out io.Writer) *Logger {
	l := &Logger{logrus.New()}
	l.Level = level
	l.Out = out

	return l
}
