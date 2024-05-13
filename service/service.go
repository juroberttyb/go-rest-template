package service

import (
	"context"
	"time"

	"github.com/A-pen-app/kickstart/models"
)

type issueOption struct {
	ttl time.Duration
}
type IssueOption func(*issueOption) error

func WithTTL(duration time.Duration) IssueOption {
	return func(opt *issueOption) error {
		opt.ttl = duration
		return nil
	}
}

type Auth interface {
	// IssueToken returns a JWT for given userID
	IssueToken(ctx context.Context, userID string, userType models.UserType, options ...IssueOption) (string, error)
	// ValidateToken verifies if given token is issued by us and returns userID if valid
	ValidateToken(ctx context.Context, token string) (*models.Claims, error)
}

type Order interface {
	GetBoard(ctx context.Context, boardType models.OrderBoardType) (*models.Board, string, error)
	Make(ctx context.Context, action models.OrderAction, price, quantity int) error
	Take(ctx context.Context, action models.OrderAction, quantity int) error
	Delete(ctx context.Context, orderID string) error
}
