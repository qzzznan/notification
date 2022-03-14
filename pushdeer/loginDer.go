package pushdeer

import (
	"github.com/gin-gonic/gin"
	"github.com/qzzznan/notification/util"
	log "github.com/sirupsen/logrus"
	"net/http"
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

	var userToken string

	userToken = util.GenToken(token)

	log.Infoln("apple userToken:", userToken)

	c.JSON(http.StatusOK, gin.H{
		"code": 0, "content": gin.H{
			"token": userToken,
		},
	})
}
