package v1

import (
	"github.com/gin-gonic/gin"
	"time"
)

func barkResp(c *gin.Context, code int, msg string, data interface{}) {
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

func pdResp(c *gin.Context, state, code int, err error, content interface{}) {
	m := gin.H{
		"code": code,
	}
	if err != nil {
		m["error"] = err
	} else if content != nil {
		m["content"] = content
	}
	c.JSON(state, m)
}
