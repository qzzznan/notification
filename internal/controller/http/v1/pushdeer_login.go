package v1

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"strings"
)

func (r *pushDeerRoutes) fake(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": 0, "content": gin.H{
			"token": "HelloWorld",
		},
	})
}

func (r *pushDeerRoutes) apple(c *gin.Context) {
	idToken := c.Query("idToken")
	token, _ := jwt.Parse(idToken, nil)
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		err := errors.New("claims is not jwt.MapClaims")
		r.l.Errorln(err)
		pdResp(c, 400, 1, err, nil)
		return
	}
	fields, err := GetFields(claims, "sub", "email")
	if err != nil {
		r.l.Errorln(err)
		pdResp(c, 400, 1, err, nil)
		return
	}
	appleID := fields["sub"]
	email := fields["email"]
	name, _, found := strings.Cut(email, "@")
	if !found {
		name = "unknown"
	}
	uid, err := r.p.Register(c.Request.Context(), appleID, email, name)
	if err != nil {
		r.l.Errorln(err)
		pdResp(c, 400, 1, err, nil)
		return
	}
	pdResp(c, 200, 0, nil, gin.H{
		"token": uid,
	})
}

func GetFields(claims jwt.MapClaims, fields ...string) (map[string]string, error) {
	m := make(map[string]string)
	for _, v := range fields {
		f, ok := claims[v]
		if !ok {
			return nil, fmt.Errorf("field %s not found", v)
		}
		s, ok := f.(string)
		if !ok {
			return nil, fmt.Errorf("field %s is not string", v)
		}
		m[v] = s
	}
	return m, nil
}
