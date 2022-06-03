package webapi

import (
	"context"
	"fmt"
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
}

func NewPushDeerAPNs() (*PushDeerAPNsAPI, error) {
	cert, err := certificate.FromP12File("assets/c.p12", pdKeyPassword)
	if err != nil {
		return nil, err
	}
	return &PushDeerAPNsAPI{
		Client: apns2.NewClient(cert).Production(),
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
	for _, v := range devices {
		notification.DeviceToken = v.DeviceID
		res, err := api.Client.Push(notification)
		if err != nil {
			return fmt.Errorf("push message to device %d failed: %s", v.ID, err)
		}

		api.logAPNsResult(ctx, res)
	}
	return nil
}

func (api *PushDeerAPNsAPI) logAPNsResult(ctx context.Context, res *apns2.Response) {
	if l := ctx.Value("logger"); l != nil {
		if log, ok := l.(logger.Interface); ok {
			log.Infoln("push result:", res)
		}
	}
}
