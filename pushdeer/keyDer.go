package pushdeer

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func gen(c *gin.Context) {
	token := c.GetString("token")
	_ = token

	c.JSON(http.StatusOK, gin.H{
		"code": 0, "content": gin.H{
			"keys": []gin.H{
				{"id": 1, "uid": 1, "key": "adsafasdfasdf", "name": "adfsafs", "create_at": "adsfa"},
				{"id": 2, "uid": 1, "key": "adsafasdfasdf", "name": "adfsafs", "create_at": "adsfa"},
			},
		},
	})
}

func keyRename(c *gin.Context) {
	token := c.GetString("token")
	id := c.GetString("id")
	newKey := c.GetString("name")
	_ = token
	_ = id
	_ = newKey
}

func keyRegen(c *gin.Context) {
	token := c.GetString("token")
	id := c.GetString("id")
	_ = token
	_ = id

	c.JSON(http.StatusOK, gin.H{
		"code": 0, "content": gin.H{"message": "done"},
	})
}

func keyList(c *gin.Context) {
	token := c.GetString("token")
	_ = token

	c.JSON(http.StatusOK, gin.H{
		"code": 0, "content": gin.H{
			"keys": []gin.H{
				{"id": 1, "uid": 1, "key": "adsafasdfasdf"},
			},
		},
	})
}

func keyRemove(c *gin.Context) {
	token := c.GetString("token")
	id := c.GetString("id")
	_ = token
	_ = id

	c.JSON(http.StatusOK, gin.H{
		"code": 0, "content": gin.H{"message": "done"},
	})
}
