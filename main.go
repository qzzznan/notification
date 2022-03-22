package main

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/qzzznan/notification/bark"
	"github.com/qzzznan/notification/db"
	"github.com/qzzznan/notification/pushdeer"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io/ioutil"
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

	err = db.InitPostgresDB()
	if err != nil {
		log.Fatalln(err)
	}

	e := gin.Default()
	//e.Use(gin.Recovery())
	//e.Use(ginReqLogMiddleware)
	//e.Use(ginBodyLogMiddleware)

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
