package webapi

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/sideshow/apns2"
	"github.com/sideshow/apns2/certificate"
	"github.com/sideshow/apns2/payload"
	"notification/internal/entity"
	"notification/internal/usecase"
	"notification/pkg/logger"
	"time"
)

var _ usecase.PushDeerWebAPI = (*PushDeerAPNsAPI)(nil)

type PushDeerAPNsAPI struct {
	Client *apns2.Client
	l      logger.Interface
	wsConn map[string]*websocket.Conn
}

func NewPushDeerAPNs(l logger.Interface) (*PushDeerAPNsAPI, error) {
	cert, err := certificate.FromP12File("assets/c.p12", pdKeyPassword)
	if err != nil {
		return nil, err
	}
	return &PushDeerAPNsAPI{
		Client: apns2.NewClient(cert).Production(),
		l:      l,
		wsConn: make(map[string]*websocket.Conn),
	}, nil
}

func (api *PushDeerAPNsAPI) Push(ctx context.Context, devices []*entity.Device, msg *entity.Message) error {
	pl := payload.NewPayload().
		AlertTitle(msg.PushKeyName).
		AlertBody(msg.Text).
		Sound("default").
		Category("category")

	notification := &apns2.Notification{
		Topic:      pdTopic,
		Expiration: time.Now().Add(time.Hour * 24),
		Payload:    pl.MutableContent(),
	}

	var err error
	for _, v := range devices {
		if v.Type == "ios" {
			notification.DeviceToken = v.DeviceID
			err = api.apns(ctx, notification)
		} else {
			err = api.ws(ctx, v.DeviceID)
		}
		if err != nil {
			return fmt.Errorf("push message to %d [%s] failed: %s", v.ID, v.Type, err.Error())
		}
	}
	return nil
}

func (api *PushDeerAPNsAPI) ws(ctx context.Context, deviceId string) error {
	c, ok := api.wsConn[deviceId]
	if !ok {
		return nil
	}
	return c.WriteJSON(map[string]interface{}{
		"message": "ok",
		"state":   "todo",
	})
}

func (api *PushDeerAPNsAPI) apns(ctx context.Context, n *apns2.Notification) error {
	res, err := api.Client.Push(n)
	if err != nil {
		return err
	}
	api.l.Infoln("push result", res)
	return nil
}

func (api *PushDeerAPNsAPI) Register(ctx context.Context, key string, c *websocket.Conn) {
	api.wsConn[key] = c
}
