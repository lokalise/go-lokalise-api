package lokalise

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestSegmentationService_List(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	var keyID int64 = 640
	languageIso := "en_US"

	mux.HandleFunc(
		fmt.Sprintf("/projects/%s/keys/%d/segments/%s", testProjectID, keyID, languageIso),
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			testMethod(t, r, "GET")
			testHeader(t, r, apiTokenHeader, testApiToken)

			_, _ = fmt.Fprint(w, `{
                "project_id": "`+testProjectID+`",
				"segments": [
					{
						"segment_number": 1,
                        "language_iso": "en_US"
					}
				]
			}`)
		})

	r, err := client.Segments().List(testProjectID, keyID, languageIso)
	if err != nil {
		t.Errorf("Segments.List returned error: %v", err)
	}

	want := []Segment{
		{
			SegmentNumber: 1,
			LanguageIso:   languageIso,
		},
	}

	if !reflect.DeepEqual(r.Segments, want) {
		t.Errorf(assertionTemplate, "Segments.List", r.Segments, want)
	}
}

func TestSegmentationService_Retrieve(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	var keyID int64 = 640
	languageIso := "en_US"
	var segmentNumber int64 = 1

	mux.HandleFunc(
		fmt.Sprintf("/projects/%s/keys/%d/segments/%s/%d", testProjectID, keyID, languageIso, segmentNumber),
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			testMethod(t, r, "GET")
			testHeader(t, r, apiTokenHeader, testApiToken)

			_, _ = fmt.Fprint(w, `{
                "project_id": "`+testProjectID+`",
                "key_id": 640, 
                "language_iso": "en_US",
				"segment": {
                	"segment_number": 1,
                	"language_iso": "en_US"
                }
			}`)
		})

	r, err := client.Segments().Retrieve(testProjectID, keyID, languageIso, segmentNumber)
	if err != nil {
		t.Errorf("Segments.Retrieve returned error: %v", err)
	}

	want := Segment{
		SegmentNumber: segmentNumber,
		LanguageIso:   languageIso,
	}

	if !reflect.DeepEqual(r.Segment, want) {
		t.Errorf(assertionTemplate, "Segments.Retrieve", r.Segment, want)
	}
}

func TestSegmentationService_Update(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	var keyID int64 = 640
	languageIso := "en_US"
	var segmentNumber int64 = 1

	mux.HandleFunc(
		fmt.Sprintf("/projects/%s/keys/%d/segments/%s/%d", testProjectID, keyID, languageIso, segmentNumber),
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			testMethod(t, r, "PUT")
			testHeader(t, r, apiTokenHeader, testApiToken)
			data := `{
				"value": "Quick brown fox jumps over the lazy dog.",
				"is_fuzzy": false,
				"is_reviewed": true,
				"custom_translation_status_ids": [
					2, 3
				]
			}`

			req := new(bytes.Buffer)
			_ = json.Compact(req, []byte(data))

			testBody(t, r, req.String())

			_, _ = fmt.Fprint(w, `{
                "project_id": "`+testProjectID+`",
                "key_id": 640, 
                "language_iso": "en_US",
				"segment": {
                	"segment_number": 1,
                	"language_iso": "en_US"
                }
			}`)
		})

	r, err := client.Segments().Update(testProjectID, 640, "en_US", 1, SegmentUpdateRequest{
		Value:                      "Quick brown fox jumps over the lazy dog.",
		IsFuzzy:                    Bool(false),
		IsReviewed:                 Bool(true),
		CustomTranslationStatusIds: []int64{2, 3},
	})
	if err != nil {
		t.Errorf("Segments.Update returned error: %v", err)
	}

	want := Segment{
		SegmentNumber: segmentNumber,
		LanguageIso:   languageIso,
	}

	if !reflect.DeepEqual(r.Segment, want) {
		t.Errorf(assertionTemplate, "Segments.Update", r.Segment, want)
	}
}
