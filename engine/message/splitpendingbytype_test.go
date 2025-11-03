package message

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/target/goalert/notification"
)

func TestSplitPendingByType(t *testing.T) {
	msgs := []Message{
		{SentAt: time.Unix(1, 0), Type: notification.MessageTypeAlert},
		{Type: notification.MessageTypeAlertBundle},
		{Type: notification.MessageTypeAlert},
		{Type: notification.MessageTypeTest},
	}

	match, remainder := splitPendingByType(msgs, notification.MessageTypeAlertBundle, notification.MessageTypeTest)
	assert.ElementsMatch(t, []Message{
		{Type: notification.MessageTypeAlertBundle},
		{Type: notification.MessageTypeTest},
	}, match)
	assert.ElementsMatch(t, []Message{
		{SentAt: time.Unix(1, 0), Type: notification.MessageTypeAlert},
		{Type: notification.MessageTypeAlert},
	}, remainder)

}

func TestSplitPendingByTypeFiltersAcknowledgedAlerts(t *testing.T) {
	msgs := []Message{
		{Type: notification.MessageTypeAlert, AlertStatus: notification.AlertStateAcknowledged},
		{Type: notification.MessageTypeAlert, AlertStatus: notification.AlertStateUnacknowledged},
		{Type: notification.MessageTypeAlert, AlertStatus: notification.AlertStateClosed},
	}

	match, remainder := splitPendingByType(msgs, notification.MessageTypeAlert)
	assert.ElementsMatch(t, []Message{
		{Type: notification.MessageTypeAlert, AlertStatus: notification.AlertStateUnacknowledged},
	}, match)
	assert.ElementsMatch(t, []Message{
		{Type: notification.MessageTypeAlert, AlertStatus: notification.AlertStateAcknowledged},
		{Type: notification.MessageTypeAlert, AlertStatus: notification.AlertStateClosed},
	}, remainder)
}
