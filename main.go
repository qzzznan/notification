package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/qzzznan/notification/bark"
	"github.com/qzzznan/notification/pushdeer"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	log.SetLevel(log.DebugLevel)

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalln(err)
	}
}

func main() {
	err := bark.InitPushClient()
	if err != nil {
		log.Fatalln(err)
	}

	e := gin.Default()
	pushdeer.InitDerHandler(e)

	addr := fmt.Sprintf("%s:%s",
		viper.GetString("ip"),
		viper.GetString("port"))

	log.Infof("server start at %s", addr)

	if err = e.Run(addr); err != nil {
		log.Fatalln(err)
	}
}
