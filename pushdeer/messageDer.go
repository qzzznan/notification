package pushdeer

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func push(c *gin.Context) {
	pushKey := c.GetString("pushkey") // keygen
	text := c.GetString("text")
	_ = pushKey
	_ = text

	c.JSON(http.StatusOK, gin.H{
		"result": []gin.H{
			{"counts": 1, "logs": []string{}, "success": "ok"},
		},
	})
}

func msgList(c *gin.Context) {
	token := c.GetString("token")
	_ = token

	c.JSON(http.StatusOK, gin.H{
		"code": 0, "content": gin.H{
			"message": []gin.H{
				{"id": 1, "text": "hello", "desp": "", "type": "md", "create_at": "time.ASCII"},
			},
		},
	})
}

func msgRemove(c *gin.Context) {
	token := c.GetString("token")
	msgID := c.GetString("id")
	_ = token
	_ = msgID

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"content": gin.H{
			"message": "done",
		},
	})
}
