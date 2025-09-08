package webpush

import (
	"context"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/base64"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SherClockHolmes/webpush-go"
	"github.com/target/goalert/config"
	"github.com/target/goalert/notification/nfymsg"
)

func TestSender(t *testing.T) {
	priv, vapidPub, err := webpush.GenerateVAPIDKeys()
	if err != nil {
		t.Fatal(err)
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	}))
	defer srv.Close()

	// generate fake subscription keys
	var privSub []byte
	var x, y *big.Int
	privSub, x, y, err = elliptic.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	_ = privSub
	subPub := elliptic.Marshal(elliptic.P256(), x, y)
	p256dh := base64.RawURLEncoding.EncodeToString(subPub)
	authBytes := make([]byte, 16)
	if _, err := rand.Read(authBytes); err != nil {
		t.Fatal(err)
	}
	auth := base64.RawURLEncoding.EncodeToString(authBytes)

	dest := NewDest(srv.URL, auth, p256dh)
	msg := nfymsg.Test{Base: nfymsg.Base{ID: "1", Dest: dest}}

	cfg := config.Config{}
	cfg.WebPush.Enable = true
	cfg.WebPush.VAPIDPublicKey = vapidPub
	cfg.WebPush.VAPIDPrivateKey = priv
	cfg.WebPush.TTL = 30

	ctx := cfg.Context(context.Background())

	s := NewSender(ctx, srv.Client())

	if _, err := s.SendMessage(ctx, msg); err != nil {
		t.Fatal(err)
	}
}
