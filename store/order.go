package store

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/A-pen-app/kickstart/config"
	"github.com/A-pen-app/kickstart/database"
	"github.com/A-pen-app/kickstart/models"
	"github.com/A-pen-app/logging"
	"github.com/A-pen-app/tracing"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type orderStore struct {
	db *sqlx.DB
}

// NewOrder returns an implementation of store.Order
func NewOrder(db *sqlx.DB) Order {
	return &orderStore{
		db: db,
	}
}

func (s *orderStore) GetLiveOrders(ctx context.Context, action models.OrderAction) ([]*models.Order, error) {
	if !config.GetBool("TESTING") {
		defer tracing.Start(ctx, "store.get.orders.live").End()
	}

	orders := []*models.Order{}
	query := `
		SELECT 
			id,
			action,
			price,
			quantity,
			created_at
		FROM public.order
		WHERE 
	`
	conditions := []string{
		"action = ?",
	}
	values := []interface{}{
		action,
	}
	query = query + strings.Join(conditions, " AND ") + " ORDER BY created_at DESC"

	switch action {
	case models.Buy:
		query = query + ", price ASC" // (latest order, highest price)
	case models.Sell:
		query = query + ", price DESC" // (latest order, lowest price)
	default:
		err := errors.New("invalid order action")
		logging.Errorw(ctx, "store get live orders failed", "err", err)
		return nil, err
	}

	query = s.db.Rebind(query)
	if err := s.db.Select(&orders, query, values...); err != nil {
		if err == sql.ErrNoRows {
			return orders, nil
		}
		logging.Errorw(ctx, "store get orders failed", "err", err)
		return nil, parseError(err)
	}
	return orders, nil
}

func (s *orderStore) Make(ctx context.Context, action models.OrderAction, price, quantity int) error {
	query := `
		INSERT INTO public.order (
			action,
			price,
			quantity
		)
		VALUES (
			?,
			?,
			?
		)
	`
	values := []interface{}{
		action,
		price,
		quantity,
	}
	query = s.db.Rebind(query)
	if _, err := s.db.Exec(query, values...); err != nil {
		logging.Errorw(ctx, "store make order failed", "err", err)
		return parseError(err)
	}
	return nil
}

// FIXME: currently only buy action is supported
func (s *orderStore) Take(ctx context.Context, action models.OrderAction, quantity int) (int, error) {
	var latestPrice int

	db := database.GetPostgres()
	if err := database.Transaction(db, func(tx *sqlx.Tx) error {
		orders := []*models.Order{}
		query := `
			SELECT 
				id,
				price,
				quantity,
				created_at
			FROM public.order
			WHERE 
		`
		conditions := []string{
			"action = ?",
		}
		values := []interface{}{
			models.Sell,
		}
		query = query + strings.Join(conditions, " AND ") + " ORDER BY created_at DESC, price ASC" // (latest order, lowest price)

		query = tx.Rebind(query)
		if err := tx.Select(&orders, query, values...); err != nil {
			logging.Errorw(ctx, "store get orders failed", "err", err)
			return parseError(err)
		}

		orderIDs := []string{}
		// FIXME: currently this take orders until quantity is 0, but if selling orders are not enough, it should be handled
		for _, order := range orders {
			if quantity == 0 {
				break
			}
			if order.Quantity > quantity {
				query := `
					UPDATE public.order
					SET
						quantity=?
					WHERE
					id=?
				`
				values := []interface{}{
					order.Quantity - quantity,
					order.ID,
				}
				query = tx.Rebind(query)
				if _, err := tx.Exec(query, values...); err != nil {
					logging.Errorw(ctx, "store update order in take order failed", "err", err)
					return parseError(err)
				}
				quantity = 0
				latestPrice = order.Price
			} else {
				orderIDs = append(orderIDs, order.ID)
				quantity = quantity - order.Quantity
			}
		}

		if len(orderIDs) > 0 {
			query := `
				DELETE FROM public.order 
				WHERE 
				id = ANY(?)
			`
			values := []interface{}{
				pq.StringArray(orderIDs),
			}
			query = tx.Rebind(query)
			if _, err := tx.Exec(query, values...); err != nil {
				logging.Errorw(ctx, "store delete taken orders failed", "err", err)
				return parseError(err)
			}
		}
		return nil
	}); err != nil {
		logging.Errorw(ctx, "store take order action failed", "err", err)
		return 0, err
	}
	return latestPrice, nil
}

func (s *orderStore) Delete(ctx context.Context, orderID string) error {
	query := `
		DELETE FROM public.order 
		WHERE 
		id = ?
	`
	values := []interface{}{
		orderID,
	}
	query = s.db.Rebind(query)
	if _, err := s.db.Exec(query, values...); err != nil {
		logging.Errorw(ctx, "store delete order failed", "err", err, "orderID", orderID)
		return parseError(err)
	}
	return nil
}
