package pushdeer

import (
	"context"
	"github.com/gin-gonic/gin"
	db2 "github.com/qzzznan/notification/db"
	log "github.com/sirupsen/logrus"
)

var DB db2.DeerDB

func InitDerHandler(e *gin.Engine) {
	// TODO: 移到其他地方去
	s := &db2.PostgresDB{}
	err := s.Init(context.Background())
	if err != nil {
		log.Fatalln(err)
	}
	DB = s
	// -----

	{
		login := e.Group("/login")
		login.GET("/fake", fake)
		login.POST("/idtoken", apple)
	}
	{
		user := e.Group("/user")
		user.POST("/merge", merge)
		user.POST("/info", info)
	}
	{
		device := e.Group("/device")
		device.POST("/reg", reg)
		device.POST("/list", list)
		device.POST("/rename", rename)
		device.POST("/remove", remove)
	}
	{
		key := e.Group("/key")
		key.POST("/gen", gen)
		key.POST("/rename", keyRename)
		key.POST("/regen", keyRegen)
		key.POST("/list", keyList)
		key.POST("/remove", keyRemove)
	}
	{
		msg := e.Group("/message")
		msg.POST("/push", push)
		msg.POST("/list", msgList)
		msg.POST("/remove", msgRemove)
	}
}
