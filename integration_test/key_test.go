package integration_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/lokalise/go-lokalise-api"
)

/*func TestKeyList(t *testing.T) {
	client, err := lokalise.New(os.Getenv("lokalise_token"))
	if err != nil {
		t.Fatalf("client instantiation: %v", err)
	}

	keys := client.Keys()
	keys.SetDebug(true)
	keys.SetListOptions(lokalise.KeyListOptions{
		Limit:               3,
		IncludeTranslations: 1,
	})

	resp, err := keys.List("373182575d64e892ba8ab2.58226357")

	if err != nil {
		t.Fatalf("request err: %v", err)
	}

	respJson, _ := json.MarshalIndent(resp, "", "  ")
	t.Log("\n", string(respJson))
}*/

func TestKeyRetrieve(t *testing.T) {
	client, err := lokalise.New(os.Getenv("lokalise_token"))
	if err != nil {
		t.Fatalf("client instantiation: %v", err)
	}

	keys := client.Keys()
	// keys.SetDebug(true)

	resp, err := keys.Retrieve("373182575d64e892ba8ab2.58226357", 26835175)

	if err != nil {
		t.Fatalf("request err: %v", err)
	}

	respJson, _ := json.MarshalIndent(resp, "", "  ")
	t.Log("\n", string(respJson))
}

/*func TestKeyCreate(t *testing.T) {
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
}*/
