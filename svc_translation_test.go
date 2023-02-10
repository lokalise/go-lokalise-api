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
		t.Errorf(assertionTemplate, "Translations.List", r.Translations, want)
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
					"is_unverified": true,
					"is_reviewed": false,
					"reviewed_by": 0,
					"words": 2,
					"custom_translation_statuses": [],
					"task_id": 123	
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
		IsUnverified:              true,
		IsReviewed:                false,
		ReviewedBy:                0,
		Words:                     2,
		CustomTranslationStatuses: []TranslationStatus{},
		TaskID:                    123,
	}

	if !reflect.DeepEqual(r.Translation, want) {
		t.Errorf(assertionTemplate, "Translations.Retrieve", r.Translation, want)
	}
}

func TestTranslationService_Retrieve_NullTaskID(t *testing.T) {
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
					"is_unverified": true,
					"is_reviewed": false,
					"reviewed_by": 0,
					"words": 2,
					"custom_translation_statuses": [],
					"task_id": null	
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
		IsUnverified:              true,
		IsReviewed:                false,
		ReviewedBy:                0,
		Words:                     2,
		CustomTranslationStatuses: []TranslationStatus{},
		TaskID:                    0,
	}

	if !reflect.DeepEqual(r.Translation, want) {
		t.Errorf(assertionTemplate, "Translations.Retrieve", r.Translation, want)
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
				"is_unverified": false,
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
		Translation:  "Quick brown fox jumps over the lazy dog.",
		IsUnverified: Bool(false),
		IsReviewed:   true,
	})
	if err != nil {
		t.Errorf("Translations.Update returned error: %v", err)
	}

	want := Translation{
		TranslationID: 640,
	}

	if !reflect.DeepEqual(r.Translation, want) {
		t.Errorf(assertionTemplate, "Translations.Update", r.Translation, want)
	}
}

func TestNewTranslation_MarshalJSON(t *testing.T) {
	translations := []NewTranslation{
		{
			LanguageISO: "en",
			Translation: "simple text",
		},
		{
			LanguageISO: "en",
			Translation: "{\"one\":\"oneText\",\"other\":\"otherText\"}",
		},
	}

	want := JsonCompact(`
	[
		{
		   "language_iso": "en",
		   "translation": "simple text"
		},
		{
		   "language_iso": "en",
		   "translation": {
			  "one": "oneText",
			  "other": "otherText"
		   }
		}
 	]`)

	marshal, err := json.Marshal(translations)
	if err != nil {
		t.Errorf("NewTranslation marshalling returned error %s", err)
	}

	if string(marshal) != want {
		t.Errorf("NewTranslation marshalling mismatch: %+v, want %+v", string(marshal), want)
	}
}

func TestUpdateTranslation_MarshalJSON(t *testing.T) {
	translations := []UpdateTranslation{
		{
			Translation: "simple text",
		},
		{
			Translation: "{\"one\":\"oneText\",\"other\":\"otherText\"}",
		},
	}

	want := JsonCompact(`
	[
		{
		   "translation": "simple text"
		},
		{
		   "translation": {
			  "one": "oneText",
			  "other": "otherText"
		   }
		}
 	]`)

	marshal, err := json.Marshal(translations)
	if err != nil {
		t.Errorf("UpdateTranslation marshalling returned error %s", err)
	}

	if string(marshal) != want {
		t.Errorf("UpdateTranslation marshalling mismatch: %+v, want %+v", string(marshal), want)
	}
}

func TestNewTranslation_UnmarshalJSON(t *testing.T) {
	raw := `
	[
		{
		   "language_iso": "en",
		   "translation": "simple text"
		},
		{
		   "language_iso": "en",
		   "translation": {
			  "one": "oneText",
			  "other": "otherText"
		   }
		}
 	]`

	want := []NewTranslation{
		{
			LanguageISO: "en",
			Translation: "simple text",
		},
		{
			LanguageISO: "en",
			Translation: "{\"one\":\"oneText\",\"other\":\"otherText\"}",
		},
	}

	var result []NewTranslation

	err := json.Unmarshal([]byte(raw), &result)
	if err != nil {
		t.Errorf("NewTranslation unmarshalling returned error %s", err)
	}

	if !reflect.DeepEqual(want, result) {
		t.Errorf("NewTranslation unmarshalling mismatch: %+v, want %+v", result, want)
	}
}

func TestUpdateTranslation_UnmarshalJSON(t *testing.T) {
	raw := `
	[
		{
		   "translation": "simple text"
		},
		{
		   "translation": {
			  "one": "oneText",
			  "other": "otherText"
		   }
		}
 	]`

	want := []UpdateTranslation{
		{
			Translation: "simple text",
		},
		{
			Translation: "{\"one\":\"oneText\",\"other\":\"otherText\"}",
		},
	}

	var result []UpdateTranslation

	err := json.Unmarshal([]byte(raw), &result)
	if err != nil {
		t.Errorf("UpdateTranslation unmarshalling returned error %s", err)
	}

	if !reflect.DeepEqual(want, result) {
		t.Errorf("UpdateTranslation unmarshalling mismatch: %+v, want %+v", result, want)
	}
}

func JsonCompact(text string) string {
	compactedBuffer := new(bytes.Buffer)
	err := json.Compact(compactedBuffer, []byte(text))
	if err != nil {
		panic(fmt.Sprintf("Invalid test data definition %+v", err))
	}
	return compactedBuffer.String()
}
