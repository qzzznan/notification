package pushdeer

import (
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gin-gonic/gin"
	"github.com/qzzznan/notification/db"
	"github.com/qzzznan/notification/model"
	"github.com/qzzznan/notification/util"
	"net/http"
	"strconv"
)

func gen(c *gin.Context) {
	token := c.Query("token")
	id, err := db.GetUserID(token)
	if err != nil {
		util.FillRsp(c, 400, 1, fmt.Errorf("invalid token"), nil)
		return
	}
	name := c.Query("name")
	if name == "" {
		name = gofakeit.PetName()
	}

	key := gofakeit.LetterN(64)
	err = db.InsertPushKey(&model.PushKey{
		UserID: fmt.Sprintf("%d", id),
		Key:    key,
		Name:   "",
	})
	if err != nil {
		util.FillRsp(c, 400, 1, err, nil)
		return
	}

	keys, err := db.GetAllPushKey(id)
	if err != nil {
		util.FillRsp(c, 400, 1, err, nil)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0, "content": gin.H{
			"keys": keys,
		},
	})
}

func keyRename(c *gin.Context) {
	token := c.Query("token")
	kidStr := c.Query("id")
	newName := c.Query("name")

	_, err := db.GetUserID(token)
	if err != nil {
		util.FillRsp(c, 400, 1, err, nil)
		return
	}

	kid, err := strconv.ParseInt(kidStr, 10, 64)
	if err != nil {
		util.FillRsp(c, 400, 1, err, nil)
		return
	}

	err = db.UpdatePushKey(kid, newName, "")
	if err != nil {
		util.FillRsp(c, 400, 1, err, nil)
		return
	}

	c.JSON(200, gin.H{
		"code": 0, "content": gin.H{
			"message": "done",
		},
	})
}

func keyRegen(c *gin.Context) {
	token := c.Query("token")
	k := c.Query("id")

	_, err := db.GetUserID(token)
	if err != nil {
		util.FillRsp(c, 400, 1, err, nil)
		return
	}

	kid, err := strconv.ParseInt(k, 10, 64)
	if err != nil {
		util.FillRsp(c, 400, 1, err, nil)
		return
	}

	ks := gofakeit.LetterN(64)
	err = db.UpdatePushKey(kid, "", ks)
	if err != nil {
		util.FillRsp(c, 400, 1, err, nil)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0, "content": gin.H{"message": "done"},
	})
}

func keyList(c *gin.Context) {
	token := c.Query("token")
	uid, err := db.GetUserID(token)
	if err != nil {
		util.FillRsp(c, 400, 1, err, nil)
		return
	}

	keys, err := db.GetAllPushKey(uid)
	if err != nil {
		util.FillRsp(c, 400, 1, err, nil)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0, "content": gin.H{
			"keys": keys,
		},
	})
}

func keyRemove(c *gin.Context) {
	token := c.Query("token")
	id := c.Query("id")
	_ = token
	_ = id

	//TODO: remove key

	c.JSON(http.StatusOK, gin.H{
		"code": 0, "content": gin.H{"message": "done"},
	})
}
