package bark

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/qzzznan/notification/model"
	"github.com/samber/lo"
	"github.com/sideshow/apns2"
	"github.com/sideshow/apns2/payload"
	"github.com/sideshow/apns2/token"
	"golang.org/x/net/http2"
	"net/http"
	"strings"
	"time"
)

var client *apns2.Client

func InitPushClient() error {
	key, err := token.AuthKeyFromBytes([]byte(apnsPrivateKey))
	if err != nil {
		return err
	}

	CAs, err := x509.SystemCertPool()
	if err != nil {
		return err
	}

	lo.ForEach[string](apnsCAs, func(t string, _ int) {
		CAs.AppendCertsFromPEM([]byte(t))
	})

	client = &apns2.Client{
		Token: &token.Token{
			AuthKey: key,
			KeyID:   keyID,
			TeamID:  teamID,
		},
		HTTPClient: &http.Client{
			Transport: &http2.Transport{
				DialTLS: apns2.DialTLS,
				TLSClientConfig: &tls.Config{
					RootCAs: CAs,
				},
			},
			Timeout: apns2.HTTPClientTimeout,
		},
		Host: apns2.HostProduction,
	}
	return nil
}

func PushMessage(msg *model.APNsMessage) error {
	pl := payload.NewPayload().
		AlertTitle(msg.Title).
		AlertBody(msg.Body).
		Sound(msg.Sound).
		Category(msg.Category)

	group, ok := msg.Data["group"]
	if ok {
		g, ok := group.(string)
		if ok {
			pl = pl.ThreadID(g)
		}
	}

	for k, v := range msg.Data {
		pl.Custom(strings.ToLower(k), fmt.Sprintf("%v", v))
	}

	rsp, err := client.Push(&apns2.Notification{
		DeviceToken: msg.DeviceToken,
		Topic:       topic,
		Payload:     pl.MutableContent(),
		Expiration:  time.Now().Add(time.Hour * 24),
	})

	if err != nil {
		return err
	}
	if rsp.StatusCode != http.StatusOK {
		return fmt.Errorf("bark push %v error", rsp.Reason)
	}
	return nil
}
