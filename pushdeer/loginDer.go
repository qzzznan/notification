package pushdeer

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/qzzznan/notification/util"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

func fake(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": 0, "content": gin.H{
			"token": "HelloWorld",
		},
	})
}

func apple(c *gin.Context) {
	token := c.Query("idToken")

	log.Infoln("apple token:", token)

	userToken, err := DB.GetUserToken(token)
	if strings.Contains(err.Error(), "no rows in result set") {
		userToken = util.GenUID()
		err = DB.SaveLoginToken(token, userToken)
		if err != nil {
			c.JSON(1, gin.H{
				"error": fmt.Sprintf("apple id login %v", err),
			})
			return
		}
	} else {
		c.JSON(1, gin.H{
			"error": fmt.Sprintf("apple id login %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0, "content": gin.H{
			"token": userToken,
		},
	})
}
