package service

import (
	"context"
	"fmt"
	"time"

	"github.com/A-pen-app/cache"
	"github.com/A-pen-app/kickstart/models"
	"github.com/A-pen-app/kickstart/store"
	"github.com/A-pen-app/logging"
	"github.com/A-pen-app/mq"
)

type orderSvc struct {
	c store.Order
	q mq.MQ
}

// NewOrder returns an implementation of service.Order
func NewOrder(c store.Order, q mq.MQ) Order {
	return &orderSvc{
		c: c,
		q: q,
	}
}

func (s *orderSvc) New(ctx context.Context, userID, orderID string, email *string) error {
	if err := s.c.New(ctx, userID, orderID, email); err != nil {
		logging.Errorw(ctx, "attend order failed", "err", err, "orderID", orderID, "userID", userID)
		return err
	}
	return nil
}

func (s *orderSvc) Take(ctx context.Context, userID, orderID string, email *string) error {
	if err := s.c.New(ctx, userID, orderID, email); err != nil {
		logging.Errorw(ctx, "attend order failed", "err", err, "orderID", orderID, "userID", userID)
		return err
	}

	go func(ctx context.Context) {
		// send email to user
		go func(ctx context.Context) {
			if err := s.q.Send("mail", struct {
				Address string
				Content string
			}{
				Address: "user@gmail.com",
				Content: "your order has been fulfilled",
			}); err != nil {
				logging.Errorw(ctx, "send email failed", "err", err)
			}
		}(ctx)

		// send sms message to user
		go func(ctx context.Context) {
			if err := s.q.Send("sms", struct {
				Number  string
				Content string
			}{
				Number:  "0911122233",
				Content: "your order has been fulfilled",
			}); err != nil {
				logging.Errorw(ctx, "send sms failed", "err", err)
			}
		}(ctx)
	}(ctx)
	return nil
}

func (s *orderSvc) Get(ctx context.Context, userID, orderID string) (*models.Order, error) {
	order, err := s.c.Get(ctx, orderID)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (s *orderSvc) GetOrders(ctx context.Context, userID string, next string, count int, filter models.OrderStatus) ([]*models.Order, string, error) {
	var cacheKey string
	var f func() ([]*models.Order, error)

	logging.Debug(ctx, string(filter))

	limit := 500
	switch filter {
	case models.Attended:
		cacheKey = fmt.Sprintf("get_orders.attended.%s", userID)
		f = func() ([]*models.Order, error) {
			logging.Debug(ctx, "actually getting attending orders")
			return s.c.GetMyOrders(ctx, userID, limit, false)
		}
	case models.Attending:
		logging.Debug(ctx, "getting attending orders")
		cacheKey = fmt.Sprintf("get_orders.attending.%s", userID)
		f = func() ([]*models.Order, error) {
			logging.Debug(ctx, "actually getting attending orders")
			return s.c.GetMyOrders(ctx, userID, limit, true)
		}
	default:
		cacheKey = "get_orders"
		f = func() ([]*models.Order, error) {
			return s.c.GetOrders(ctx, limit)
		}
	}

	var orders []*models.Order
	if err := cache.Get(ctx, cacheKey, &orders); err == nil {
		logging.Infow(ctx, "got orders from cache")
		if orders == nil { // if last get returns [], it will be cached as null, conerting for ease of frontend integration
			orders = []*models.Order{}
		}
	} else if err == cache.ErrorNotFound {
		logging.Infow(ctx, "get orders response not found in cache", "err", err)

		orders, err = f()
		if err != nil {
			return nil, "", err
		}
		if err := cache.SetWithTTL(ctx, cacheKey, orders, time.Minute); err != nil {
			logging.Errorw(ctx, "set orders cache failed", "err", err)
			return nil, "", err
		}
	} else {
		logging.Errorw(ctx, "unexpected error while getting orders from cache", "err", err)
		return nil, "", err
	}

	if next != "" {
		for i, order := range orders {
			if order.ID == next {
				orders = orders[i:]
				break
			}
		}
	}

	// prepare next cursor
	next = ""
	if len(orders) > count { // more elements available
		next = orders[count].ID
		orders = orders[:count]
	}

	return orders, next, nil
}

// this cloud be placed under service package as aggregator package
func aggregateOrders(ctx context.Context, orders []*models.Order) error {
	return nil
}
