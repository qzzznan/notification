package log

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io"
	"os"
	"os/signal"
	"syscall"
)

var (
	logFilePath string
	pidFilePath string
)

var (
	reopenSig  = make(chan os.Signal, 1)
	logFile    *os.File
	changeList []func(writer io.Writer)

	l     *logrus.Logger
	pidFD *os.File
)

func InitLogConfig() (err error) {
	changeList = make([]func(writer io.Writer), 0)

	m := viper.GetStringMapString("log")
	logFilePath = m["log_file_path"]
	pidFilePath = m["pid_file_path"]

	err = lockPIDFile()
	if err != nil {
		return err
	}

	l = logrus.New()
	l.SetOutput(os.Stderr)

	logFile, err = os.OpenFile(logFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return
	}

	var level logrus.Level
	level, err = logrus.ParseLevel(m["level"])
	if err != nil {
		return err
	}

	l.SetLevel(level)
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
		old := logFile
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
			_ = old.Sync()
			_ = old.Close()
		}
	}
}

func lockPIDFile() error {
	pid := os.Getpid()
	var err error
	pidFD, err = os.OpenFile(pidFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("open pid file failed: %v", err)
	}
	// 建议性锁
	err = syscall.Flock(int(pidFD.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	if err != nil {
		return fmt.Errorf("lock pid file failed: %v", err)
	}
	_, err = pidFD.WriteString(fmt.Sprintf("%d", pid))
	if err != nil {
		return fmt.Errorf("write pid file failed: %v", err)
	}
	err = pidFD.Sync()
	if err != nil {
		return fmt.Errorf("sync pid file failed: %v", err)
	}
}
