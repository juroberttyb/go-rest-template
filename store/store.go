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
	Make(ctx context.Context, action models.OrderAction, price, amount int) error
	Take(ctx context.Context, action models.OrderAction, amount int) (int, error)
	Delete(ctx context.Context, orderID string) error
}
