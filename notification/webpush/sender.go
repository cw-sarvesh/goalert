package webpush

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/mail"
	"net/url"
	"strings"

	webpushlib "github.com/SherClockHolmes/webpush-go"

	"github.com/target/goalert/config"
	"github.com/target/goalert/gadb"
	"github.com/target/goalert/notification"
	"github.com/target/goalert/notification/nfydest"
	"github.com/target/goalert/permission"
	"github.com/target/goalert/util/log"
	"github.com/target/goalert/validation"
	"github.com/target/goalert/validation/validate"
)

const (
	pushIconURL = "builtin://push"
)

type subscription struct {
	Endpoint string `json:"endpoint"`
	Keys     struct {
		Auth   string `json:"auth"`
		P256dh string `json:"p256dh"`
	} `json:"keys"`
}

type pushPayload struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	URL   string `json:"url,omitempty"`
	Type  string `json:"type,omitempty"`
}

var (
	_ nfydest.Provider      = (*Sender)(nil)
	_ nfydest.MessageSender = (*Sender)(nil)
	_ nfydest.DestValidator = (*Sender)(nil)
)

// ID implements nfydest.Provider.
func (s *Sender) ID() string { return DestTypeWebPush }

// TypeInfo implements nfydest.Provider.
func (s *Sender) TypeInfo(ctx context.Context) (*nfydest.TypeInfo, error) {
	cfg := config.FromContext(ctx)
	return &nfydest.TypeInfo{
		Type:                       DestTypeWebPush,
		Name:                       "Browser Push",
		IconURL:                    pushIconURL,
		IconAltText:                "Browser Push",
		Enabled:                    cfg.WebPush.Enable,
		SupportsAlertNotifications: true,
		SupportsStatusUpdates:      true,
		SupportsUserVerification:   true,
		UserVerificationRequired:   true,
		RequiredFields:             nil,
		UserDisclaimer:             "Notifications will be delivered to browsers where you have enabled Web Push on your profile.",
	}, nil
}

// ValidateField implements nfydest.Provider.
func (s *Sender) ValidateField(ctx context.Context, fieldID, value string) error {
	return validation.NewGenericError("unknown field ID")
}

// DisplayInfo implements nfydest.Provider.
func (s *Sender) DisplayInfo(ctx context.Context, args map[string]string) (*nfydest.DisplayInfo, error) {
	return &nfydest.DisplayInfo{
		Text:        "Browser Push",
		IconURL:     pushIconURL,
		IconAltText: "Browser Push",
	}, nil
}

// ValidateDest ensures the destination arguments are scoped correctly.
func (s *Sender) ValidateDest(ctx context.Context, dest gadb.DestV1) error {
	for field, value := range dest.Args {
		if field != FieldUserID {
			return validation.NewFieldError(field, "unexpected field")
		}
		if value == "" {
			return validation.NewFieldError(field, "user id is required")
		}
		if _, err := validate.ParseUUID(field, value); err != nil {
			return err
		}
	}
	return nil
}

// SendMessage implements nfydest.MessageSender.
func (s *Sender) SendMessage(ctx context.Context, msg notification.Message) (*notification.SentMessage, error) {
	cfg := config.FromContext(ctx)
	if cfg.WebPush.VAPIDPublicKey == "" || cfg.WebPush.VAPIDPrivateKey == "" {
		return nil, fmt.Errorf("web push VAPID keys not configured")
	}

	userID := permission.UserID(ctx)
	if userID == "" {
		return nil, fmt.Errorf("web push requires user context")
	}

	payload, err := buildPayload(cfg, msg)
	if err != nil {
		return nil, err
	}

	payloadData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal web push payload: %w", err)
	}

	subscriber := normalizeSubscriberAddress(ctx, cfg)

	subs, err := s.loadSubscriptions(ctx, userID)
	if err != nil {
		return nil, err
	}
	if len(subs) == 0 {
		return &notification.SentMessage{
			State:        notification.StateFailedPerm,
			StateDetails: "no registered browsers for web push",
		}, nil
	}

	options := &webpushlib.Options{
		Subscriber:      subscriber,
		VAPIDPublicKey:  cfg.WebPush.VAPIDPublicKey,
		VAPIDPrivateKey: cfg.WebPush.VAPIDPrivateKey,
		TTL:             60,
		Urgency:         webpushlib.UrgencyHigh,
	}

	var delivered int
	for _, sub := range subs {
		resp, sendErr := webpushlib.SendNotification(payloadData, &webpushlib.Subscription{
			Endpoint: sub.Endpoint,
			Keys: webpushlib.Keys{
				Auth:   sub.Keys.Auth,
				P256dh: sub.Keys.P256dh,
			},
		}, options)

		suffix := endpointSuffix(sub.Endpoint)
		if resp != nil {
			if resp.Body != nil {
				_ = resp.Body.Close()
			}

			if resp.StatusCode == http.StatusGone || resp.StatusCode == http.StatusNotFound {
				if delErr := s.deleteSubscription(ctx, sub.Endpoint); delErr != nil {
					log.Logf(ctx, "webpush: failed to delete expired subscription (%s): %v", suffix, delErr)
				} else {
					log.Logf(ctx, "webpush: removed expired subscription (%s)", suffix)
				}
			}
		}

		if sendErr != nil {
			log.Logf(ctx, "webpush: send failed for endpoint (%s): %v", suffix, sendErr)
			continue
		}

		delivered++
	}

	if delivered == 0 {
		return &notification.SentMessage{
			State:        notification.StateFailedPerm,
			StateDetails: "web push delivery failed for all subscriptions",
		}, nil
	}

	return &notification.SentMessage{
		State: notification.StateSent,
	}, nil
}

