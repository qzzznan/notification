package bark

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/qzzznan/notification/db"
	"github.com/qzzznan/notification/log"
	"github.com/qzzznan/notification/model"
	"github.com/qzzznan/notification/util"
	"time"
)

func InitBarkHandler(e *gin.Engine) {
	e.GET("/ping", ping)
	e.GET("/health", health)
	e.GET("/info", info)
	e.GET("/register", register)
	e.GET("/:device_key/:content", pushV1)
}

func ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"code":      200,
		"message":   "pong",
		"timestamp": time.Now().Unix(),
	})
}

func health(c *gin.Context) {
	c.String(200, "ok")
}

func info(c *gin.Context) {
	c.JSON(200, gin.H{
		// TODO: add more info
		"version": 0,
		"build":   "",
		"arch":    "",
		"commit":  "",
		"devices": []string{},
	})
}

func register(c *gin.Context) {
	dt := c.Query("devicetoken")
	k := c.Query("key")
	if k != "" && dt == db.GetBarkToken(k) {
		j(c, 200, "success", gin.H{
			"key":          k,
			"device_key":   k,
			"device_token": dt,
		})
		return
	}

	if dt == "" {
		j(c, 400, "device token is required", nil)
		return
	}

	o, err := db.GetBarkDevice("", dt)
	if err != nil && err != sql.ErrNoRows {
		log.Errorf("get device error: %v", err)
		j(c, 400, "get device error", nil)
		return
	}
	if err != sql.ErrNoRows {
		j(c, 200, "success", gin.H{
			"key":          o.DeviceKey,
			"device_key":   o.DeviceKey,
			"device_token": o.DeviceToken,
		})
		return
	}
	if k == "" {
		k = util.GenUID()
	}

	err = db.InsertBarkDevice(k, dt)
	if err != nil {
		log.Errorf("register device error: %v", err)
		j(c, 400, "register device error", nil)
		return
	}

	j(c, 200, "success", gin.H{
		"key":          k,
		"device_key":   k,
		"device_token": dt,
	})
}

func pushV1(c *gin.Context) {
	deviceKey := c.Param("device_key")
	content := c.Param("content")

	m := make(map[string]interface{})
	for k, v := range c.Request.URL.Query() {
		if len(v) != 0 {
			m[k] = v[0]
		}
	}

	token := db.GetBarkToken(deviceKey)
	if token == "" {
		j(c, 400, "device key is not found", nil)
		return
	}

	// TODO: save message
	err := PushMessage(&model.APNsMessage{
		DeviceToken: token,
		Title:       c.Query("title"),
		Category:    c.Query("category"),
		Body:        content,
		Sound:       c.DefaultQuery("sound", "1107"),
	})
	if err != nil {
		j(c, 400, "push message error", nil)
		return
	}
	j(c, 200, "success", nil)
}

func j(c *gin.Context, code int, msg string, data interface{}) {
	h := gin.H{
		"code":      code,
		"message":   msg,
		"timestamp": time.Now().Unix(),
	}
	if data != nil {
		h["data"] = data
	}
	c.JSON(code, h)
}
