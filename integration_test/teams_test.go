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

	resp, err := client.Teams().List()

	if err != nil {
		t.Fatalf("request err: %v", err)
	}
	t.Log("team Id:", resp.Teams[0].TeamID)
	// t.Logf("teams %+v", resp.Teams)
	// t.Logf("paged %+v", resp.Paged)

	team0, _ := json.MarshalIndent(resp, "", "  ")
	t.Log(string(team0))
}
