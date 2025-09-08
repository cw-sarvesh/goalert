package webpush

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/SherClockHolmes/webpush-go"
	"github.com/target/goalert/config"
	"github.com/target/goalert/notification"
	"github.com/target/goalert/notification/nfydest"
)

// Sender implements nfydest.MessageSender for Web Push notifications.
type Sender struct {
	client *http.Client
}

// NewSender initializes a new web push sender using the provided HTTP client.
// If client is nil http.DefaultClient will be used.
func NewSender(_ context.Context, client *http.Client) *Sender {
	if client == nil {
		client = http.DefaultClient
	}
	return &Sender{client: client}
}

var _ nfydest.MessageSender = &Sender{}

// SendMessage sends a push message using Web Push protocol.
func (s *Sender) SendMessage(ctx context.Context, msg notification.Message) (*notification.SentMessage, error) {
	cfg := config.FromContext(ctx)

	// ensure webpush is enabled and keys are set
	if !cfg.WebPush.Enable {
		return nil, nfydest.ErrNotEnabled
	}

	sub := &webpush.Subscription{
		Endpoint: msg.DestArg(FieldEndpoint),
		Keys: webpush.Keys{
			Auth:   msg.DestArg(FieldAuthKey),
			P256dh: msg.DestArg(FieldP256DH),
		},
	}

	payload := map[string]string{"msg": msg.MsgID()}

	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	resp, err := webpush.SendNotification(data, sub, &webpush.Options{
		Subscriber:      cfg.General.PublicURL,
		VAPIDPublicKey:  cfg.WebPush.VAPIDPublicKey,
		VAPIDPrivateKey: cfg.WebPush.VAPIDPrivateKey,
		TTL:             cfg.WebPush.TTL,
		HTTPClient:      s.client,
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return &notification.SentMessage{State: notification.StateFailedTemp}, nil
	}

	return &notification.SentMessage{State: notification.StateSent}, nil
}
