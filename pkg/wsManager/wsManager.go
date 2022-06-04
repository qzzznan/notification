package wsManager

import (
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

type WsManager struct {
	u *websocket.Upgrader
}

func New() *WsManager {
	return &WsManager{
		u: &websocket.Upgrader{
			HandshakeTimeout: time.Second * 5,
			ReadBufferSize:   4096,
			WriteBufferSize:  4096,
		},
	}
}

func (wm *WsManager) Upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (*websocket.Conn, error) {
	c, err := wm.u.Upgrade(w, r, responseHeader)
	if err != nil {
		return nil, err
	}
	return c, nil
}
