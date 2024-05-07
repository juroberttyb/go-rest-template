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
	New(ctx context.Context, userID, kickstartID string, email *string) error
	Get(ctx context.Context, kickstartID string) (*models.Order, error)
	GetOrders(ctx context.Context, limit int) ([]*models.Order, error)
	GetOrderIDs(ctx context.Context, uid string, limit int, status models.OrderStatus) ([]string, error)
	GetMyOrders(ctx context.Context, uid string, limit int, isAttending bool) ([]*models.Order, error)
	GetIsAttended(ctx context.Context, uid, kickstartID string) (bool, error)
}
