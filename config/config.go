package config

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

const (
	logFilePath = `/var/log/notification.log`
	pidFilePath = `/var/run/notification.pid`
)

var (
	reopenSig = make(chan os.Signal, 1)
	LogFile   *os.File
)

func InitLogConfig() (err error) {
	pid := os.Getpid()
	err = os.WriteFile(pidFilePath, []byte(fmt.Sprintf("%d", pid)), 0644)
	if err != nil {
		return
	}

	LogFile, err = os.OpenFile(logFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return
	}

	log.SetLevel(log.DebugLevel)
	log.SetOutput(LogFile)

	signal.Notify(reopenSig, syscall.SIGUSR1)

	go func() {
		for {
			select {
			case <-reopenSig:
				_ = LogFile.Close()
				LogFile, err = os.OpenFile(logFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
				if err != nil {
					log.Error(err)
				} else {
					log.SetOutput(LogFile)
				}
			}
		}
	}()
	return
}
