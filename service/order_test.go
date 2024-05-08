package service

import (
	"context"
	"testing"

	"github.com/A-pen-app/kickstart/models"
	"github.com/A-pen-app/kickstart/mq"
	"github.com/A-pen-app/kickstart/store"
	mock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetOrders(t *testing.T) {
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
