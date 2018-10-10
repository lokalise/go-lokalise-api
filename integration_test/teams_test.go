// +build integration_test

package integration_test_test

import (
	"context"
	"testing"

	"github.com/lokalise/lokalise-go-sdk/lokalise"
)

func TestGetTeamUser(t *testing.T) {
	client, err := lokalise.NewClient("20a1a5f5154f826b32eb92a3da1ea7579ce08b12")
	if err != nil {
		t.Fatalf("client instantiation: %v", err)
	}
	resp, err := client.Teams.RetrieveTeamUser(context.Background(), 178017, 5715)
	if err != nil {
		t.Fatalf("request err: %v", err)
	}
	t.Logf("team id %d", resp.TeamID)
	t.Logf("user email %s", resp.TeamUser.Email)
}