func (s *Sender) loadSubscriptions(ctx context.Context, userID string) ([]subscription, error) {
	rows, err := s.db.QueryContext(ctx, `
		select data
		from user_web_push_subscriptions
		where user_id = $1::uuid
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("query web push subscriptions: %w", err)
	}
	defer rows.Close()

	var out []subscription
	for rows.Next() {
		var raw json.RawMessage
		if err := rows.Scan(&raw); err != nil {
			log.Logf(ctx, "webpush: scan subscription failed: %v", err)
			continue
		}
		var sub subscription
		if err := json.Unmarshal(raw, &sub); err != nil {
			log.Logf(ctx, "webpush: invalid subscription payload: %v", err)
			continue
		}
		if sub.Endpoint == "" || sub.Keys.Auth == "" || sub.Keys.P256dh == "" {
			log.Logf(ctx, "webpush: subscription missing required fields; skipping")
			continue
		}
		out = append(out, sub)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate web push subscriptions: %w", err)
	}

	return out, nil
}

func (s *Sender) deleteSubscription(ctx context.Context, endpoint string) error {
	_, err := s.db.ExecContext(ctx, `
		delete from user_web_push_subscriptions
		where endpoint = $1
	`, endpoint)
	if err != nil {
		return fmt.Errorf("delete web push subscription: %w", err)
	}
	return nil
}

func buildPayload(cfg config.Config, msg notification.Message) (pushPayload, error) {
	appName := cfg.ApplicationName()
	switch m := msg.(type) {
	case notification.Alert:
		body := strings.TrimSpace(m.Summary)
		if body == "" {
			body = fmt.Sprintf("Alert #%d is active.", m.AlertID)
		}
		return pushPayload{
			Type:  "alert",
			Title: fmt.Sprintf("Alert #%d · %s", m.AlertID, m.ServiceName),
			Body:  body,
			URL:   fmt.Sprintf("/alerts/%d", m.AlertID),
		}, nil
	case notification.AlertBundle:
		return pushPayload{
			Type:  "alert-bundle",
			Title: fmt.Sprintf("%s Alerts", m.ServiceName),
			Body:  fmt.Sprintf("%d unacknowledged alerts", m.Count),
			URL:   fmt.Sprintf("/services/%s/alerts", m.ServiceID),
		}, nil
	case notification.AlertStatus:
		body := strings.TrimSpace(m.LogEntry)
		if body == "" {
			body = fmt.Sprintf("Alert #%d status updated.", m.AlertID)
		}
		return pushPayload{
			Type:  "alert-status",
			Title: fmt.Sprintf("Alert #%d update", m.AlertID),
			Body:  body,
			URL:   fmt.Sprintf("/alerts/%d", m.AlertID),
		}, nil
	case notification.Test:
		return pushPayload{
			Type:  "test",
			Title: fmt.Sprintf("%s Test Message", appName),
			Body:  "This is a test notification.",
			URL:   "/profile",
		}, nil
	case notification.Verification:
		return pushPayload{
			Type:  "verification",
			Title: fmt.Sprintf("%s Verification Code", appName),
			Body:  fmt.Sprintf("Enter code %s to verify this device.", m.Code),
			URL:   "/profile",
		}, nil
	default:
		return pushPayload{}, fmt.Errorf("web push: unsupported message type %T", msg)
	}
}

func normalizeSubscriberAddress(ctx context.Context, cfg config.Config) string {
	subscriber := strings.TrimSpace(cfg.WebPush.SubscriberEmail)
	if subscriber != "" {
		addr, err := mail.ParseAddress(subscriber)
		if err != nil || addr.Address == "" {
			log.Logf(ctx, "webpush: invalid subscriber email %q: %v", subscriber, err)
			subscriber = ""
		} else {
			subscriber = strings.ToLower(addr.Address)
		}
	}

	if subscriber != "" {
		return subscriber
	}

	subscriber = "no-reply@localhost"
	if pubURL := cfg.PublicURL(); pubURL != "" {
		if u, err := url.Parse(pubURL); err == nil {
			host := u.Hostname()
			switch host {
			case "", "localhost", "127.0.0.1", "::1":
				// keep default placeholder for local development
			default:
				subscriber = "no-reply@" + host
			}
		}
	}

	return subscriber
}

func endpointSuffix(endpoint string) string {
	endpoint = strings.TrimSpace(endpoint)
	if endpoint == "" {
		return ""
	}
	const suffixLen = 16
	if len(endpoint) <= suffixLen {
		return endpoint
	}
	return endpoint[len(endpoint)-suffixLen:]
}
