package service

import (
	"context"
	"testing"
	"time"

	"github.com/A-pen-app/kickstart/models"
	"github.com/A-pen-app/kickstart/mq"
	"github.com/A-pen-app/kickstart/store"
	mock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetOrders(t *testing.T) {
	orderStore := new(store.MockOrder)

	layout := "2006-01-02 15:04:05"
	createdAt, err := time.Parse(layout, "2024-01-22 04:06:20")
	if err != nil {
		t.Errorf("err: %s", err)
		return
	}
	t.Log("test 2")
	updatedAt, err := time.Parse(layout, "2024-01-22 04:06:20")
	if err != nil {
		t.Errorf("err: %s", err)
		return
	}
	hostingAt, err := time.Parse(layout, "2025-01-22 05:38:39")
	if err != nil {
		t.Errorf("err: %s", err)
		return
	}
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
			CreatedAt:        createdAt,
			UpdatedAt:        updatedAt,
			HostingAt:        hostingAt,
		},
	}, nil)
	orderStore.On("GetOrderIDs", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([]string{}, nil)
	mq := new(mq.MockMQ)
	kickstartSvc := NewOrder(orderStore, mq)

	kickstarts, _, err := kickstartSvc.GetOrders(context.Background(), "b3b646b7-7e37-4ed9-a4b3-11503b94763c", "", 50, models.Normal)
	if err != nil {
		t.Errorf("err: %s", err)
		return
	}
	orderStore.AssertExpectations(t)

	for _, m := range kickstarts {
		require.Equal(t, m.Status, models.Normal, "status should be normal")
	}
}

// Test all models.OrderStatus cases
func TestAggregateOrdersAllStatus(t *testing.T) {
	// Create a new context
	ctx := context.Background()

	// Create some kickstarts
	kickstarts := []*models.Order{
		{ID: "1", Tags: []string{"tag1", "tag2"}},
		{ID: "2", Tags: []string{"specialty1", "tag3"}},
		{ID: "3", Tags: []string{"tag3"}},
		{ID: "4", Tags: []string{"tag4"}},
	}

	// Create some attending IDs
	attendingIDs := []string{"1"}
	// Create some attended IDs
	attendedIDs := []string{"3"}

	// Call the aggregateOrders function
	err := aggregateOrders(ctx, kickstarts, attendingIDs, attendedIDs)
	if err != nil {
		t.Errorf("err: %s", err)
		return
	}

	// Check the status of each kickstart
	require.Equal(t, models.Attending, kickstarts[0].Status)
	require.Equal(t, models.Recommended, kickstarts[1].Status)
	require.Equal(t, models.Attended, kickstarts[2].Status)
	require.Equal(t, models.Normal, kickstarts[3].Status)
}
