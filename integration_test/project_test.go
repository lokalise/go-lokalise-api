package integration_test

import (
	"os"
	"testing"

	"github.com/lokalise/go-lokalise-api"
)

func TestProjectsList(t *testing.T) {
	client, err := lokalise.New(os.Getenv("lokalise_token") /*lokalise.WithRetryCount(5)*/)
	if err != nil {
		t.Fatalf("client instantiation: %v", err)
	}

	resp, err := client.Projects().List(lokalise.ProjectsOptions{
		TeamID: 170090,
	})

	if err != nil {
		t.Fatalf("request err: %v", err)
	}
	t.Logf("teams %+v", resp.Projects)
	t.Logf("paged %+v", resp.Paged)
}

func TestProjectsRetrieve(t *testing.T) {
	client, err := lokalise.New(os.Getenv("lokalise_token"))
	if err != nil {
		t.Fatalf("client instantiation: %v", err)
	}
	resp, err := client.Projects().Retrieve("3002780358964f9bab5a92.87762498")
	if err != nil {
		t.Fatalf("request err: %v", err)
	}
	t.Logf("project id %s", resp.ProjectID)
	t.Logf("project name %s", resp.Name)
}
