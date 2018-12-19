// +build integration_test

package integration_test_test

import (
	"context"
	"flag"
	"fmt"
	"testing"

	"github.com/17media/go-lokalise-api/lokalise"
	"github.com/17media/go-lokalise-api/model"
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

	resp, err := client.Keys.List(context.Background(), "8749166159795d02ca9bd6.54233877", lokalise.ListKeysOptions{
		IncludeTranslations: true,
		PageOptions: lokalise.PageOptions{
			Limit: 10,
			Page:  1,
		},
	})

	if err != nil {
		t.Fatalf("request err: %v", err)
	}
	fmt.Println(resp.Keys)
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

	resp, err := client.Keys.Create(context.Background(), "8749166159795d02ca9bd6.54233877", []lokalise.KeyOptions{lokalise.KeyOptions{
		KeyName:   "test_from_readper",
		Platforms: []string{"ios", "android"},
		Translations: []model.Translation{
			model.Translation{
				LanguageISO: "en",
				Translation: "test",
			},
			model.Translation{
				LanguageISO: "id",
				Translation: "testID",
			},
		},
	},
	})

	if err != nil {
		t.Fatalf("request err: %v", err)
	}
	t.Logf("keys %+v", resp.Keys)
}
