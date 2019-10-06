package integration_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/lokalise/go-lokalise-api"
)

func TestProjectsList(t *testing.T) {
	client, err := lokalise.New(os.Getenv("lokalise_token") /*lokalise.WithRetryCount(5)*/)
	if err != nil {
		t.Fatalf("client instantiation: %v", err)
	}

	projects := client.Projects()
	projects.SetDebug(true)
	// projects.SetListOptions(...)

	resp, err := projects.
		WithListOptions(lokalise.ProjectListOptions{
			IncludeStat:     "0", // should set statistics to null
			IncludeSettings: "0", // should set settings to null
			Page:            0,   // shouldn't affect to request url
		}).
		List()

	if err != nil {
		t.Fatalf("request err: %v", err)
	}

	respJson, _ := json.MarshalIndent(resp, "", "  ")
	t.Log("\n", string(respJson))
}

// func TestProjectsRetrieve(t *testing.T) {
// 	client, err := lokalise.New(os.Getenv("lokalise_token"))
// 	if err != nil {
// 		t.Fatalf("client instantiation: %v", err)
// 	}
// 	resp, err := client.Projects().Retrieve("373182575d64e892ba8ab2.58226357")
// 	if err != nil {
// 		t.Fatalf("request err: %v", err)
// 	}
// 	t.Logf("project id %s", resp.ProjectID)
// 	t.Logf("project name %s", resp.Name)
// }
