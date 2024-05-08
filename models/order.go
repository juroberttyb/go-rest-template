package models

import (
	"time"
)

type OrderAction string

const (
	Buy  OrderAction = "buy"
	Sell OrderAction = "sell"
)

type OrderBoardType string

const (
	Live    OrderBoardType = "live"
	History OrderBoardType = "history"
	Removed OrderBoardType = "removed"
)

type Order struct {
	ID     string      `json:"id" db:"id" example:"uuid"`
	Action OrderAction `json:"action" db:"action" example:"buy"`
	// using int instead of float64 to avoid floating point precision issue
	Price     int       `json:"price" db:"price" example:"10"`
	Amount    int       `json:"amount" db:"amount" example:"100"`
	CreatedAt time.Time `json:"created_at" db:"created_at" example:"2021-01-01T00:00:00Z"`
}

type Board struct {
	LatestPrice int      `json:"latest_price"`
	BuyOrders   []*Order `json:"buy_orders"`
	SellOrders  []*Order `json:"sell_orders"`
}
