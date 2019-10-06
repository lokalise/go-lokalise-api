package integration_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/lokalise/go-lokalise-api"
)

func TestTeamsList(t *testing.T) {
	client, err := lokalise.New(os.Getenv("lokalise_token"))
	if err != nil {
		t.Fatalf("client instantiation: %v", err)
	}

	teams := client.Teams()
	teams.SetDebug(true)
	// teams.SetPageOptions(lokalise.PageOptions{Limit: 50})

	resp, err := teams.WithListOptions(lokalise.PageOptions{Limit: 50}).List()

	if err != nil {
		t.Fatalf("request err: %v", err)
	}
	// t.Logf("teams %+v", resp.Teams)
	// t.Logf("paged %+v", resp.Paged)

	respJson, _ := json.MarshalIndent(resp, "", "  ")
	t.Log("\n", string(respJson))
}
