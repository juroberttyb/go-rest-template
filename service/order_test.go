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

	orderStore.On("GetOrders", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([]*models.Order{
		{
			ID:               "7849583d-197c-48de-b48a-ce81cc26eca2",
			CreatorID:        "b7c6fc25-0cc8-4b5b-a162-2d784fa9c0d9",
			IsActive:         true,
			IsDeleted:        false,
			Status:           "attending",
			Title:            "了不起的標題",
			Content:          "這是一個了不起的活動",
			ParticipantCount: 0,
			Tags:             []string{},
		},
	}, nil)
	mq := new(mq.MockMQ)
	mq.On("GetOrderIDs", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([]string{}, nil)
	kickstartSvc := NewOrder(orderStore, mq)

	orders, _, err := kickstartSvc.GetOrders(context.Background(), "b3b646b7-7e37-4ed9-a4b3-11503b94763c", "", 50, models.Normal)
	if err != nil {
		t.Errorf("err: %s", err)
		return
	}
	orderStore.AssertExpectations(t)

	require.Equal(t, len(orders) > 0, true, "expect at least one order exists")
}

// Test all models.Order status
func TestAggregatedOrdersStatus(t *testing.T) {
	// Create a new context
	ctx := context.Background()

	// Create some orders
	orders := []*models.Order{
		{ID: "1", Tags: []string{"tag1", "tag2"}},
		{ID: "2", Tags: []string{"specialty1", "tag3"}},
		{ID: "3", Tags: []string{"tag3"}},
		{ID: "4", Tags: []string{"tag4"}},
	}

	// Call the aggregateOrders function
	err := aggregateOrders(ctx, orders)
	if err != nil {
		t.Errorf("err: %s", err)
		return
	}

	// Check the status of each kickstart
	require.Equal(t, "order created", "order created")
	require.Equal(t, "order remoevd", "order remoevd")
	require.Equal(t, "order fulfilled", "order fulfilled")
}
