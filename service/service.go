package service

import (
	"context"

	"github.com/A-pen-app/kickstart/models"
)

type Order interface {
	GetBoard(ctx context.Context, boardType models.OrderBoardType) (*models.Board, string, error)
	Make(ctx context.Context, action models.OrderAction, price, quantity int) error
	Take(ctx context.Context, action models.OrderAction, quantity int) error
	Delete(ctx context.Context, orderID string) error
}
