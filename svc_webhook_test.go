package lokalise

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestWebhookService_Create(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(
		fmt.Sprintf("/projects/%s/webhooks", testProjectID),
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			testMethod(t, r, "POST")
			testHeader(t, r, apiTokenHeader, testApiToken)
			data := `{
				"url": "https://my.domain.com/webhook",
				"events": [
					"project.translation.updated"
				],
				"event_lang_map": [
					{
						"event": "project.translation.updated",
						"lang_iso_codes": [
							"en_GB"
						]
					}
				]
			}`

			req := new(bytes.Buffer)
			_ = json.Compact(req, []byte(data))

			testBody(t, r, req.String())

			_, _ = fmt.Fprint(w, `{
				"project_id": "`+testProjectID+`",
				"webhook": {
					"webhook_id": "138c1ffa0ad94848f01f980e7f2f2af19d1bd553",
					"url": "https://my.domain.com/webhook",
					"secret": "ccbe0d72a8377b2c91db25d6281192eade03158c",
					"events": [
						"project.translation.updated"
					],
					"event_lang_map": [
						{
							"event": "project.translation.updated",
							"lang_iso_codes": [
								"en_GB"
							]
						}
					]
				}
			}`)
		})

	r, err := client.Webhooks().Create(testProjectID, CreateWebhook{
		URL: "https://my.domain.com/webhook",
		Events: []string{
			"project.translation.updated",
		},
		EventLangMap: []EventLang{
			{
				Event: "project.translation.updated",
				LangISOCodes: []string{
					"en_GB",
				},
			},
		},
	})
	if err != nil {
		t.Errorf("Webhooks.Create returned error: %v", err)
	}

	want := Webhook{
		WebhookID: "138c1ffa0ad94848f01f980e7f2f2af19d1bd553",
		URL:       "https://my.domain.com/webhook",
		Secret:    "ccbe0d72a8377b2c91db25d6281192eade03158c",
		Events: []string{
			"project.translation.updated",
		},
		EventLangMap: []EventLang{
			{
				Event: "project.translation.updated",
				LangISOCodes: []string{
					"en_GB",
				},
			},
		},
	}

	if !reflect.DeepEqual(r.Webhook, want) {
		t.Errorf("Webhooks.Create returned %+v, want %+v", r.Webhook, want)
	}
}

func TestWebhookService_Delete(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(
		fmt.Sprintf("/projects/%s/webhooks/%s", testProjectID, "138c1ffa0ad94848f01f980e7f2f2af19d1bd553"),
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			testMethod(t, r, "DELETE")
			testHeader(t, r, apiTokenHeader, testApiToken)

			_, _ = fmt.Fprint(w, `{
				"project_id": "`+testProjectID+`",
				"webhook_deleted": true
			}`)
		},
	)

	r, err := client.Webhooks().Delete(testProjectID, "138c1ffa0ad94848f01f980e7f2f2af19d1bd553")
	if err != nil {
		t.Errorf("Webhooks.Delete returned error: %v", err)
	}

	want := DeleteWebhookResponse{
		WithProjectID: WithProjectID{
			ProjectID: testProjectID,
		},
		Deleted: true,
	}

	if !reflect.DeepEqual(r, want) {
		t.Errorf("Webhooks.Delete returned %+v, want %+v", r, want)
	}
}

func TestWebhookService_List(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(
		fmt.Sprintf("/projects/%s/webhooks", testProjectID),
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			testMethod(t, r, "GET")
			testHeader(t, r, apiTokenHeader, testApiToken)

			_, _ = fmt.Fprint(w, `{
				"project_id": "`+testProjectID+`",
				"webhooks": [
					{
						"webhook_id": "138c1ffa0ad94848f01f980e7f2f2af19d1bd553"
					}
				]
			}`)
		})

	r, err := client.Webhooks().List(testProjectID)
	if err != nil {
		t.Errorf("Webhooks.List returned error: %v", err)
	}

	want := []Webhook{
		{
			WebhookID: "138c1ffa0ad94848f01f980e7f2f2af19d1bd553",
		},
	}

	if !reflect.DeepEqual(r.Webhooks, want) {
		t.Errorf("Webhooks.List returned %+v, want %+v", r.Webhooks, want)
	}
}

func TestWebhookService_Retrieve(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(
		fmt.Sprintf("/projects/%s/webhooks/%s", testProjectID, "138c1ffa0ad94848f01f980e7f2f2af19d1bd553"),
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			testMethod(t, r, "GET")
			testHeader(t, r, apiTokenHeader, testApiToken)

			_, _ = fmt.Fprint(w, `{
				"project_id": "`+testProjectID+`",
				"webhook": {
					"webhook_id": "138c1ffa0ad94848f01f980e7f2f2af19d1bd553"
				}
			}`)
		})

	r, err := client.Webhooks().Retrieve(testProjectID, "138c1ffa0ad94848f01f980e7f2f2af19d1bd553")
	if err != nil {
		t.Errorf("Webhooks.Retrieve returned error: %v", err)
	}

	want := Webhook{
		WebhookID: "138c1ffa0ad94848f01f980e7f2f2af19d1bd553",
	}

	if !reflect.DeepEqual(r.Webhook, want) {
		t.Errorf("Webhooks.Retrieve returned %+v, want %+v", r.Webhook, want)
	}
}

func TestWebhookService_Update(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(
		fmt.Sprintf("/projects/%s/webhooks/%s", testProjectID, "138c1ffa0ad94848f01f980e7f2f2af19d1bd553"),
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			testMethod(t, r, "PUT")
			testHeader(t, r, apiTokenHeader, testApiToken)
			data := `{
				"events": [
					"project.translation.proofread"
				],
				"event_lang_map": [
					{
						"event": "project.translation.proofread",
						"lang_iso_codes": [
							"en_GB"
						]
					}
				]
			}`

			req := new(bytes.Buffer)
			_ = json.Compact(req, []byte(data))

			testBody(t, r, req.String())

			_, _ = fmt.Fprint(w, `{
				"project_id": "`+testProjectID+`",
				"webhook": {
					"webhook_id": "138c1ffa0ad94848f01f980e7f2f2af19d1bd553"
				}
			}`)
		})

	r, err := client.Webhooks().Update(
		testProjectID,
		"138c1ffa0ad94848f01f980e7f2f2af19d1bd553", UpdateWebhook{
			Events: []string{"project.translation.proofread"},
			EventLangMap: []EventLang{
				{
					Event: "project.translation.proofread",
					LangISOCodes: []string{
						"en_GB",
					},
				},
			},
		},
	)
	if err != nil {
		t.Errorf("Webhooks.Update returned error: %v", err)
	}

	want := Webhook{
		WebhookID: "138c1ffa0ad94848f01f980e7f2f2af19d1bd553",
	}

	if !reflect.DeepEqual(r.Webhook, want) {
		t.Errorf("Webhooks.Update returned %+v, want %+v", r.Webhook, want)
	}
}
