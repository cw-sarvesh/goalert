package webpush

import (
	"database/sql"

	"github.com/target/goalert/gadb"
)

// DestTypeWebPush identifies the notification destination type for browser push notifications.
const DestTypeWebPush = "builtin-web-push"

// FieldUserID is the destination argument key that ties a Web Push contact method to a specific user.
const FieldUserID = "user_id"

// Sender implements nfydest.Provider and nfydest.MessageSender for Web Push.
type Sender struct {
	db *sql.DB
}

// NewSender constructs a new Web Push sender backed by the provided database handle.
func NewSender(db *sql.DB) *Sender {
	return &Sender{db: db}
}

// NewDest constructs a destination representing browser push notifications. If userID is provided,
// it will be associated with the destination so that it is unique per user.
func NewDest(userID string) gadb.DestV1 {
	dest := gadb.NewDestV1(DestTypeWebPush)
	if userID != "" {
		dest.SetArg(FieldUserID, userID)
	}
	return dest
}
