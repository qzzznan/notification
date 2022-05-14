package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/qzzznan/notification/bark"
	"github.com/qzzznan/notification/config"
	"github.com/qzzznan/notification/db"
	"github.com/qzzznan/notification/log"
	"github.com/qzzznan/notification/pushdeer"
	"github.com/qzzznan/notification/ws"
	"github.com/spf13/viper"
	"io"
	"os"
)

func init() {
	config.LoadConfig()

	handleError(
		log.InitLogConfig,
		db.InitPostgresDB,
		db.InitRedis,
		db.InitBoltDB,
		bark.InitPushClient,
		pushdeer.InitPushClient,
	)
}

func main() {
	var err error
	if !viper.GetBool("http.debug") {
		gin.SetMode(gin.ReleaseMode)
	}
	e := gin.New()

	lw := &log.LoggerWriter{}
	lc := gin.LoggerConfig{
		Output: lw,
	}
	log.RegisterResetLogFile(func(w io.Writer) {
		lw.SetWriter(w)
	})

	e.Use(gin.Recovery())
	e.Use(gin.LoggerWithConfig(lc))
	err = e.SetTrustedProxies(nil)
	if err != nil {
		log.Fatalln(err)
	}

	pushdeer.InitDerHandler(e)
	bark.InitBarkHandler(e)
	ws.InitWsHandler(e)

	addr := viper.GetString("http.addr")
	log.Infof("server start at %s", addr)

	if err = e.Run(addr); err != nil {
		log.Fatalln(err)
	}
}

type handleF = func() error

func handleError(fl ...handleF) {
	for i, f := range fl {
		if err := f(); err != nil {
			log.Errorln("index:", i, err)
			fmt.Println("index:", i, err)
			os.Exit(1)
		}
	}
}
