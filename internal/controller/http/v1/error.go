package v1

import (
	"github.com/gin-gonic/gin"
	"notification/pkg/logger"
	"runtime"
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

func pdResp(c *gin.Context, l logger.Interface, state, code int, err error, content interface{}) {
	if err != nil {
		_, fl, ln, _ := runtime.Caller(1)
		l.Errorf("%s:%d %s\n", fl, ln, err.Error())
	}
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
