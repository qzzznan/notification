package pushdeer

import "github.com/gin-gonic/gin"

func InitDerHandler(e *gin.Engine) {
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
