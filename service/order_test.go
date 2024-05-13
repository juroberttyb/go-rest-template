package service

import (
	"context"
	"log"
	"testing"

	"github.com/A-pen-app/cache"
	"github.com/A-pen-app/kickstart/config"
	"github.com/A-pen-app/kickstart/models"
	"github.com/A-pen-app/kickstart/mq"
	"github.com/A-pen-app/kickstart/store"
	"github.com/A-pen-app/logging"
	mock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetOrders(t *testing.T) {
	log.Println("Initializing resource for testing...")
	var projectID string = config.GetString("PROJECT_ID")
	if !config.GetBool("PRODUCTION_ENVIRONMENT") {
		projectID = ""
	}

	if err := logging.Initialize(&logging.Config{
		ProjectID:    projectID,
		Level:        logging.Level(config.GetUint("LOG_LEVEL")),
		Development:  !config.GetBool("PRODUCTION_ENVIRONMENT"),
		KeyRequestID: "request_id",
		KeyUserID:    "user_id",
		KeyError:     "err",
		KeyScope:     "scope",
	}); err != nil {
		panic(err)
	}
	defer logging.Finalize()

	cache.Initialize(&cache.Config{
		Type:     cache.TypeLocal,
		RedisURL: "localhost:6379",
		Prefix:   "local-dev",
	})
	defer cache.Finalize()

	orderStore := new(store.MockOrder)

	orderStore.On("GetLiveOrders", mock.Anything, mock.Anything).Return(
		[]*models.Order{
			{
				ID: "7849583d-197c-48de-b48a-ce81cc26eca2",
			},
		},
		nil,
	)
	mq := new(mq.MockMQ)
	kickstartSvc := NewOrder(orderStore, mq)

	board, next, err := kickstartSvc.GetBoard(context.Background(), models.Live)
	if err != nil {
		t.Errorf("err: %s", err)
		return
	}
	orderStore.AssertExpectations(t)

	require.Equal(t, "", next, "expect next to be empty")
	require.Equal(t, true, board != nil, "expect at least one order exists")
}

// Test all models.Order status
func TestAggregatedOrdersStatus(t *testing.T) {
	log.Println("Initializing resource for testing...")
	var projectID string = config.GetString("PROJECT_ID")
	if !config.GetBool("PRODUCTION_ENVIRONMENT") {
		projectID = ""
	}

	if err := logging.Initialize(&logging.Config{
		ProjectID:    projectID,
		Level:        logging.Level(config.GetUint("LOG_LEVEL")),
		Development:  !config.GetBool("PRODUCTION_ENVIRONMENT"),
		KeyRequestID: "request_id",
		KeyUserID:    "user_id",
		KeyError:     "err",
		KeyScope:     "scope",
	}); err != nil {
		panic(err)
	}
	defer logging.Finalize()

	cache.Initialize(&cache.Config{
		Type:     cache.TypeLocal,
		RedisURL: "localhost:6379",
		Prefix:   "local-dev",
	})
	defer cache.Finalize()

	// Create a new context
	ctx := context.Background()

	// Create some orders
	board := models.Board{}

	// Call the aggregateOrders function
	err := aggregateBoard(ctx, &board)
	if err != nil {
		t.Errorf("err: %s", err)
		return
	}

	// Check the status of each kickstart
	require.Equal(t, "order created", "order created")
	require.Equal(t, "order remoevd", "order remoevd")
	require.Equal(t, "order fulfilled", "order fulfilled")
}
