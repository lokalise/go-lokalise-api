package lokalise

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestOrderService_Create(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(
		fmt.Sprintf("/teams/%d/orders", 1),
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			testMethod(t, r, "POST")
			testHeader(t, r, apiTokenHeader, testApiToken)
			data := `{
				"project_id": "` + testProjectID + `",
				"card_id": 12345,
				"briefing": "Terms of use of our app.",
				"source_language_iso": "en_US",
				"target_language_isos": [
					"ru",
					"fr",
					"it"
				],
				"keys": [
					12213,
					12214,
					12215,
					21216
				],
				"provider_slug": "gengo",
				"translation_tier": 1
			}`

			req := new(bytes.Buffer)
			_ = json.Compact(req, []byte(data))

			testBody(t, r, req.String())

			_, _ = fmt.Fprint(w, `{
				"order_id": "20181231AAAA",
				"project_id": "`+testProjectID+`",
				"card_id": 12345,
				"status": "in progress",
				"created_at": "2018-12-31 12:00:00 (Etc/UTC)",
				"created_at_timestamp": 1546257600,
				"created_by": 420,
				"created_by_email": "jonn@company.com",
				"source_language_iso": "en_US",
				"target_language_isos": [
					"ru",
					"fr",
					"it"
				],
				"keys": [
					12213,
					12214,
					12215,
					21216
				],
				"source_words": {
					"ru": 256,
					"fr": 222,
					"it": 256
				},
				"provider_slug": "gengo",
				"translation_style": "friendly",
				"translation_tier": 1,
				"translation_tier_name": "Native speakers",
				"briefing": "Terms of use of our app.",
				"total": 177.90,
				"dry_run": false
			}`)
		})

	r, err := client.Orders().Create(1, CreateOrder{
		ProjectID:         testProjectID,
		CardID:            12345,
		Briefing:          "Terms of use of our app.",
		SourceLangISO:     "en_US",
		TargetLangISOs:    []string{"ru", "fr", "it"},
		Keys:              []int{12213, 12214, 12215, 21216},
		ProviderSlug:      "gengo",
		TranslationTierID: 1,
	})
	if err != nil {
		t.Errorf("Orders.Create returned error: %v", err)
	}

	want := Order{
		WithProjectID: WithProjectID{
			ProjectID: testProjectID,
		},
		WithCreationTime: WithCreationTime{
			CreatedAt:   "2018-12-31 12:00:00 (Etc/UTC)",
			CreatedAtTs: 1546257600,
		},
		WithCreationUser: WithCreationUser{
			CreatedBy:      420,
			CreatedByEmail: "jonn@company.com",
		},
		OrderID:       "20181231AAAA",
		CardID:        12345,
		Status:        "in progress",
		SourceLangISO: "en_US",
		TargetLangISOs: []string{
			"ru",
			"fr",
			"it",
		},
		Keys: []int64{
			12213,
			12214,
			12215,
			21216,
		},
		SourceWords: map[string]int64{
			"ru": 256,
			"fr": 222,
			"it": 256,
		},
		ProviderSlug:        "gengo",
		TranslationStyle:    "friendly",
		TranslationTierID:   1,
		TranslationTierName: "Native speakers",
		Briefing:            "Terms of use of our app.",
		Total:               177.90,
	}

	if !reflect.DeepEqual(r, want) {
		t.Errorf("Orders.Create returned %+v, want %+v", r, want)
	}
}

func TestOrderService_List(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(
		fmt.Sprintf("/teams/%d/orders", 1),
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			testMethod(t, r, "GET")
			testHeader(t, r, apiTokenHeader, testApiToken)

			_, _ = fmt.Fprint(w, `{
			"orders": [
					{
						"order_id": "20181231AAAA"
					}
				]
			}`)
		})

	r, err := client.Orders().List(1)
	if err != nil {
		t.Errorf("Orders.List returned error: %v", err)
	}

	want := []Order{
		{
			OrderID: "20181231AAAA",
		},
	}

	if !reflect.DeepEqual(r.Orders, want) {
		t.Errorf("Orders.List returned %+v, want %+v", r.Orders, want)
	}
}

func TestOrderService_Retrieve(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(
		fmt.Sprintf("/teams/%d/orders/%s", 1, "20181231AAAA"),
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			testMethod(t, r, "GET")
			testHeader(t, r, apiTokenHeader, testApiToken)

			_, _ = fmt.Fprint(w, `{
				"order_id": "20181231AAAA"
			}`)
		})

	r, err := client.Orders().Retrieve(1, "20181231AAAA")
	if err != nil {
		t.Errorf("Orders.Retrieve returned error: %v", err)
	}

	want := Order{
		OrderID: "20181231AAAA",
	}

	if !reflect.DeepEqual(r, want) {
		t.Errorf("Orders.Retrieve returned %+v, want %+v", r, want)
	}
}
