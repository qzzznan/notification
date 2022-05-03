package log

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"os/signal"
	"syscall"
)

const (
	logFilePath = `/var/log/notification/notification.log`
	//logFilePath = `./notification.log`
	pidFilePath = `/var/run/notification.pid`
	//pidFilePath = `./notification.pid`
)

var (
	reopenSig  = make(chan os.Signal, 1)
	logFile    *os.File
	changeList []func(writer io.Writer)

	l *logrus.Logger
)

func init() {
	changeList = make([]func(writer io.Writer), 0)
}

func InitLogConfig() (err error) {
	l = logrus.New()
	l.SetOutput(os.Stderr)

	pid := os.Getpid()
	err = os.WriteFile(pidFilePath, []byte(fmt.Sprintf("%d", pid)), 0644)
	if err != nil {
		return
	}

	logFile, err = os.OpenFile(logFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return
	}

	l.SetLevel(logrus.DebugLevel)
	RegisterResetLogFile(func(w io.Writer) {
		l.SetOutput(w)
	})

	signal.Notify(reopenSig, syscall.SIGUSR1)

	go handleChange()

	return
}

func RegisterResetLogFile(handle func(w io.Writer)) {
	handle(logFile)
	changeList = append(changeList, handle)
}

func handleChange() {
	var err error
	for range reopenSig {
		_ = logFile.Close()
		logFile, err = os.OpenFile(logFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			Errorln("reopen log file failed", err)
			fmt.Println("reopen log file failed", err)
		} else {
			for _, h := range changeList {
				h(logFile)
			}
			Println("reopen log file success")
			fmt.Println("reopen log file success")
		}
	}
}
