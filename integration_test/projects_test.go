// +build integration_test

package integration_test_test

import (
	"context"
	"flag"
	"testing"

	"github.com/17media/go-lokalise-api/lokalise"
)

func TestGetProjects(t *testing.T) {
	flag.Parse()
	if token == "" {
		t.Errorf("set token flag to run integration tests")
		return
	}
	client, err := lokalise.NewClient(token, lokalise.WithRetryCount(5))
	if err != nil {
		t.Fatalf("client instantiation: %v", err)
	}

	resp, err := client.Projects.List(context.Background(), lokalise.ProjectsOptions{
		TeamID: 170090,
		PageOptions: lokalise.PageOptions{
			Limit: 0,
			Page:  1,
		},
	})

	if err != nil {
		t.Fatalf("request err: %v", err)
	}
	t.Logf("teams %+v", resp.Projects)
	t.Logf("paged %+v", resp.Paged)
}

func TestGetProject(t *testing.T) {
	flag.Parse()
	if token == "" {
		t.Errorf("set token flag to run integration tests")
		return
	}
	client, err := lokalise.NewClient(token)
	if err != nil {
		t.Fatalf("client instantiation: %v", err)
	}
	resp, err := client.Projects.Retrieve(context.Background(), "3002780358964f9bab5a92.87762498")
	if err != nil {
		t.Fatalf("request err: %v", err)
	}
	t.Logf("project id %s", resp.ProjectID)
	t.Logf("project name %s", resp.Name)
}
