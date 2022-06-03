package v1

import (
	"github.com/gin-gonic/gin"
	"notification/internal/entity"
	"notification/internal/usecase"
	"notification/pkg/logger"
	"time"
)

type barkRoutes struct {
	b usecase.Bark
	l logger.Interface
}

func newBarkRouter(g *gin.RouterGroup, l logger.Interface, b usecase.Bark) {
	e := g.Group("/bark")
	r := &barkRoutes{b, l}
	e.GET("/ping", r.ping)
	e.GET("/health", r.health)
	e.GET("/info", r.info)
	e.GET("/register", r.register)
	e.GET("/:key/:content", r.push)
}

func (r *barkRoutes) ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"code":      200,
		"message":   "pong",
		"timestamp": time.Now().Unix(),
	})
}

func (r *barkRoutes) health(c *gin.Context) {
	c.String(200, "ok")
}

func (r *barkRoutes) info(c *gin.Context) {
	c.JSON(200, gin.H{
		// TODO: add more info
		"version": 0,
		"build":   "",
		"arch":    "",
		"commit":  "",
		"devices": []string{},
	})
}

func (r *barkRoutes) register(c *gin.Context) {
	token := c.Query("devicetoken")
	key := c.Query("key")

	ent := &entity.BarkDevice{
		DeviceToken: token,
		DeviceKey:   key,
	}
	err := r.b.Register(c.Request.Context(), ent)
	if err != nil {
		r.l.Errorln("register bark device", err)
		barkResp(c, 400, "internal error", nil)
		return
	}
	barkResp(c, 200, "success", gin.H{
		"key":          ent.DeviceKey,
		"device_key":   ent.DeviceKey,
		"device_token": ent.DeviceToken,
	})
}

func (r *barkRoutes) push(c *gin.Context) {
	key := c.Param("key")

	m := make(map[string]interface{})
	for k, v := range c.Request.URL.Query() {
		if len(v) != 0 {
			m[k] = v[0]
		}
	}

	msg := &entity.APNsMessage{
		Title:    c.Query("title"),
		Category: c.Query("category"),
		Body:     c.Param("content"),
		Sound:    c.DefaultQuery("sound", "1107"),
		Data:     m,
	}

	err := r.b.Push(c.Request.Context(), key, msg)
	if err != nil {
		r.l.Errorln("push bark message", err)
		barkResp(c, 400, "internal error", nil)
	}
	barkResp(c, 200, "success", nil)
}
