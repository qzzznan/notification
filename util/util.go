package util

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"github.com/satori/go.uuid"
)

var Validate = validator.New()

func GenUID() string {
	return uuid.NewV4().String()
}

func GenToken(token string) string {
	checkSum := sha256.Sum256([]byte(token))
	arr := make([]byte, 0, len(checkSum))
	for _, v := range checkSum {
		arr = append(arr, v)
	}
	return hex.EncodeToString(arr) + "-SHA256"
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

func FillRsp(c *gin.Context, state int, code int, err error, content interface{}) {
	if err != nil || code != 0 {
		c.JSON(state, gin.H{
			"code":  code,
			"error": err,
		})
		return
	}
	c.JSON(state, gin.H{
		"code":    code,
		"content": content,
	})
}
