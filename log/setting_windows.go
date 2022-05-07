package log

import (
	"github.com/sirupsen/logrus"
	"io"
	"os"
)

var (
	l *logrus.Logger
)

func InitLogConfig() error {
	l = logrus.New()
	l.SetOutput(os.Stderr)
	return nil
}

func RegisterResetLogFile(handler func(w io.Writer)) {
	handler(os.Stderr)
}
