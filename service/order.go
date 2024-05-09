package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/A-pen-app/cache"
	"github.com/A-pen-app/kickstart/models"
	"github.com/A-pen-app/kickstart/store"
	"github.com/A-pen-app/kickstart/util"
	"github.com/A-pen-app/logging"
	"github.com/A-pen-app/mq"
)

type orderSvc struct {
	LatestPrice int
	BoardGuard  sync.Mutex
	c           store.Order
	q           mq.MQ
}

const DEFAULT_PRICE int = 10

// NewOrder returns an implementation of service.Order
func NewOrder(c store.Order, q mq.MQ) Order {
	return &orderSvc{
		LatestPrice: DEFAULT_PRICE,
		BoardGuard:  sync.Mutex{},
		c:           c,
		q:           q,
	}
}

// return type is ([]*models.Order, string, error) corresponding to (orders, next, error)
func (s *orderSvc) GetBoard(ctx context.Context, boardType models.OrderBoardType) (*models.Board, string, error) {
	var cacheKey string
	// return value of f is (buyOrders, sellOrders, getBuyOrders, getSellOrders, error)
	var f func() (*models.Board, error)

	switch boardType {
	case models.Live:
		cacheKey = fmt.Sprintf("get_orders.%s", models.Live)
		f = func() (*models.Board, error) {
			spawned := 0
			w, errCh := sync.WaitGroup{}, make(chan error, spawned)

			buyOrders := []*models.Order{}
			spawned += 1
			w.Add(1)
			go func(ctx context.Context, orders *[]*models.Order) {
				defer w.Done()
				var err error
				if *orders, err = s.c.GetLiveOrders(ctx, models.Buy); err != nil {
					logging.Errorw(ctx, "get live buy orders failed", "err", err)
					errCh <- err
				}
			}(ctx, &buyOrders)

			sellOrders := []*models.Order{}
			spawned += 1
			w.Add(1)
			go func(ctx context.Context, orders *[]*models.Order) {
				defer w.Done()
				var err error
				if sellOrders, err = s.c.GetLiveOrders(ctx, models.Sell); err != nil {
					logging.Errorw(ctx, "get live sell orders failed", "err", err)
					errCh <- err
				}
			}(ctx, &sellOrders)

			w.Wait()
			if err := util.ChErrHandler(errCh, spawned); err != nil {
				return nil, err
			}
			return &models.Board{
				BuyOrders:  buyOrders,
				SellOrders: sellOrders,
			}, nil
		}
	case models.History:
	// FIXME: finish this
	case models.Removed:
	// FIXME: finish this
	default:
		err := fmt.Errorf("unexpected board type: %s", boardType)
		logging.Errorw(ctx, "service unexpected board type accessed in getBoard", "err", err, "boardType", boardType)
		return nil, "", err
	}

	var board *models.Board
	// here we assume cache providing the single source of truth, the integrity of the order list need to be maintained for each modification to the list
	// eg: make, take should also trigger cache update
	if err := cache.Get(ctx, cacheKey, &board); err == nil {
		if board == nil { // if last get returns [], it will be cached as null, conerting for ease of frontend integration
			board = &models.Board{
				LatestPrice: s.LatestPrice,
				BuyOrders:   []*models.Order{},
				SellOrders:  []*models.Order{},
			}
		}
		if board.BuyOrders == nil {
			board.BuyOrders = []*models.Order{}
		}
		if board.SellOrders == nil {
			board.SellOrders = []*models.Order{}
		}
	} else if err == cache.ErrorNotFound {
		logging.Infow(ctx, "get orders response not found in cache", "err", err)

		board, err = f()
		if err != nil {
			return nil, "", err
		}
		board.LatestPrice = s.LatestPrice

		// FIXME: add Board aggregation
		if err := aggregateBoard(ctx, board); err != nil {
			logging.Errorw(ctx, "aggregate orders failed", "err", err)
			return nil, "", err
		}

		// to return order list faster, use a goroutine to set cache
		go func(ctx context.Context) {
			if err := cache.SetWithTTL(ctx, cacheKey, board, time.Second); err != nil {
				logging.Errorw(ctx, "set buy orders cache failed", "err", err)
			}
		}(ctx)
	} else {
		logging.Errorw(ctx, "unexpected error while getting orders from cache", "err", err)
		return nil, "", err
	}

	// FIXME: next page token assignment
	next := ""

	return board, next, nil
}

func (s *orderSvc) Make(ctx context.Context, action models.OrderAction, price, quantity int) error {
	if action == models.Buy && price >= s.LatestPrice {
		err := fmt.Errorf("price to buy %d should not be lower than latest price %d", price, s.LatestPrice)
		logging.Errorw(ctx, "service make order failed", "err", err)
		return err
	}
	if action == models.Sell && price <= s.LatestPrice {
		err := fmt.Errorf("price to sell %d should not be higher than latest price %d", price, s.LatestPrice)
		logging.Errorw(ctx, "service make order failed", "err", err)
		return err
	}

	s.BoardGuard.Lock()
	// FIXME: should update to cache after make order
	if err := s.c.Make(ctx, action, price, quantity); err != nil {
		logging.Errorw(ctx, "service make order failed", "err", err)
		return err
	}
	s.BoardGuard.Unlock()

	go func(ctx context.Context) {
		// send email to user
		go func(ctx context.Context) {
			if err := s.q.Send("mail", struct {
				Address string
				Content string
			}{
				Address: "user@gmail.com",
				Content: "your order has been created",
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
				Content: "your order has been created",
			}); err != nil {
				logging.Errorw(ctx, "send sms failed", "err", err)
			}
		}(ctx)
	}(ctx)
	return nil
}

func (s *orderSvc) Take(ctx context.Context, action models.OrderAction, quantity int) error {
	if action != models.Buy {
		err := fmt.Errorf("currently, only buy action is supported, got %s", action)
		logging.Errorw(ctx, "service take order failed", "err", err)
		return err
	}

	// FIXME: should update to cache after take order
	s.BoardGuard.Lock()
	latestPrice, err := s.c.Take(ctx, action, quantity)
	if err != nil {
		logging.Errorw(ctx, "service take order failed", "err", err)
		return err
	}
	// in case orders are all taken
	if latestPrice == 0 {
		s.LatestPrice = DEFAULT_PRICE
	} else {
		s.LatestPrice = latestPrice
	}
	s.BoardGuard.Unlock()

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

func (s *orderSvc) Delete(ctx context.Context, orderID string) error {
	// FIXME: should update to cache after delete order
	s.BoardGuard.Lock()
	if err := s.c.Delete(ctx, orderID); err != nil {
		logging.Errorw(ctx, "attend order failed", "err", err, "orderID", orderID)
		return err
	}
	s.BoardGuard.Unlock()
	return nil
}

// this cloud be placed under service package as aggregator package
func aggregateBoard(ctx context.Context, board *models.Board) error {
	return nil
}
