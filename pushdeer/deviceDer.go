package pushdeer

import (
	"github.com/gin-gonic/gin"
	"github.com/qzzznan/notification/db"
	"github.com/qzzznan/notification/model"
	"github.com/qzzznan/notification/util"
	"net/http"
	"strconv"
)

func reg(c *gin.Context) {
	regInfo := &model.RegInfo{}
	err := c.BindQuery(&regInfo)
	if err != nil {
		util.FillRsp(c, 200, 1, err, nil)
		return
	}

	err = util.Validate.Struct(regInfo)
	if err != nil {
		util.FillRsp(c, 200, 1, err, nil)
		return
	}
	id, err := db.GetUserIDStr(regInfo.Token)
	if err != nil {
		util.FillRsp(c, http.StatusForbidden, 1, err, nil)
		return
	}

	err = db.InsertDevice(&model.Device{
		UserID:   id,
		DeviceID: regInfo.DeviceID,
		Type:     "ios",
		IsClip:   regInfo.IsClip,
		Name:     regInfo.Name,
	})
	if err != nil {
		util.FillRsp(c, http.StatusForbidden, 1, err, nil)
		return
	}
	all, err := db.GetAllDevice(id)
	if err != nil {
		util.FillRsp(c, http.StatusForbidden, 1, err, nil)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0, "content": gin.H{
			"devices": all,
		},
	})
}

func list(c *gin.Context) {
	token := c.Query("token")
	id, err := db.GetUserIDStr(token)
	if err != nil {
		util.FillRsp(c, http.StatusForbidden, 1, err, nil)
		return
	}
	all, err := db.GetAllDevice(id)
	if err != nil {
		util.FillRsp(c, http.StatusForbidden, 1, err, nil)
		return
	}

	util.FillRsp(c, 200, 0, nil, gin.H{
		"devices": all,
	})
}

func rename(c *gin.Context) {
	token := c.Query("token")
	id := c.Query("id")
	name := c.Query("name")

	_, err := db.GetUserIDStr(token)
	if err != nil {
		util.FillRsp(c, http.StatusForbidden, 1, err, nil)
		return
	}

	deviceID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		util.FillRsp(c, http.StatusForbidden, 1, err, nil)
		return
	}
	err = db.UpdateDeviceName(deviceID, name)
	if err != nil {
		util.FillRsp(c, http.StatusForbidden, 1, err, nil)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
	})
}

func remove(c *gin.Context) {
	token := c.Query("token")
	id := c.Query("id")

	_, err := db.GetUserIDStr(token)
	if err != nil {
		util.FillRsp(c, http.StatusForbidden, 1, err, nil)
		return
	}

	err = db.RemoveDevice(id)
	if err != nil {
		util.FillRsp(c, http.StatusForbidden, 1, err, nil)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"content": gin.H{
			"message": "done",
		},
	})
}
