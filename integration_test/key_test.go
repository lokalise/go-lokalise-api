package integration_test

import (
	"os"
	"testing"

	"github.com/lokalise/go-lokalise-api"
)

func TestKeysList(t *testing.T) {
	client, err := lokalise.New(os.Getenv("lokalise_token"))
	if err != nil {
		t.Fatalf("client instantiation: %v", err)
	}

	resp, err := client.Keys().List("3002780358964f9bab5a92.87762498", lokalise.ListKeysOptions{
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

func TestKeysCreate(t *testing.T) {
	client, err := lokalise.New(os.Getenv("lokalise_token"))
	if err != nil {
		t.Fatalf("client instantiation: %v", err)
	}

	resp, err := client.Keys().Create(
		"3002780358964f9bab5a92.87762498",
		[]lokalise.Key{{
			KeyName:   "integration_test",
			Platforms: []string{"ios", "android"},
			Translations: []lokalise.Translation{
				{
					LanguageISO: "en_US",
					Translation: "test",
				},
			},
		},
		},
	)

	if err != nil {
		t.Fatalf("request err: %v", err)
	}
	t.Logf("keys %+v", resp.Keys)
}
