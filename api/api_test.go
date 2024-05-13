package api

import (
	"log"
	"os"
	"testing"

	"github.com/A-pen-app/cache"
	"github.com/A-pen-app/kickstart/config"
	"github.com/A-pen-app/logging"
)

func TestMain(m *testing.M) {
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

	log.Println("Resource Initialized. Running tests...")
	exitVal := m.Run()

	os.Exit(exitVal)
}
