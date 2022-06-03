package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"notification/internal/entity"
	"strconv"
	"time"
)

func (r *pushDeerRoutes) push(c *gin.Context) {
	pushKey := c.Query("pushkey")
	text := c.Query("text")
	desp := c.Query("desp")
	typ := c.DefaultQuery("type", "markdown")

	err := r.p.PushMessage(c.Request.Context(), pushKey, &entity.Message{
		Text:   text,
		Note:   desp,
		Type:   typ,
		SendAt: time.Now(),
	})

	if err != nil {
		pdResp(c, r.l, 400, 1, err, nil)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0, "content": gin.H{
			"result": []string{
				`{"count": 1, "logs": [], "success": "ok"`,
			},
		},
	})
}

func (r *pushDeerRoutes) msgList(c *gin.Context) {
	token := c.Query("token")
	limit := c.Query("limit")

	lim, err := strconv.ParseInt(limit, 10, 32)
	if err != nil {
		lim = 100
	}
	ctx := c.Request.Context()
	msg, err := r.p.GetMessage(ctx, token, 0, uint64(lim))
	if err != nil {
		pdResp(c, r.l, 400, 1, err, nil)
		return
	}
	pdResp(c, r.l, 200, 0, nil, gin.H{
		"messages": msg,
	})
}

func (r *pushDeerRoutes) msgRm(c *gin.Context) {
	token := c.Query("token")
	msgID := c.Query("id")

	ctx := c.Request.Context()
	err := r.p.RemoveMessage(ctx, token, msgID)
	if err != nil {
		pdResp(c, r.l, 400, 1, err, nil)
		return
	}

	pdResp(c, r.l, 200, 0, nil, gin.H{
		"message": "done",
	})
}
