package integration_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/lokalise/go-lokalise-api"
)

func TestScreenshotList(t *testing.T) {
	client, err := lokalise.New(os.Getenv("lokalise_token"))
	if err != nil {
		t.Fatalf("client instantiation: %v", err)
	}

	scr := client.Screenshots()
	scr.SetDebug(true)

	resp, err := scr.List("373182575d64e892ba8ab2.58226357")

	if err != nil {
		t.Fatalf("request err: %v", err)
	}

	respJson, _ := json.MarshalIndent(resp, "", "  ")
	t.Log("\n", string(respJson))
}
