package webpush

import "github.com/target/goalert/gadb"

const (
	DestTypeWebPush = "builtin-webpush"

	FieldEndpoint = "endpoint"
	FieldAuthKey  = "auth"
	FieldP256DH   = "p256dh"
)

// NewDest constructs a web push destination using endpoint and keys.
func NewDest(endpoint, auth, p256dh string) gadb.DestV1 {
	return gadb.NewDestV1(DestTypeWebPush, FieldEndpoint, endpoint, FieldAuthKey, auth, FieldP256DH, p256dh)
}
