package store

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/A-pen-app/kickstart/models"
	"github.com/A-pen-app/logging"
	"github.com/A-pen-app/tracing"
	"github.com/jmoiron/sqlx"
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

func (s *orderStore) New(ctx context.Context, userID, kickstartID string, email *string) error {
	query := `
	INSERT INTO public.user_kickstart_relation (
		user_id,
		kickstart_id,
		email,
		created_at,
		updated_at
	)
	VALUES (
		?,
		?,
		?,
		now()::timestamp(0),
		now()::timestamp(0)
	) ON CONFLICT 
	(
		user_id,
		kickstart_id
	) DO NOTHING`
	values := []interface{}{
		userID,
		kickstartID,
		email,
	}
	query = s.db.Rebind(query)
	if _, err := s.db.Exec(query, values...); err != nil {
		logging.Errorw(ctx, "attend kickstart failed", "err", err, "userID", userID, "kickstartID", kickstartID)
		return parseError(err)
	}
	return nil
}

func (s *orderStore) Get(ctx context.Context, kickstartID string) (*models.Order, error) {
	kickstart := models.Order{}
	query := `
	SELECT 
		id,
		content,
		created_at, 
		updated_at,
		ending_at,
		speaker,
		point,
		apply_info,
		charging_fee,
		participant_count,
		picture_url,
		max_participant_count,
		title,
		is_deleted,
		is_active,
		creator_id,
		tags,
		awards,
		coins
	FROM public.kickstart
	WHERE 
		id=? 
		AND 
		is_deleted=false
	`
	values := []interface{}{
		kickstartID,
	}
	query = s.db.Rebind(query)
	if err := s.db.QueryRowx(query, values...).StructScan(&kickstart); err != nil {
		logging.Errorw(ctx, "get kickstart failed", "err", err, "kickstartID", kickstartID)
		return nil, parseError(err)
	}
	return &kickstart, nil
}

func (s *orderStore) GetOrders(ctx context.Context, limit int) ([]*models.Order, error) {
	defer tracing.Start(ctx, "store.get.kickstarts").End()

	kickstarts := []*models.Order{}
	query := `
	SELECT 
		id,
		content,
		created_at, 
		updated_at,
		hosting_at,
		ending_at,
		speaker,
		point,
		apply_info,
		charging_fee,
		location,
		participant_count,
		picture_url,
		max_participant_count,
		title,
		is_deleted,
		is_active,
		creator_id,
		tags,
		contact,
		hoster,
		url,
		awards,
		coins
	FROM public.kickstart
	WHERE 
	`
	conditions := []string{
		"is_deleted=false",
		"is_active=true",
	}
	values := []interface{}{}
	query = query + strings.Join(conditions, " AND ") + " ORDER BY hosting_at DESC LIMIT ?"
	values = append(values, limit)

	query = s.db.Rebind(query)
	if err := s.db.Select(&kickstarts, query, values...); err != nil {
		logging.Errorw(ctx, "get kickstarts failed", "err", err)
		return nil, parseError(err)
	}
	return kickstarts, nil
}

func (s *orderStore) GetMyOrders(ctx context.Context, uid string, limit int, isAttending bool) ([]*models.Order, error) {
	defer tracing.Start(ctx, "store.get.my.kickstarts").End()

	kickstarts := []*models.Order{}
	query := `
	SELECT 
		id,
		content,
		created_at, 
		updated_at,
		hosting_at,
		ending_at,
		speaker,
		point,
		apply_info,
		charging_fee,
		location,
		participant_count,
		picture_url,
		max_participant_count,
		title,
		is_deleted,
		is_active,
		creator_id,
		tags,
		contact,
		hoster,
		url,
		awards,
		coins
	FROM public.kickstart m
	WHERE 
	`
	conditions := []string{
		`
		EXISTS ( 
			SELECT 1
			FROM user_kickstart_relation umr
			WHERE m.id = umr.kickstart_id
			AND umr.user_id = ?
		)
		`,
		"is_deleted=false",
		"is_active=true",
	}
	if isAttending {
		conditions = append(conditions, "hosting_at > now()::timestamp(0)")
	} else {
		conditions = append(conditions, "hosting_at < now()::timestamp(0)")
	}
	query = query + strings.Join(conditions, " AND ") + " ORDER BY hosting_at DESC LIMIT ?"

	values := []interface{}{
		uid,
		limit,
	}
	logging.Debug(ctx, fmt.Sprintf("values %v", values))

	query = s.db.Rebind(query)
	if err := s.db.Select(&kickstarts, query, values...); err != nil {
		logging.Errorw(ctx, "get kickstarts failed", "err", err)
		return nil, parseError(err)
	}
	return kickstarts, nil
}

func (s *orderStore) GetOrderIDs(ctx context.Context, uid string, limit int, status models.OrderStatus) ([]string, error) {
	defer tracing.Start(ctx, "store.get.attended.kickstarts").End()

	kickstartIDs := []string{}
	query := `
	SELECT 
		id
	FROM public.kickstart m
	WHERE 
	`
	conditions := []string{
		`
		EXISTS ( 
			SELECT 1
			FROM user_kickstart_relation umr
			WHERE m.id = umr.kickstart_id
			AND umr.user_id = ?
		)
		`,
		"is_deleted=false",
		"is_active=true",
	}
	switch status {
	case models.Attending:
		conditions = append(conditions, "hosting_at > now()::timestamp(0)")
	case models.Attended:
		conditions = append(conditions, "hosting_at < now()::timestamp(0)")
	default:
	}
	query = query + strings.Join(conditions, " AND ") + " ORDER BY created_at DESC LIMIT ?"

	values := []interface{}{
		uid,
		limit,
	}
	logging.Debug(ctx, fmt.Sprintf("values %v", values))

	query = s.db.Rebind(query)
	if err := s.db.Select(&kickstartIDs, query, values...); err != nil && err != sql.ErrNoRows {
		logging.Errorw(ctx, "get kickstarts failed", "err", err)
		return nil, parseError(err)
	}
	return kickstartIDs, nil
}

func (s *orderStore) GetIsAttended(ctx context.Context, uid, kickstartID string) (bool, error) {
	defer tracing.Start(ctx, "store.get.kickstart.is_attended").End()

	query := `
		SELECT
		EXISTS (
			SELECT 
				1
			FROM public.user_kickstart_relation
			WHERE 
	`
	conditions := []string{
		"is_deleted=false",
		"user_id=?",
		"kickstart_id=?",
	}
	query = query + strings.Join(conditions, " AND ")
	query += `)` // ending ')' for EXISTS

	values := []interface{}{
		uid,
		kickstartID,
	}

	exist := false
	query = s.db.Rebind(query)
	if err := s.db.Get(&exist, query, values...); err != nil {
		logging.Errorw(ctx, "get kickstarts failed", "err", err)
		return false, parseError(err)
	}
	return exist, nil
}
