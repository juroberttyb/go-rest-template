package store

import (
	context "context"
	"fmt"
	"testing"

	"github.com/A-pen-app/kickstart/database"
	"github.com/A-pen-app/kickstart/models"
	"github.com/stretchr/testify/require"
)

func TestOrderDBIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping system integration test")
	}

	db := database.GetPostgres()
	orderStore := NewOrder(db)

	ctx := context.Background()
	sellPrice := 50
	err := orderStore.Make(ctx, models.Sell, sellPrice, 10)
	if err != nil {
		t.Fatalf("make sell order failed: %s", err.Error())
	}

	err = orderStore.Make(ctx, models.Buy, 5, 20)
	if err != nil {
		t.Fatalf("make buy order failed: %s", err.Error())
	}

	newPrice, err := orderStore.Take(ctx, models.Buy, 2)
	if err != nil {
		t.Fatalf("take buy order failed: %s", err.Error())
	}
	require.Equal(t, sellPrice, newPrice, fmt.Sprintf("expect new price to be %d, the default price", sellPrice))

	buyOrders, err := orderStore.GetLiveOrders(ctx, models.Buy)
	if err != nil {
		t.Fatalf("get live buy orders failed: %s", err.Error())
	}
	require.Equal(t, true, len(buyOrders) > 0, "expect at least 1 buy order")

	sellOrders, err := orderStore.GetLiveOrders(ctx, models.Sell)
	if err != nil {
		t.Fatalf("get live sell orders failed: %s", err.Error())
	}
	require.Equal(t, true, len(sellOrders) > 0, "expect at least 1 sell order")
	require.Equal(t, 8, sellOrders[0].Quantity, "expect sell order quantity to be 8")
}
