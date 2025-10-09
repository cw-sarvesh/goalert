package twilio

import (
	"context"
	"testing"

	"github.com/target/goalert/config"
)

func TestSMSTypeInfoEnabledWithTwilio(t *testing.T) {
	var s SMS

	cfg := config.Config{}
	cfg.Twilio.Enable = true

	info, err := s.TypeInfo(cfg.Context(context.Background()))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !info.Enabled {
		t.Fatal("expected SMS contact method to be enabled when Twilio is enabled")
	}
}

func TestSMSTypeInfoDisabledByConfig(t *testing.T) {
	var s SMS

	cfg := config.Config{}
	cfg.Twilio.Enable = true
	cfg.Twilio.DisableSMSContactMethod = true

	info, err := s.TypeInfo(cfg.Context(context.Background()))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info.Enabled {
		t.Fatal("expected SMS contact method to be disabled when DisableSMSContactMethod is true")
	}
}
