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

func GetWsConn(deviceToken string) *websocket.Conn {
	return m[deviceToken]
}

func WsHandler(c *gin.Context) {
	deviceToken := c.Query("device_token")
	if deviceToken == "" {
		c.JSON(400, gin.H{
			"message": "device_token is required",
		})
		return
	}
	con, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}
	m[deviceToken] = con

	data := &struct{}{}
	for {
		err = con.ReadJSON(data)
		if err != nil {

		}
		err = con.WriteJSON(gin.H{"message": "hello"})
		if err != nil {

		}
	}
}
