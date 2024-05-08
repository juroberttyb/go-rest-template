package api

import (
	"encoding/json"
	"io"
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
