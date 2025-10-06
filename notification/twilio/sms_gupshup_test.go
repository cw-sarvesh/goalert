package twilio

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/target/goalert/config"
	"github.com/target/goalert/gadb"
	"github.com/target/goalert/notification"
	"github.com/target/goalert/notification/gupshup"
	"github.com/target/goalert/notification/nfymsg"
)

func TestSendViaGupshup(t *testing.T) {
	t.Parallel()

	formCh := make(chan url.Values, 1)
	headerCh := make(chan http.Header, 1)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headerCh <- r.Header.Clone()

		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)

		vals, err := url.ParseQuery(string(body))
		require.NoError(t, err)

		formCh <- vals

		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(map[string]string{"messageId": "msg-123"})
		require.NoError(t, err)
	}))
	defer srv.Close()

	sms := &SMS{
		limit:   newReplyLimiter(),
		gupshup: gupshup.NewClient(gupshup.Config{BaseURL: srv.URL, APIKey: "secret", Source: "GSRC", HTTPClient: srv.Client()}),
	}

	cfg := config.Config{}
	cfg.General.ApplicationName = "GoAlert"
	cfg.General.PublicURL = "https://goalert.example"
	cfg.Gupshup.Enable = true
	ctx := cfg.Context(context.Background())

	msg := notification.Alert{AlertID: 42, Summary: "Example alert"}

	sent, err := sms.sendViaGupshup(ctx, sms.gupshup, msg, "+15555551234")
	require.NoError(t, err)
	require.Equal(t, notification.StateSent, sent.State)
	require.Equal(t, "msg-123", sent.ExternalID)

	hdr := <-headerCh
	require.Equal(t, "secret", hdr.Get("apikey"))
	require.Equal(t, "application/x-www-form-urlencoded", hdr.Get("Content-Type"))

	vals := <-formCh
	require.Equal(t, "SMS", vals.Get("channel"))
	require.Equal(t, "GSRC", vals.Get("source"))
	require.Equal(t, "+15555551234", vals.Get("destination"))

	text := vals.Get("message")
	require.Contains(t, text, "GoAlert: Alert #42: Example alert")
	require.Contains(t, text, "https://goalert.example/alerts/42")
	require.NotContains(t, text, "Reply '")
}

func TestSendMessageUsesGupshupWhenEnabled(t *testing.T) {
	t.Parallel()

	formCh := make(chan url.Values, 1)
	headerCh := make(chan http.Header, 1)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headerCh <- r.Header.Clone()

		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)

		vals, err := url.ParseQuery(string(body))
		require.NoError(t, err)

		formCh <- vals

		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(map[string]string{"messageId": "msg-abc"})
		require.NoError(t, err)
	}))
	defer srv.Close()

	sms := &SMS{
		limit:   newReplyLimiter(),
		gupshup: gupshup.NewClient(gupshup.Config{BaseURL: srv.URL, APIKey: "secret", Source: "GSRC", HTTPClient: srv.Client()}),
	}

	cfg := config.Config{}
	cfg.General.ApplicationName = "GoAlert"
	cfg.General.PublicURL = "https://goalert.example"
	cfg.Twilio.Enable = true
	cfg.Twilio.FromNumber = "+19999999999"
	cfg.Gupshup.Enable = true
	ctx := cfg.Context(context.Background())

	dest := gadb.NewDestV1(DestTypeTwilioSMS, FieldPhoneNumber, "+15555551234")
	msg := notification.Alert{
		Base:      nfymsg.Base{ID: "msg-001", Dest: dest},
		AlertID:   42,
		Summary:   "Example alert",
		ServiceID: "svc-123",
	}

	sent, err := sms.SendMessage(ctx, msg)
	require.NoError(t, err)
	require.Equal(t, notification.StateSent, sent.State)
	require.Equal(t, "msg-abc", sent.ExternalID)

	hdr := <-headerCh
	require.Equal(t, "secret", hdr.Get("apikey"))
	require.Equal(t, "application/x-www-form-urlencoded", hdr.Get("Content-Type"))

	vals := <-formCh
	require.Equal(t, "SMS", vals.Get("channel"))
	require.Equal(t, "GSRC", vals.Get("source"))
	require.Equal(t, "+15555551234", vals.Get("destination"))

	text := vals.Get("message")
	require.Contains(t, text, "GoAlert: Alert #42: Example alert")
	require.Contains(t, text, "https://goalert.example/alerts/42")
	require.NotContains(t, text, "Reply '")
}
