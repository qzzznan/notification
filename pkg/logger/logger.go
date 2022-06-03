package logger

import "github.com/sirupsen/logrus"

var _ Interface = (*Logger)(nil)

type Interface interface {
	logrus.FieldLogger
}

type Logger struct {
	*logrus.Logger
}

func New(level string) (*Logger, error) {
	l, err := logrus.ParseLevel(level)
	if err != nil {
		return nil, err
	}
	logger := logrus.New()
	logger.SetLevel(l)
	logger.ReportCaller = true
	logger.SetFormatter(&logrus.TextFormatter{
		DisableColors: false,
		ForceColors:   true,
	})
	//logger.AddHook()
	return &Logger{logger}, nil
}
