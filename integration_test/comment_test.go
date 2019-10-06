package integration_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/lokalise/go-lokalise-api"
)

func TestCommentsList(t *testing.T) {
	client, err := lokalise.New(os.Getenv("lokalise_token"))
	if err != nil {
		t.Fatalf("client instantiation: %v", err)
	}

	comments := client.Comments()
	comments.SetDebug(true)

	resp, err := comments.ListProject("373182575d64e892ba8ab2.58226357")

	if err != nil {
		t.Fatalf("request err: %v", err)
	}
	// t.Logf("teams %+v", resp.Teams)
	// t.Logf("paged %+v", resp.Paged)

	respJson, _ := json.MarshalIndent(resp, "", "  ")
	t.Log("\n", string(respJson))
}
