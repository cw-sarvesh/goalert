package message

import (
	"github.com/target/goalert/gadb"
	"github.com/target/goalert/notification"
)

// splitPendingByType will split a list of messages returning only unsent and matching at least one of the provided
// types. Any sent or non-type-matching messages will be returned in the remainder.
func splitPendingByType(messages []Message, types ...gadb.EnumOutgoingMessagesType) (matching, remainder []Message) {
mainLoop:
	for _, msg := range messages {
		if !msg.SentAt.IsZero() {
			remainder = append(remainder, msg)
			continue
		}

		for _, t := range types {
			if msg.Type != t {
				continue
			}

			if msg.Type == notification.MessageTypeAlert {
				if msg.AlertStatus == notification.AlertStateAcknowledged || msg.AlertStatus == notification.AlertStateClosed {
					remainder = append(remainder, msg)
					continue mainLoop
				}
			}

			matching = append(matching, msg)
			continue mainLoop
		}

		// doesn't match any specified types, keep it
		remainder = append(remainder, msg)
	}

	return matching, remainder
}
