package pushdeer

import (
	"github.com/gin-gonic/gin"
	"github.com/qzzznan/notification/db"
	"github.com/qzzznan/notification/util"
	"net/http"
)

func merge(c *gin.Context) {
}

func info(c *gin.Context) {
	token := c.Query("token")
	user, err := db.GetUser(token, "")
	if err != nil {
		util.FillRsp(c, 200, 1, err, nil)
		return
	}

	user.Level = 1
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"content": user,
	})
}
