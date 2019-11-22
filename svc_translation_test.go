package lokalise

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestTranslationService_List(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(
		fmt.Sprintf("/projects/%s/translations", testProjectID),
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			testMethod(t, r, "GET")
			testHeader(t, r, apiTokenHeader, testApiToken)

			_, _ = fmt.Fprint(w, `{
				"project_id": "`+testProjectID+`",
				"translations": [
					{
						"translation_id": 344412
					}                    
				]
			}`)
		})

	r, err := client.Translations().List(testProjectID)
	if err != nil {
		t.Errorf("Translations.List returned error: %v", err)
	}

	want := []Translation{
		{
			TranslationID: 344412,
		},
	}

	if !reflect.DeepEqual(r.Translations, want) {
		t.Errorf("Translations.List returned %+v, want %+v", r.Translations, want)
	}
}

func TestTranslationService_Retrieve(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(
		fmt.Sprintf("/projects/%s/translations/%d", testProjectID, 640),
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			testMethod(t, r, "GET")
			testHeader(t, r, apiTokenHeader, testApiToken)

			_, _ = fmt.Fprint(w, `{
				"project_id": "`+testProjectID+`",
				"translation": {
					"translation_id": 344412,
					"key_id": 553662,
					"language_iso": "en_US",
					"modified_at": "2018-12-31 12:00:00 (Etc/UTC)",
					"modified_at_timestamp": 1546257600,
					"modified_by": 420,
					"modified_by_email": "user@mycompany.com",
					"translation": "Hello, world!",
					"is_fuzzy": true,
					"is_reviewed": false,
					"reviewed_by": 0,
					"words": 2,
					"custom_translation_statuses": []
				}
			}`)
		})

	r, err := client.Translations().Retrieve(testProjectID, 640)
	if err != nil {
		t.Errorf("Keys.Retrieve returned error: %v", err)
	}

	want := Translation{
		TranslationID:             344412,
		KeyID:                     553662,
		LanguageISO:               "en_US",
		ModifiedAt:                "2018-12-31 12:00:00 (Etc/UTC)",
		ModifiedAtTs:              1546257600,
		ModifiedBy:                420,
		ModifiedByEmail:           "user@mycompany.com",
		Translation:               "Hello, world!",
		IsFuzzy:                   true,
		IsReviewed:                false,
		ReviewedBy:                0,
		Words:                     2,
		CustomTranslationStatuses: []TranslationStatus{},
	}

	if !reflect.DeepEqual(r.Translation, want) {
		t.Errorf("Keys.Retrieve returned %+v, want %+v", r.Translation, want)
	}
}

func TestTranslationService_Update(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(
		fmt.Sprintf("/projects/%s/translations/%d", testProjectID, 640),
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			testMethod(t, r, "PUT")
			testHeader(t, r, apiTokenHeader, testApiToken)
			data := `{
				"translation": "Quick brown fox jumps over the lazy dog.",
				"is_fuzzy": false,
				"is_reviewed": true
			}`

			req := new(bytes.Buffer)
			_ = json.Compact(req, []byte(data))

			testBody(t, r, req.String())

			_, _ = fmt.Fprint(w, `{
				"project_id": "`+testProjectID+`",
				"translation": {
					"translation_id": 640
				}   
			}`)
		})

	r, err := client.Translations().Update(testProjectID, 640, UpdateTranslation{
		Translation: "Quick brown fox jumps over the lazy dog.",
		IsFuzzy:     Bool(false),
		IsReviewed:  true,
	})
	if err != nil {
		t.Errorf("Translations.Update returned error: %v", err)
	}

	want := Translation{
		TranslationID: 640,
	}

	if !reflect.DeepEqual(r.Translation, want) {
		t.Errorf("Translations.Update returned %+v, want %+v", r.Translation, want)
	}
}
