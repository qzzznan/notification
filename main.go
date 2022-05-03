package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/qzzznan/notification/bark"
	"github.com/qzzznan/notification/db"
	"github.com/qzzznan/notification/log"
	"github.com/qzzznan/notification/pushdeer"
	"io"
	"io/ioutil"
	"os"
)

var (
	ip   string
	port string
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
	defer flag.Parse()
	flag.StringVar(&ip, "addr", "localhost", "ip address")
	flag.StringVar(&port, "port", "8080", "port")
}

func main() {
	handleError(
		log.InitLogConfig,
		db.InitPostgresDB,
		bark.InitPushClient,
		func() error {
			return pushdeer.InitPushClient("./static/c.p12")
		},
	)

	var err error
	e := gin.New()

	lc := gin.LoggerConfig{}
	log.RegisterResetLogFile(func(w io.Writer) {
		lc.Output = w
	})

	e.Use(gin.Recovery())
	e.Use(gin.LoggerWithConfig(lc))
	err = e.SetTrustedProxies(nil)
	if err != nil {
		log.Fatalln(err)
	}

	//e.Use(ginReqLogMiddleware)
	//e.Use(ginBodyLogMiddleware)

	pushdeer.InitDerHandler(e)
	bark.InitBarkHandler(e)

	addr := fmt.Sprintf("%s:%s", ip, port)

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
			log.Errorln("index:", i, err)
			fmt.Println("index:", i, err)
			os.Exit(1)
		}
	}
}
