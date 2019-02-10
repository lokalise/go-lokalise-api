// +build integration_test

package integration_test_test

import (
	"context"
	"flag"
	"fmt"
	"testing"

	"github.com/lokalise/go-lokalise-api/lokalise"
	"github.com/lokalise/go-lokalise-api/model"
)

func TestListKeys(t *testing.T) {
	flag.Parse()
	if token == "" {
		t.Errorf("set token flag to run integration tests")
		return
	}
	client, err := lokalise.NewClient(token, lokalise.WithRetryCount(5))
	if err != nil {
		t.Fatalf("client instantiation: %v", err)
	}

	resp, err := client.Keys.List(context.Background(), "3002780358964f9bab5a92.87762498", lokalise.ListKeysOptions{
		IncludeTranslations: true,
		PageOptions: lokalise.PageOptions{
			Limit: 10,
			Page:  1,
		},
	})

	if err != nil {
		t.Fatalf("request err: %v", err)
	}
	t.Logf("keys %+v", resp.Keys)
	t.Logf("paged %+v", resp.Paged)
}

func TestCreateKeys(t *testing.T) {
	flag.Parse()
	if token == "" {
		t.Errorf("set token flag to run integration tests")
		return
	}
	client, err := lokalise.NewClient(token, lokalise.WithRetryCount(5))
	if err != nil {
		t.Fatalf("client instantiation: %v", err)
	}

	resp, err := client.Keys.Create(context.Background(), "3002780358964f9bab5a92.87762498", []lokalise.KeyOptions{lokalise.KeyOptions{
		KeyName:   "integration_test",
		Platforms: []string{"ios", "android"},
		Translations: []model.Translation{
			model.Translation{
				LanguageISO: "en_US",
				Translation: "test",
			},
		},
	},
	})

	if err != nil {
		t.Fatalf("request err: %v", err)
	}
	t.Logf("keys %+v", resp.Keys)
}
