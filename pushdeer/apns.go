package pushdeer

import (
	"fmt"
	"github.com/qzzznan/notification/log"
	"github.com/qzzznan/notification/model"
	"github.com/sideshow/apns2"
	"github.com/sideshow/apns2/certificate"
	"github.com/sideshow/apns2/payload"
	"time"
)

const (
	topic       = "com.pushdeer.self.ios"
	keyID       = "66M7BD2GCV"
	teamID      = "HUJ6HAE4VU"
	keyPassword = "64wtMhU4mULj"
)

var client *apns2.Client

func InitPushClient(file string) error {
	cert, err := certificate.FromP12File(file, keyPassword)
	if err != nil {
		return err
	}

	client = apns2.NewClient(cert).Production()
	return nil
}

func PushMessage(msg *model.Message, devices []*model.Device) error {
	pl := payload.NewPayload().
		AlertTitle(msg.PushKeyName).
		AlertBody(msg.Text).
		Sound("default").
		Category("category")

	notification := &apns2.Notification{
		Topic:      topic,
		Expiration: time.Now().Add(time.Hour * 24),
		Payload:    pl.MutableContent(),
	}
	for _, v := range devices {
		notification.DeviceToken = v.DeviceID
		res, err := client.Push(notification)
		if err != nil {
			return fmt.Errorf("push message to device %d failed: %s", v.ID, err)
		}
		log.Infoln("push result:", res)
	}
	return nil
}
