package pushdeer

import (
	"github.com/gin-gonic/gin"
	"github.com/qzzznan/notification/db"
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

	userToken, err := db.ExistUser("", token)
	if err != nil {
		log.Errorln("insert user failed:", err)
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	if userToken == "" {
		userToken = util.GenUID()
		err = db.InsertUser(userToken, token)
		if err != nil {
			log.Errorln("insert user failed:", err)
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}
	}

	log.Infoln("apple userToken:", userToken)

	c.JSON(http.StatusOK, gin.H{
		"code": 0, "content": gin.H{
			"token": userToken,
		},
	})
}
