package webpush

import (
	"context"

	"github.com/target/goalert/gadb"
)

// RemoveUserSubscriptions deletes any stored push subscriptions for the given user.
func RemoveUserSubscriptions(ctx context.Context, db gadb.DBTX, userID string) error {
	if userID == "" {
		return nil
	}

	_, err := db.ExecContext(ctx, `
		delete from user_web_push_subscriptions
		where user_id = $1::uuid
	`, userID)
	return err
}
