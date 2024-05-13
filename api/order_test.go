package api

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/A-pen-app/cache"
	"github.com/A-pen-app/kickstart/config"
	"github.com/A-pen-app/kickstart/database"
	"github.com/A-pen-app/kickstart/global"
	"github.com/A-pen-app/kickstart/models"
	"github.com/A-pen-app/kickstart/tests"
	"github.com/A-pen-app/logging"
	"github.com/A-pen-app/tracing"
	"github.com/stretchr/testify/require"
)

func TestGetBoardAPIIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping system integration test")
	}

	os.Setenv("TESTING", "true") // to inform different parts of the application that we are testing and perform accordingly

	log.Println("Initializing resource for testing...")
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
		TracerName:            "kickstart",
		ServiceName:           global.ServiceName,
		DeploymentEnvironment: env,
	})
	defer tracing.Finalize(ctx)

	// Setup cache module
	//FIXME set it to production for testing purpose
	cacheType := cache.TypeLocal
	prefix := "local-dev"
	redisURL := "localhost:6379"

	if config.GetBool("PRODUCTION_ENVIRONMENT") {
		cacheType = cache.TypeRedis
		redisURL = "10.49.162.163:6379"
		prefix = config.GetString("SERVICE_NAME")
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

	// Now that we finished initializing all necessary modules,
	// let's turn on the readiness indication flag.
	global.Ready = true

	log.Println("Resource for testing initialized...")
	url := tests.BaseURL + "/board?board_type=live"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatal(err.Error())
	}

	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	client := &http.Client{}

	log.Println("API request prepared, making request to server...")
	res, err := client.Do(req)
	if err != nil {
		t.Fatal(err.Error())
	}
	defer res.Body.Close()
	log.Println("Request made, parsing response status code and result...")

	if res.StatusCode != http.StatusOK {
		t.Fatalf("status code %d != 200", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err.Error())
	}

	p := pageResp{}
	if err = json.Unmarshal(body, &p); err != nil {
		t.Fatal(err.Error())
	}
	rawBoard, err := json.Marshal(p.Data)
	if err != nil {
		t.Fatal(err.Error())
	}
	board := models.Board{}
	if err = json.Unmarshal(rawBoard, &board); err != nil {
		t.Fatal(err.Error())
	}
	t.Logf("LatestPrice %d", board.LatestPrice)
	t.Logf("Board: %+v", board)

	require.Equal(t, true, board.LatestPrice != 0, "expect latest price to be non-zero")

	log.Println("GetBoard API test finishing...")
}
