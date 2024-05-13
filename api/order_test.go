package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"testing"

	"github.com/A-pen-app/cache"
	"github.com/A-pen-app/kickstart/config"
	"github.com/A-pen-app/kickstart/models"
	"github.com/A-pen-app/kickstart/tests"
	"github.com/A-pen-app/logging"
	"github.com/stretchr/testify/require"
)

func TestGetBoardAPIIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping system integration test")
	}

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

	url := tests.BaseURL + "/board?board_type=live"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatal(err.Error())
	}

	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		t.Fatal(err.Error())
	}
	defer res.Body.Close()

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
}
