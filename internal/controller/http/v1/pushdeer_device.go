package v1

import (
	"github.com/gin-gonic/gin"
	"notification/internal/entity"
)

func (r *pushDeerRoutes) reg(c *gin.Context) {
	info := &entity.RegInfo{}
	err := c.BindQuery(info)
	if err != nil {
		r.l.Errorln(err)
		pdResp(c, 400, 1, err, nil)
		return
	}
	err = r.v.Struct(info)
	if err != nil {
		r.l.Errorln(err)
		pdResp(c, 400, 1, err, nil)
		return
	}

	devices, err := r.p.RegisterDevice(c.Request.Context(), info)
	if err != nil {
		r.l.Errorln(err)
		pdResp(c, 400, 1, err, nil)
		return
	}

	pdResp(c, 200, 0, nil, gin.H{
		"devices": devices,
	})
}

func (r *pushDeerRoutes) list(c *gin.Context) {
	token := c.Query("token")
	devices, err := r.p.GetAllDevice(c.Request.Context(), token)
	if err != nil {
		r.l.Errorln(err)
		pdResp(c, 400, 1, err, nil)
		return
	}
	pdResp(c, 200, 0, nil, gin.H{
		"devices": devices,
	})
}

func (r *pushDeerRoutes) rename(c *gin.Context) {
	token := c.Query("token")
	deviceID := c.Query("id")
	name := c.Query("name")

	ctx := c.Request.Context()
	err := r.p.ValidateToken(ctx, token)
	if err != nil {
		r.l.Errorln(err)
		pdResp(c, 400, 1, err, nil)
		return
	}
	err = r.p.RenameDevice(ctx, deviceID, name)
	if err != nil {
		r.l.Errorln(err)
		pdResp(c, 400, 1, err, nil)
		return
	}
	pdResp(c, 200, 0, nil, nil)
}

func (r *pushDeerRoutes) remove(c *gin.Context) {
	token := c.Query("token")
	deviceID := c.Query("id")

	ctx := c.Request.Context()
	err := r.p.ValidateToken(ctx, token)
	if err != nil {
		r.l.Errorln(err)
		pdResp(c, 400, 1, err, nil)
		return
	}
	err = r.p.RemoveDevice(ctx, deviceID)
	if err != nil {
		r.l.Errorln(err)
		pdResp(c, 400, 1, err, nil)
		return
	}
	pdResp(c, 200, 0, nil, gin.H{
		"message": "done",
	})
}

func (r *pushDeerRoutes) upgrade(c *gin.Context) {
	ctx := c.Request.Context()
	deviceID := c.Query("device_id")
	err := r.p.Upgrade(ctx, deviceID, c.Writer, c.Request, nil)
	if err != nil {
		r.l.Errorln(err)
		pdResp(c, 400, 1, err, nil)
		return
	}
}
