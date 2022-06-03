package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"notification/internal/usecase"
	"notification/pkg/logger"
)

type UseCase struct {
	B usecase.Bark
	P usecase.PushDeer
}

func NewRouter(e *gin.Engine, l logger.Interface, u *UseCase) {
	e.Use(gin.Logger())
	e.Use(gin.Recovery())

	e.GET("/healthz", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	h := e.Group("/v1")
	{
		newBarkRouter(h, l, u.B)
		newPushDeerRouter(h, l, u.P)
	}
}
