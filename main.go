package main

import (
	"bytes"
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
	e.Use(gin.Recovery())
	e.Use(ginBodyLogMiddleware)

	pushdeer.InitDerHandler(e)

	addr := fmt.Sprintf("%s:%s",
		viper.GetString("ip"),
		viper.GetString("port"))

	log.Infof("server start at %s", addr)

	if err = e.Run(addr); err != nil {
		log.Fatalln(err)
	}
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func ginBodyLogMiddleware(c *gin.Context) {
	blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
	c.Writer = blw
	c.Next()
	fmt.Println("Response body: " + blw.body.String())
}
