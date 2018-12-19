// +build integration_test

package integration_test_test

import (
	"context"
	"flag"
	"testing"

	"github.com/17media/go-lokalise-api/lokalise"
	"github.com/17media/go-lokalise-api/model"
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
	resp, err := client.TeamUsers.Retrieve(context.Background(), 170090, 5715)
	if err != nil {
		t.Fatalf("request err: %v", err)
	}
	t.Logf("team id %d", resp.TeamID)
	t.Logf("user email %s", resp.TeamUser.Email)
}

func TestUpdateTeamUserRole(t *testing.T) {
	flag.Parse()
	if token == "" {
		t.Errorf("set token flag to run integration tests")
		return
	}
	client, err := lokalise.NewClient(token)
	if err != nil {
		t.Fatalf("client instantiation: %v", err)
	}
	resp, err := client.TeamUsers.UpdateRole(context.Background(), 170090, 5715, model.TeamUserRoleAdmin)
	if err != nil {
		t.Fatalf("request err: %v", err)
	}
	t.Logf("team id %d", resp.TeamID)
	t.Logf("user email %s", resp.TeamUser.Email)
	t.Logf("role %s", resp.TeamUser.Role)
}

func TestGetTeamUsers(t *testing.T) {
	flag.Parse()
	if token == "" {
		t.Errorf("set token flag to run integration tests")
		return
	}
	client, err := lokalise.NewClient(token, lokalise.WithRetryCount(5))
	if err != nil {
		t.Fatalf("client instantiation: %v", err)
	}

	resp, err := client.TeamUsers.List(context.Background(), 170090, lokalise.PageOptions{
		Limit: 10,
		Page:  1,
	})

	if err != nil {
		t.Fatalf("request err: %v", err)
	}
	t.Logf("team id %d", resp.TeamID)
	t.Logf("users %v", resp.TeamUsers)
	t.Logf("paged %+v", resp.Paged)
}

func TestGetTeams(t *testing.T) {
	flag.Parse()
	if token == "" {
		t.Errorf("set token flag to run integration tests")
		return
	}
	client, err := lokalise.NewClient(token, lokalise.WithRetryCount(5))
	if err != nil {
		t.Fatalf("client instantiation: %v", err)
	}

	resp, err := client.Teams.List(context.Background(), lokalise.PageOptions{
		Limit: 1,
		Page:  1,
	})

	if err != nil {
		t.Fatalf("request err: %v", err)
	}
	t.Logf("teams %+v", resp.Teams)
	t.Logf("paged %+v", resp.Paged)
}
