package store

import (
	context "context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/A-pen-app/cache"
	"github.com/A-pen-app/kickstart/config"
	"github.com/A-pen-app/kickstart/database"
	"github.com/A-pen-app/kickstart/global"
	"github.com/A-pen-app/kickstart/models"
	"github.com/A-pen-app/logging"
	"github.com/A-pen-app/tracing"
	"github.com/stretchr/testify/require"
)

func TestOrderDBIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping system integration test")
	}

	os.Setenv("TESTING", "true") // to inform different parts of the application that we are testing and perform accordingly

	// Create root context.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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

	// Setup tracing module
	env := "development"
	if config.GetBool("PRODUCTION_ENVIRONMENT") {
		env = "production"
	}
	tracing.Initialize(ctx, &tracing.Config{
		ProjectID:             config.GetString("PROJECT_ID"),
		TracerName:            "kickstart-api",
		ServiceName:           global.ServiceName,
		DeploymentEnvironment: env,
	})
	defer tracing.Finalize(ctx)

	cache.Initialize(&cache.Config{
		Type:     cache.TypeLocal,
		RedisURL: "localhost:6379",
		Prefix:   "local-dev",
	})
	defer cache.Finalize()

	// Setup database module.
	database.Initialize(ctx)
	defer database.Finalize()

	db := database.GetPostgres()
	orderStore := NewOrder(db)

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
