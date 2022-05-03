package pushdeer

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/qzzznan/notification/db"
	"github.com/qzzznan/notification/log"
	"github.com/qzzznan/notification/util"
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
	idToken := c.Query("idToken")

	log.Infoln("login with apple id:", idToken)

	//TODO: verify idToken
	token, _ := jwt.Parse(idToken, nil)
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		util.FillRsp(c, http.StatusBadRequest, 1, fmt.Errorf("id token invalid"), nil)
		return
	}
	fields, err := util.GetFields(claims, "sub", "email")
	if err != nil {
		util.FillRsp(c, http.StatusBadRequest, 1, fmt.Errorf("jwt fields error"), nil)
		return
	}
	appleID := fields["sub"]
	email := fields["email"]
	name, _, found := strings.Cut(email, "@")
	if !found {
		name = "unknown"
	}

	uuid, err := db.ExistUser(appleID)
	if err != nil {
		log.Errorln("insert user failed:", err)
		util.FillRsp(c, http.StatusInternalServerError, 1,
			fmt.Errorf("db error"),
			nil)
		return
	}
	if uuid == "" {
		uuid = util.GenUID()
		err = db.InsertUser(appleID, email, name, uuid)
		if err != nil {
			log.Errorln("insert user failed:", err)

			util.FillRsp(c, http.StatusInternalServerError, 1,
				fmt.Errorf("db error"),
				nil)
			return
		}
	}

	util.FillRsp(c, http.StatusOK, 0, nil, gin.H{
		"token": uuid,
	})
}
