package log

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"os/signal"
	"sync"
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

	l     *logrus.Logger
	pidFD *os.File
)

func init() {
	changeList = make([]func(writer io.Writer), 0)
	LockPIDFile()
}

func InitLogConfig() (err error) {
	l = logrus.New()
	l.SetOutput(os.Stderr)

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

func LockPIDFile() {
	pid := os.Getpid()
	var err error
	pidFD, err = os.OpenFile(pidFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		panic(fmt.Sprintf("open pid file failed: %v", err))
	}
	_, err = pidFD.WriteString(fmt.Sprintf("%d", pid))
	if err != nil {
		panic(fmt.Sprintf("write pid file failed: %v", err))
	}
	err = syscall.Flock(int(pidFD.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	if err != nil {
		panic(fmt.Sprintf("lock pid file failed: %v", err))
	}
}

type LoggerWriter struct {
	io.Writer
	sync.Mutex
}

func (l *LoggerWriter) SetWriter(w io.Writer) {
	l.Lock()
	defer l.Unlock()
	l.Writer = w
}
