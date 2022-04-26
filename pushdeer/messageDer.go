package pushdeer

import (
	"github.com/gin-gonic/gin"
	"github.com/qzzznan/notification/db"
	"github.com/qzzznan/notification/model"
	"net/http"
)

func push(c *gin.Context) {
	pushKey := c.Query("pushkey")
	text := c.Query("text")
	desp := c.Query("desp")
	typ := c.Query("type")
	_ = text
	_ = desp
	_ = typ

	ki, err := db.GetPushKeyInfo(pushKey)
	if err != nil {

	}

	err = db.AddMessage(&model.Message{
		PushKeyName: ki.Name,
	})
	if err != nil {

	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0, "content": gin.H{
			"result": []string{
				`{"count": 1, "logs": [], "success": "ok"`,
			},
		},
	})
}

func msgList(c *gin.Context) {
	token := c.GetString("token")
	limit := c.Query("limit")
	_ = token
	_ = limit

	c.JSON(http.StatusOK, gin.H{
		"code": 0, "content": gin.H{
			"messages": []gin.H{
				{
					"id":           1,
					"uid":          "114",
					"text":         "hello",
					"desp":         "",
					"type":         "markdown",
					"created_at":   "2022-04-26T15:49:12.111111Z",
					"pushkey_name": "iphone",
				},
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
