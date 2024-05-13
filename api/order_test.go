package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"testing"

	"github.com/A-pen-app/kickstart/models"
	"github.com/A-pen-app/kickstart/tests"
	"github.com/stretchr/testify/require"
)

func TestGetBoardAPIIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping system integration test")
	}

	if tests.BaseURL == tests.DEFAULT_DEV_DEPLOY_URL {
		t.Skip("skipping api integration test; base URL not set")
	}

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
		t.Errorf("failed to make request: %s", err.Error())

		log.Println("Parsing err response status code and result...")
		body, err := io.ReadAll(res.Body)
		if err != nil {
			t.Fatal(err.Error())
		}
		log.Printf("Err response: %s", string(body))
		return
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
