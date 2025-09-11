package engine

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/target/goalert/engine/message"
	"github.com/target/goalert/gadb"
	"github.com/target/goalert/notification"
	"github.com/target/goalert/notification/twilio"
	"github.com/target/goalert/user/contactmethod"
)

type cmStoreStub struct {
	cms []contactmethod.ContactMethod
}

func (s cmStoreStub) FindAll(ctx context.Context, db gadb.DBTX, userID string) ([]contactmethod.ContactMethod, error) {
	return s.cms, nil
}

func TestApplyHighPriorityOverride(t *testing.T) {
	voiceID := uuid.New()
	smsID := uuid.New()
	voiceDest := twilio.NewVoiceDest("+15555550123")
	smsDest := twilio.NewSMSDest("+15555550123")
	store := cmStoreStub{cms: []contactmethod.ContactMethod{{ID: voiceID, Dest: voiceDest}, {ID: smsID, Dest: smsDest}}}
	msg := &message.Message{UserID: "u1", DestID: notification.DestID{CMID: uuid.NullUUID{UUID: smsID, Valid: true}}, Dest: smsDest}
	meta := map[string]string{"alerts/priority": "high"}
	applyHighPriorityOverride(context.Background(), nil, store, msg, meta, "alerts/priority", "high")
	if msg.Dest.Type != twilio.DestTypeTwilioVoice {
		t.Fatalf("expected voice dest, got %s", msg.Dest.Type)
	}
	if msg.DestID.CMID.UUID != voiceID {
		t.Fatalf("expected destID to be voice contact method")
	}
}

func TestApplyHighPriorityOverrideNoMatch(t *testing.T) {
	voiceID := uuid.New()
	smsID := uuid.New()
	voiceDest := twilio.NewVoiceDest("+15555550123")
	smsDest := twilio.NewSMSDest("+15555550123")
	store := cmStoreStub{cms: []contactmethod.ContactMethod{{ID: voiceID, Dest: voiceDest}, {ID: smsID, Dest: smsDest}}}
	msg := &message.Message{UserID: "u1", DestID: notification.DestID{CMID: uuid.NullUUID{UUID: smsID, Valid: true}}, Dest: smsDest}
	meta := map[string]string{"other": "val"}
	applyHighPriorityOverride(context.Background(), nil, store, msg, meta, "alerts/priority", "high")
	if msg.Dest.Type != twilio.DestTypeTwilioSMS {
		t.Fatalf("expected SMS dest to remain, got %s", msg.Dest.Type)
	}
	if msg.DestID.CMID.UUID != smsID {
		t.Fatalf("expected destID to remain original")
	}
}
