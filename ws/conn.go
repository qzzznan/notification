package ws

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader *websocket.Upgrader
var m map[string]*websocket.Conn

func init() {
	upgrader = &websocket.Upgrader{}
	m = make(map[string]*websocket.Conn)
}

func InitWsHandler(e *gin.Engine) {
	g := e.Group("/ws")
	g.Any("/echo", echoHandler)
}
