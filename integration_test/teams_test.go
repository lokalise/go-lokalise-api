// +build integration_test

package integration_test_test

import (
	"context"
	"flag"
	"testing"

	"github.com/lokalise/lokalise-go-sdk/lokalise"
)

var (
	token string
)

func init() {
	flag.StringVar(&token, "token", "", "Lokalise API token")
}

func TestGetTeamUser(t *testing.T) {
	flag.Parse()
	if token == "" {
		t.Errorf("set token flag to run integration tests")
		return
	}
	client, err := lokalise.NewClient(token)
	if err != nil {
		t.Fatalf("client instantiation: %v", err)
	}
	resp, err := client.TeamUsers.Retrieve(context.Background(), 178017, 5715)
	if err != nil {
		t.Fatalf("request err: %v", err)
	}
	t.Logf("team id %d", resp.TeamID)
	t.Logf("user email %s", resp.TeamUser.Email)
}
