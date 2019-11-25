package lokalise

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestLanguageService_Create(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(
		fmt.Sprintf("/projects/%s/languages", testProjectID),
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			testMethod(t, r, "POST")
			testHeader(t, r, apiTokenHeader, testApiToken)
			data := `{
				"languages": [
					{
						"lang_iso": "en"
					}
				]
			}`

			req := new(bytes.Buffer)
			_ = json.Compact(req, []byte(data))

			testBody(t, r, req.String())

			_, _ = fmt.Fprint(w, `{
				"project_id": "`+testProjectID+`",
				"languages": [
					{
						"lang_id": 640,
						"lang_iso": "en",
						"lang_name": "English",
						"is_rtl": false,
						"plural_forms": [
							"one", "other"
						]
					}
				]
			}`)
		})

	r, err := client.Languages().Create(testProjectID, []NewLanguage{
		{
			LangISO: "en",
		},
	})
	if err != nil {
		t.Errorf("Languages.Create returned error: %v", err)
	}

	want := []Language{
		{
			LangID:   640,
			LangISO:  "en",
			LangName: "English",
			IsRTL:    false,
			PluralForms: []string{
				"one",
				"other",
			},
		},
	}

	if !reflect.DeepEqual(r.Languages, want) {
		t.Errorf("Languages.Create returned %+v, want %+v", r.Languages, want)
	}
}

func TestLanguageService_Delete(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(
		fmt.Sprintf("/projects/%s/languages/%d", testProjectID, 640),
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			testMethod(t, r, "DELETE")
			testHeader(t, r, apiTokenHeader, testApiToken)

			_, _ = fmt.Fprint(w, `{
				"project_id": "`+testProjectID+`",
				"language_deleted": true
			}`)
		})

	r, err := client.Languages().Delete(testProjectID, 640)
	if err != nil {
		t.Errorf("Languages.Delete returned error: %v", err)
	}

	want := DeleteLanguageResponse{
		WithProjectID: WithProjectID{
			ProjectID: testProjectID,
		},
		LanguageDeleted: true,
	}

	if !reflect.DeepEqual(r, want) {
		t.Errorf("Languages.Delete returned %+v, want %+v", r, want)
	}
}

func TestLanguageService_ListProject(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(
		fmt.Sprintf("/projects/%s/languages", testProjectID),
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			testMethod(t, r, "GET")
			testHeader(t, r, apiTokenHeader, testApiToken)

			_, _ = fmt.Fprint(w, `{
			"languages": [
				{
					"lang_id": 640
				}
			]
		}`)
		})

	r, err := client.Languages().ListProject(testProjectID)
	if err != nil {
		t.Errorf("Languages.ListProject returned error: %v", err)
	}

	want := []Language{
		{
			LangID: 640,
		},
	}

	if !reflect.DeepEqual(r.Languages, want) {
		t.Errorf("Contributors.ListProject returned %+v, want %+v", r.Languages, want)
	}
}

func TestLanguageService_ListSystem(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(
		"/system/languages",
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			testMethod(t, r, "GET")
			testHeader(t, r, apiTokenHeader, testApiToken)

			_, _ = fmt.Fprint(w, `{
				"languages": [
					{
						"lang_id": 640
					}
				]
			}`)
		})

	r, err := client.Languages().ListSystem()
	if err != nil {
		t.Errorf("Languages.ListSystem returned error: %v", err)
	}

	want := []Language{
		{
			LangID: 640,
		},
	}

	if !reflect.DeepEqual(r.Languages, want) {
		t.Errorf("Contributors.ListSystem returned %+v, want %+v", r.Languages, want)
	}
}

func TestLanguageService_Retrieve(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(
		fmt.Sprintf("/projects/%s/languages/%d", testProjectID, 640),
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			testMethod(t, r, "GET")
			testHeader(t, r, apiTokenHeader, testApiToken)

			_, _ = fmt.Fprint(w, `{
				"project_id": "`+testProjectID+`",
				"language": {
					"lang_id": 640
				}   
			}`)
		})

	r, err := client.Languages().Retrieve(testProjectID, 640)
	if err != nil {
		t.Errorf("Languages.Retrieve returned error: %v", err)
	}

	want := Language{
		LangID: 640,
	}

	if !reflect.DeepEqual(r.Language, want) {
		t.Errorf("Languages.Retrieve returned %+v, want %+v", r.Language, want)
	}
}

func TestLanguageService_Update(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(
		fmt.Sprintf("/projects/%s/languages/%d", testProjectID, 640),
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			testMethod(t, r, "PUT")
			testHeader(t, r, apiTokenHeader, testApiToken)
			data := `{
				"lang_iso": "en-US",
				"plural_forms": [
					"one", "zero", "few", "other"
				]
			}`

			req := new(bytes.Buffer)
			_ = json.Compact(req, []byte(data))

			testBody(t, r, req.String())

			_, _ = fmt.Fprint(w, `{
				"project_id": "`+testProjectID+`",
				"language": {
					"lang_id": 640
				}   
			}`)
		})

	r, err := client.Languages().Update(testProjectID, 640, UpdateLanguage{
		LangISO: "en-US",
		PluralForms: []string{
			"one", "zero", "few", "other",
		},
	})
	if err != nil {
		t.Errorf("Languages.Update returned error: %v", err)
	}

	want := Language{
		LangID: 640,
	}

	if !reflect.DeepEqual(r.Language, want) {
		t.Errorf("Languages.Update returned %+v, want %+v", r.Language, want)
	}
}
