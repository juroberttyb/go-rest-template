/*
Package store defines all interfaces and implementations for data model operations,
should be as general as possible.
*/
package store

import (
	"context"

	"github.com/A-pen-app/kickstart/models"
)

// Chat defines all chatroom related operations
type Order interface {
	GetLiveOrders(ctx context.Context, action models.OrderAction) ([]*models.Order, error)
	Make(ctx context.Context, action models.OrderAction, price, quantity int) error
	Take(ctx context.Context, action models.OrderAction, quantity int) (int, error)
	Delete(ctx context.Context, orderID string) error
}

type Crypto interface {
	CreateKey(ctx context.Context, keyRing, keyID string) error
	GetPublicKey(ctx context.Context, keyID string) (string, error)
	Decrypt(ctx context.Context, keyID string, base64Ciphertext string) ([]byte, error)
	Sign(ctx context.Context, keyID string, msg string) (string, error)
	Verify(ctx context.Context, keyID, msg, signature string) error
}
