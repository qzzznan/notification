package v1

import "github.com/gin-gonic/gin"

func (r *pushDeerRoutes) merge(c *gin.Context) {

}

func (r *pushDeerRoutes) info(c *gin.Context) {
	token := c.Query("token")
	user, err := r.p.GetUser(c.Request.Context(), token)
	if err != nil {

		pdResp(c, r.l, 400, 1, err, nil)
		return
	}
	pdResp(c, r.l, 200, 0, nil, user)
}
