package pushdeer

import (
	"github.com/gin-gonic/gin"
	"github.com/qzzznan/notification/db"
	"github.com/qzzznan/notification/log"
	"github.com/qzzznan/notification/model"
	"github.com/qzzznan/notification/util"
	"net/http"
	"strconv"
	"time"
)

func push(c *gin.Context) {
	pushKey := c.Query("pushkey")
	text := c.Query("text")
	desp := c.Query("desp")
	typ := c.Query("type")
	typ = "markdown"

	ki, err := db.GetPushKeyInfo(pushKey)
	if err != nil {
		util.FillRsp(c, 200, 1, err, nil)
		return
	}
	devices, err := db.GetAllDevice(ki.UserID)
	if err != nil {
		util.FillRsp(c, 200, 1, err, nil)
		return
	}

	msg := &model.Message{
		UserID:      ki.UserID,
		Text:        text,
		Type:        typ,
		PushKeyName: ki.Name,
		SendAt:      time.Now(),
		Note:        desp,
	}

	err = PushMessage(msg, devices)
	if err != nil {
		util.FillRsp(c, 200, 1, err, nil)
		return
	}

	err = db.AddMessage(msg)
	if err != nil {
		util.FillRsp(c, 200, 1, err, nil)
		return
	}

	log.Debug("push message to ", pushKey)

	c.JSON(http.StatusOK, gin.H{
		"code": 0, "content": gin.H{
			"result": []string{
				`{"count": 1, "logs": [], "success": "ok"`,
			},
		},
	})
}

func msgList(c *gin.Context) {
	token := c.Query("token")
	limit := c.Query("limit")

	id, err := db.GetUserIDStr(token)
	if err != nil {
		util.FillRsp(c, 200, 1, err, nil)
		return
	}

	lim, err := strconv.ParseInt(limit, 10, 32)
	if err != nil {
		lim = 100
	}

	arr, err := db.GetMessages(id, 0, int(lim))
	if err != nil {
		util.FillRsp(c, 200, 1, err, nil)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0, "content": gin.H{
			"messages": arr,
		},
	})
}

func msgRemove(c *gin.Context) {
	token := c.Query("token")
	msgID := c.Query("id")

	_, err := db.GetUserIDStr(token)
	if err != nil {
		util.FillRsp(c, 200, 1, err, nil)
		return
	}

	err = db.RemoveMessage(msgID)
	if err != nil {
		util.FillRsp(c, 200, 1, err, nil)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"content": gin.H{
			"message": "done",
		},
	})
}
