package service

import (
	"context"

	"github.com/A-pen-app/kickstart/models"
)

type Order interface {
	New(ctx context.Context, userID, kickstartID string, email *string) error
	Get(ctx context.Context, userID, kickstartID string) (*models.Order, error)
	GetOrders(ctx context.Context, userID string, next string, count int, filter models.OrderStatus) ([]*models.Order, string, error)
}
