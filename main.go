// Package main defines server initialization steps
package main

import (
	"context"
	"fmt"

	_ "go.uber.org/automaxprocs"

	"github.com/A-pen-app/cache"
	"github.com/A-pen-app/kickstart/config"
	"github.com/A-pen-app/kickstart/database"
	"github.com/A-pen-app/kickstart/global"
	"github.com/A-pen-app/kickstart/server/app"
	"github.com/A-pen-app/logging"
	"github.com/A-pen-app/mq"
	"github.com/A-pen-app/mq/pubsubLite"
	"github.com/A-pen-app/mq/rabbitmq"
	"github.com/A-pen-app/tracing"
)

func main() {
	// We're running, turn on the liveness indication flag.
	global.Alive = true

	// Create root context.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var projectID string = config.GetString("PROJECT_ID")
	if !config.GetBool("PRODUCTION_ENVIRONMENT") {
		projectID = ""
	}

	// Setup logging module.
	// NOTE: This should always be first.
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

	// Setup cache module
	//FIXME set it to production for testing purpose
	cacheType := cache.TypeLocal
	prefix := "mp-dev"
	redisURL := "localhost:6379"

	if config.GetBool("PRODUCTION_ENVIRONMENT") {
		cacheType = cache.TypeRedis
		redisURL = "10.49.162.163:6379"
		prefix = "kickstart"
	}
	cache.Initialize(&cache.Config{
		Type:     cacheType,
		RedisURL: redisURL,
		Prefix:   prefix,
	})
	defer cache.Finalize()

	// Setup database module.
	database.Initialize(ctx)
	defer database.Finalize()

	// Setup mq module.
	mq.Initialize(ctx, &mq.Config{
		Pubsub: &pubsubLite.Config{
			ProjectID:    config.GetString("PROJECT_ID"),
			RegionOrZone: "asia-east1",
			Topics: map[string]string{
				"svc-action": "",
				"notif":      "",
				"consume":    "consume",
			},
		},
		Rabbitmq: &rabbitmq.Config{
			ProjectName:     config.GetString("PROJECT_NAME"), // FIXME: need to add a listener at mq-svc to listen and handle this PROJECT_NAME's routing key
			RabbitmqConnURL: "amqps://zybsafyg:gik8pM6R_3FZv6EUQHOPchouyLqO9sj5@mustang.rmq.cloudamqp.com/zybsafyg",
		},
	})
	defer mq.Finalize()

	// Create HTTP server instance to listen on all interfaces.
	address := fmt.Sprintf("%s:%s",
		config.GetString("SERVER_LISTEN_ADDRESS"),
		config.GetString("SERVER_LISTEN_PORT"))
	server := app.CreateServer(ctx, address)

	// Now that we finished initializing all necessary modules,
	// let's turn on the readiness indication flag.
	global.Ready = true

	// Start servicing requests.
	logging.Info(ctx, "Initialization complete, listening on %s...", address)
	if err := server.ListenAndServe(); err != nil {
		logging.Info(ctx, err.Error())
	}
}
