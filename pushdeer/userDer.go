package pushdeer

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func merge(c *gin.Context) {

}

func info(c *gin.Context) {
	token, ok := c.GetPostForm("token")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "error": "token is required"})
		return
	}
	_ = token

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"content": gin.H{
			"id":        1,
			"uid":       1,
			"name":      "pushdeer",
			"email":     "",
			"level":     1,
			"apple_id":  "111",
			"create_at": "",
			"update_at": "",
		},
	})
}
