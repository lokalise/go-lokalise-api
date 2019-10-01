package integration_test

import (
	"os"
	"testing"

	"github.com/lokalise/go-lokalise-api"
)

func TestTeamUsersRetrieve(t *testing.T) {
	client, err := lokalise.New(os.Getenv("lokalise_token"))
	if err != nil {
		t.Fatalf("client instantiation: %v", err)
	}

	resp, err := client.TeamUsers().Retrieve(170090, 5715)
	if err != nil {
		t.Fatalf("request err: %v", err)
	}
	t.Logf("team id %d", resp.TeamID)
	t.Logf("user email %s", resp.TeamUser.Email)
}

func TestTeamUsersUpdateRole(t *testing.T) {
	client, err := lokalise.New(os.Getenv("lokalise_token"))
	if err != nil {
		t.Fatalf("client instantiation: %v", err)
	}

	resp, err := client.TeamUsers().UpdateRole(170090, 5715, lokalise.TeamUserRoleAdmin)
	if err != nil {
		t.Fatalf("request err: %v", err)
	}
	t.Logf("team id %d", resp.TeamID)
	t.Logf("user email %s", resp.TeamUser.Email)
	t.Logf("role %s", resp.TeamUser.Role)
}

func TestTeamUsersList(t *testing.T) {
	client, err := lokalise.New(os.Getenv("lokalise_token"), lokalise.WithRetryCount(5))
	if err != nil {
		t.Fatalf("client instantiation: %v", err)
	}

	resp, err := client.TeamUsers().List(170090)

	if err != nil {
		t.Fatalf("request err: %v", err)
	}
	t.Logf("team id %d", resp.TeamID)
	t.Logf("users %v", resp.TeamUsers)
	t.Logf("paged %+v", resp.Paged)
}
