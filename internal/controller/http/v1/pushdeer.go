package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"notification/internal/usecase"
	"notification/pkg/logger"
)

type pushDeerRoutes struct {
	p usecase.PushDeer
	l logger.Interface
	v *validator.Validate
}

func newPushDeerRouter(e *gin.RouterGroup, l logger.Interface, p usecase.PushDeer) {
	r := &pushDeerRoutes{p: p, l: l, v: validator.New()}
	g := e.Group("/pushdeer")
	{
		login := g.Group("/login")
		login.Any("/fake", r.fake)
		login.POST("/idtoken", r.apple)
	}
	{
		user := g.Group("/user")
		user.POST("/merge", r.merge)
		user.POST("/info", r.info)
	}
	{
		device := g.Group("/device")
		device.POST("/reg", r.reg)
		device.POST("/list", r.list)
		device.POST("/rename", r.rename)
		device.POST("/remove", r.remove)
	}
	{
		key := g.Group("/key")
		key.POST("/gen", r.gen)
		key.POST("/rename", r.keyRename)
		key.POST("/regen", r.keyRegen)
		key.POST("/list", r.keyList)
		key.POST("/remove", r.keyRemove)
	}
	{
		msg := g.Group("/message")
		msg.POST("/push", r.push)
		msg.POST("/list", r.msgList)
		msg.POST("/remove", r.msgRm)
		msg.POST("/clean", nil)
	}
}
