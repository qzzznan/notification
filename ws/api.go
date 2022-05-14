package ws

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/qzzznan/notification/log"
)

func GetWsConn(deviceToken string) *websocket.Conn {
	return m[deviceToken]
}

func echoHandler(c *gin.Context) {
	con, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}

	data := make(map[string]interface{})
	for {
		err = con.ReadJSON(&data)
		if err != nil {
			log.Errorln("echo ws read:", err)
			break
		}
		err = con.WriteJSON(gin.H{"message": "ok"})
		if err != nil {
			log.Errorln("echo ws write:", err)
			break
		}
	}
}
