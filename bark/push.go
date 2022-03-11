package bark

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/samber/lo"
	"github.com/sideshow/apns2"
	"github.com/sideshow/apns2/payload"
	"github.com/sideshow/apns2/token"
	"golang.org/x/net/http2"
	"net/http"
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

func pushMessage() error {
	pl := payload.NewPayload().
		AlertTitle("").
		AlertBody("").
		Sound("").
		Category("")

	rsp, err := client.Push(&apns2.Notification{
		DeviceToken: "",
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
