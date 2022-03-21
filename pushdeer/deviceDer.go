package pushdeer

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/qzzznan/notification/db"
	"github.com/qzzznan/notification/model"
	"github.com/qzzznan/notification/util"
	"net/http"
)

func reg(c *gin.Context) {
	regInfo := &model.RegInfo{}
	err := c.Bind(&regInfo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  1,
			"error": err.Error(),
		})
		return
	}

	err = util.Validate.Struct(regInfo)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"code": 1, "error": err.Error(),
		})
		return
	}
	id, err := db.GetUserID(regInfo.Token)
	if err != nil {
		util.FillRsp(c, http.StatusForbidden, 1, err, nil)
		return
	}

	err = db.InsertDevice(&model.Device{
		UserID:   id,
		DeviceID: regInfo.DeviceID,
		Type:     "none",
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
	token, ok := c.GetPostForm("token")
	if !ok {
		util.FillRsp(c, http.StatusForbidden, 1, fmt.Errorf("token is required"), nil)
		return
	}
	id, err := db.GetUserID(token)
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

func rename(c *gin.Context) {
	token := c.GetString("token")
	id := c.GetInt64("id")
	name := c.GetString("name")

	_, err := db.GetUserID(token)
	if err != nil {
		util.FillRsp(c, http.StatusForbidden, 1, err, nil)
		return
	}

	err = db.UpdateDeviceName(id, name)
	if err != nil {
		util.FillRsp(c, http.StatusForbidden, 1, err, nil)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
	})
}

func remove(c *gin.Context) {
	token := c.GetString("token")
	id := c.GetInt64("id")

	_, err := db.GetUserID(token)
	if err != nil {
		util.FillRsp(c, http.StatusForbidden, 1, err, nil)
		return
	}

	_ = id
	//TODO: remove

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"content": gin.H{
			"message": "done",
		},
	})
}
