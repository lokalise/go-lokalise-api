package lokalise

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestTranslationProviderService_List(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(
		fmt.Sprintf("/teams/%d/translation_providers", 444),
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			testMethod(t, r, "GET")
			testHeader(t, r, apiTokenHeader, testApiToken)

			_, _ = fmt.Fprint(w, `{
				"translation_providers": [
					{
						"provider_id": 1
					}
				]
			}`)
		})

	r, err := client.TranslationProviders().List(444)
	if err != nil {
		t.Errorf("TranslationProviders.List returned error: %v", err)
	}

	want := []TranslationProvider{
		{
			ProviderID: 1,
		},
	}

	if !reflect.DeepEqual(r.TranslationProviders, want) {
		t.Errorf("TranslationProviders.List returned %+v, want %+v", r.TranslationProviders, want)
	}
}

func TestTranslationProviderService_Retrieve(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(
		fmt.Sprintf("/teams/%d/translation_providers/%d", 444, 1),
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			testMethod(t, r, "GET")
			testHeader(t, r, apiTokenHeader, testApiToken)

			_, _ = fmt.Fprint(w, `{
				"provider_id": 1,
				"name": "Gengo",
				"slug": "gengo",
				"price_pair_min": "0.00",
				"website_url": "https://gengo.com",
				"description": "At Gengo, our mission is to provide language services to everyone and connect a global community. Our network of over 18,000 translators are tested and qualified to meet stringent project standards. Translators are based around the world and in every timezone, which enables us to support over 35 languages and work towards the Gengo mission.",
				"tiers": [
					{
						"tier_id": 1,
						"title": "Native speaker"
					}
				],
				"pairs": [
					{
						"tier_id": 1,
						"from_lang_iso": "en_US",
						"from_lang_name": "English (United States)",
						"to_lang_iso": "ru",
						"to_lang_name": "Russian",
						"price_per_word": "0.07"
					}
				]
			}`)
		})

	r, err := client.TranslationProviders().Retrieve(444, 1)
	if err != nil {
		t.Errorf("TranslationProviders.Retrieve returned error: %v", err)
	}

	want := TranslationProvider{
		ProviderID:   1,
		Name:         "Gengo",
		Slug:         "gengo",
		PricePairMin: "0.00",
		WebsiteURL:   "https://gengo.com",
		Description:  "At Gengo, our mission is to provide language services to everyone and connect a global community. Our network of over 18,000 translators are tested and qualified to meet stringent project standards. Translators are based around the world and in every timezone, which enables us to support over 35 languages and work towards the Gengo mission.",
		Tiers: []TranslationTier{
			{
				TierID: 1,
				Title:  "Native speaker",
			},
		},
		Pairs: []TranslationPair{
			{
				TierID:       1,
				FromLangISO:  "en_US",
				FromLangName: "English (United States)",
				ToLangISO:    "ru",
				ToLangName:   "Russian",
				PricePerWord: "0.07",
			},
		},
	}

	if !reflect.DeepEqual(r, want) {
		t.Errorf("TranslationProviders.Retrieve returned %+v, want %+v", r, want)
	}
}
