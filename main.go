package main

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/qzzznan/notification/bark"
	"github.com/qzzznan/notification/config"
	"github.com/qzzznan/notification/db"
	"github.com/qzzznan/notification/pushdeer"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
)

func init() {
	/*
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		if err := viper.ReadInConfig(); err != nil {
			log.Fatalln(err)
		}
	*/
}

func main() {
	handleError(
		config.InitLogConfig,
		db.InitPostgresDB,
		bark.InitPushClient,
		func() error {
			return pushdeer.InitPushClient("./static/c.p12")
		},
	)

	var err error
	e := gin.Default()
	err = e.SetTrustedProxies(nil)
	if err != nil {
		log.Fatalln(err)
	}
	e.Use(gin.LoggerWithWriter(config.LogFile))

	//e.Use(ginReqLogMiddleware)
	//e.Use(ginBodyLogMiddleware)

	pushdeer.InitDerHandler(e)
	bark.InitBarkHandler(e)

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
	log.Debugln("Response body: " + blw.body.String())
}

func ginReqLogMiddleware(c *gin.Context) {
	defer c.Next()

	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return
	}
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(data))
	log.Debugln("Request Header:", c.Request.Header)
	log.Debugln("Request URL:", c.Request.RequestURI)
	log.Debugln("Request Body:", string(data))
}

type handleF = func() error

func handleError(fl ...handleF) {
	for i, f := range fl {
		if err := f(); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "index:", i, err)
			os.Exit(1)
		}
	}
}
