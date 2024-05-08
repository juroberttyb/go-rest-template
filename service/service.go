package service

import (
	"context"

	"github.com/A-pen-app/kickstart/models"
)

type Order interface {
	GetBoard(ctx context.Context, boardType models.OrderBoardType) (*models.Board, string, error)
	Make(ctx context.Context, action models.OrderAction, price, amount int) error
	Take(ctx context.Context, action models.OrderAction, amount int) error
	Delete(ctx context.Context, orderID string) error
}
