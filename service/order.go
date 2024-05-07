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

func (s *orderSvc) New(ctx context.Context, userID, kickstartID string, email *string) error {
	if err := s.c.New(ctx, userID, kickstartID, email); err != nil {
		logging.Errorw(ctx, "attend kickstart failed", "err", err, "kickstartID", kickstartID, "userID", userID)
		return err
	}
	return nil
}

func (s *orderSvc) Get(ctx context.Context, userID, kickstartID string) (*models.Order, error) {
	kickstart, err := s.c.Get(ctx, kickstartID)
	if err != nil {
		return nil, err
	}
	return kickstart, nil
}

func (s *orderSvc) GetOrders(ctx context.Context, userID string, next string, count int, filter models.OrderStatus) ([]*models.Order, string, error) {
	var cacheKey string
	var f func() ([]*models.Order, error)

	logging.Debug(ctx, string(filter))

	limit := 500
	switch filter {
	case models.Attended:
		cacheKey = fmt.Sprintf("get_kickstarts.attended.%s", userID)
		f = func() ([]*models.Order, error) {
			logging.Debug(ctx, "actually getting attending kickstarts")
			return s.c.GetMyOrders(ctx, userID, limit, false)
		}
	case models.Attending:
		logging.Debug(ctx, "getting attending kickstarts")
		cacheKey = fmt.Sprintf("get_kickstarts.attending.%s", userID)
		f = func() ([]*models.Order, error) {
			logging.Debug(ctx, "actually getting attending kickstarts")
			return s.c.GetMyOrders(ctx, userID, limit, true)
		}
	default:
		cacheKey = "get_kickstarts"
		f = func() ([]*models.Order, error) {
			return s.c.GetOrders(ctx, limit)
		}
	}

	var kickstarts []*models.Order
	if err := cache.Get(ctx, cacheKey, &kickstarts); err == nil {
		logging.Infow(ctx, "got kickstarts from cache")
		if kickstarts == nil { // if last get returns [], it will be cached as null, conerting for ease of frontend integration
			kickstarts = []*models.Order{}
		}
	} else if err == cache.ErrorNotFound {
		logging.Infow(ctx, "get kickstarts response not found in cache", "err", err)

		kickstarts, err = f()
		if err != nil {
			return nil, "", err
		}
		if err := cache.SetWithTTL(ctx, cacheKey, kickstarts, time.Minute); err != nil {
			logging.Errorw(ctx, "set kickstarts cache failed", "err", err)
			return nil, "", err
		}
	} else {
		logging.Errorw(ctx, "unexpected error while getting kickstarts from cache", "err", err)
		return nil, "", err
	}

	if next != "" {
		for i, kickstart := range kickstarts {
			if kickstart.ID == next {
				kickstarts = kickstarts[i:]
				break
			}
		}
	}

	// prepare next cursor
	next = ""
	if len(kickstarts) > count { // more elements available
		next = kickstarts[count].ID
		kickstarts = kickstarts[:count]
	}
	return kickstarts, next, nil
}

// this cloud be placed under service package as aggregator package
func aggregateOrders(ctx context.Context, kickstarts []*models.Order, attendingIDs []string, attendedIDs []string) error {
	return nil
}
