package v1

import "github.com/gin-gonic/gin"

func (r *pushDeerRoutes) gen(c *gin.Context) {
	token := c.Query("token")
	name := c.Query("name")

	keys, err := r.p.GenNewKey(c.Request.Context(), token, name)
	if err != nil {
		pdResp(c, r.l, 400, 1, err, nil)
		return
	}
	pdResp(c, r.l, 200, 0, nil, gin.H{
		"keys": keys,
	})
}

func (r *pushDeerRoutes) keyRename(c *gin.Context) {
	token := c.Query("token")
	kid := c.Query("id")
	newName := c.Query("name")

	err := r.p.RenameKey(c.Request.Context(), token, kid, newName)
	if err != nil {
		pdResp(c, r.l, 400, 1, err, nil)
		return
	}
	pdResp(c, r.l, 200, 0, nil, gin.H{
		"message": "done",
	})
}

func (r *pushDeerRoutes) keyRegen(c *gin.Context) {
	token := c.Query("token")
	kid := c.Query("id")

	err := r.p.RegenKey(c.Request.Context(), token, kid)
	if err != nil {
		pdResp(c, r.l, 400, 1, err, nil)
		return
	}
	pdResp(c, r.l, 200, 0, nil, gin.H{
		"message": "done",
	})
}

func (r *pushDeerRoutes) keyList(c *gin.Context) {
	token := c.Query("token")

	keys, err := r.p.GetAllKey(c.Request.Context(), token)
	if err != nil {
		pdResp(c, r.l, 400, 1, err, nil)
		return
	}
	pdResp(c, r.l, 200, 0, nil, gin.H{
		"keys": keys,
	})
}

func (r *pushDeerRoutes) keyRemove(c *gin.Context) {
	token := c.Query("token")
	id := c.Query("id")
	err := r.p.RemoveKey(c.Request.Context(), token, id)
	if err != nil {
		pdResp(c, r.l, 400, 1, err, nil)
		return
	}
	pdResp(c, r.l, 200, 0, nil, gin.H{
		"message": "done",
	})
}
