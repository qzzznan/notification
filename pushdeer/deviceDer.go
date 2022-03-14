package pushdeer

import (
	"github.com/gin-gonic/gin"
	"github.com/qzzznan/notification/model"
	"net/http"
)

func reg(c *gin.Context) {
	regInfo := &model.RegInfo{}
	err := c.Bind(regInfo)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{
			"code": 1, "error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0, "content": gin.H{
			"devices": []gin.H{
				{"id": 2, "uid": 1, "name": "", "type": "all", "device_id": "", "is_clip": 0},
			}, // return all devices
		},
	})
}

func list(c *gin.Context) {
	token, ok := c.GetPostForm("token")
	_ = token
	_ = ok

	c.JSON(http.StatusOK, gin.H{
		"code": 0, "content": gin.H{
			"devices": []gin.H{
				{"id": 1, "uid": 1, "name": "", "type": "all", "device_id": "", "is_clip": 0},
			},
		},
	})
}

func rename(c *gin.Context) {
	token := c.GetString("token")
	id := c.GetString("id")
	name := c.GetString("name")
	_ = token
	_ = id
	_ = name

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
	})
}

func remove(c *gin.Context) {
	token, ok := c.GetPostForm("token")
	id, ok := c.GetPostForm("id")
	_ = token
	_ = ok
	_ = id

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"content": gin.H{
			"message": "done",
		},
	})
}
