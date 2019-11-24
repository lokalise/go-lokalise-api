package lokalise

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestTranslationStatusService_Create(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(
		fmt.Sprintf("/projects/%s/custom_translation_statuses", testProjectID),
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			testMethod(t, r, "POST")
			testHeader(t, r, apiTokenHeader, testApiToken)
			data := `{
				"title": "Reviewed by doctors",
				"color": "#ff9f1a"
			}`

			req := new(bytes.Buffer)
			_ = json.Compact(req, []byte(data))

			testBody(t, r, req.String())

			_, _ = fmt.Fprint(w, `{
				"project_id": "`+testProjectID+`",
				"custom_translation_status": {
					"status_id": 124,
					"title": "Reviewed by doctors",
					"color": "#ff9f1a"
				 }
			}`)
		})

	r, err := client.TranslationStatuses().Create(testProjectID, NewTranslationStatus{
		Title: "Reviewed by doctors",
		Color: "#ff9f1a",
	})
	if err != nil {
		t.Errorf("TranslationStatuses.Create returned error: %v", err)
	}

	want := TranslationStatus{
		StatusID: 124,
		Title:    "Reviewed by doctors",
		Color:    "#ff9f1a",
	}

	if !reflect.DeepEqual(r.TranslationStatus, want) {
		t.Errorf("TranslationStatuses.Create returned %+v, want %+v", r.TranslationStatus, want)
	}
}

func TestTranslationStatusService_Delete(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(
		fmt.Sprintf("/projects/%s/custom_translation_statuses/%d", testProjectID, 421),
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			testMethod(t, r, "DELETE")
			testHeader(t, r, apiTokenHeader, testApiToken)

			_, _ = fmt.Fprint(w, `{
				"project_id": "`+testProjectID+`",
				"deleted": true
			}`)
		},
	)

	r, err := client.TranslationStatuses().Delete(testProjectID, 421)
	if err != nil {
		t.Errorf("TranslationStatuses.Delete returned error: %v", err)
	}

	want := DeleteTranslationStatusResponse{
		WithProjectID: WithProjectID{
			ProjectID: testProjectID,
		},
		Deleted: true,
	}

	if !reflect.DeepEqual(r, want) {
		t.Errorf("TranslationStatuses.Delete returned %+v, want %+v", r, want)
	}
}

func TestTranslationStatusService_List(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(
		fmt.Sprintf("/projects/%s/custom_translation_statuses", testProjectID),
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			testMethod(t, r, "GET")
			testHeader(t, r, apiTokenHeader, testApiToken)

			_, _ = fmt.Fprint(w, `{
				"project_id": "`+testProjectID+`",
				"custom_translation_statuses": [
					{
						"status_id": 123
					}
				]
			}`)
		})

	r, err := client.TranslationStatuses().List(testProjectID)
	if err != nil {
		t.Errorf("TranslationStatuses.List returned error: %v", err)
	}

	want := []TranslationStatus{
		{
			StatusID: 123,
		},
	}

	if !reflect.DeepEqual(r.TranslationStatuses, want) {
		t.Errorf("TranslationStatuses.List returned %+v, want %+v", r.TranslationStatuses, want)
	}
}

func TestTranslationStatusService_ListColors(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(
		fmt.Sprintf("/projects/%s/custom_translation_statuses/colors", testProjectID),
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			testMethod(t, r, "GET")
			testHeader(t, r, apiTokenHeader, testApiToken)

			_, _ = fmt.Fprint(w, `{
				"colors": [
					"#61bd4f"
				]
			}`)
		})

	r, err := client.TranslationStatuses().ListColors(testProjectID)
	if err != nil {
		t.Errorf("TranslationStatuses.ListColors returned error: %v", err)
	}

	want := []string{"#61bd4f"}

	if !reflect.DeepEqual(r.Colors, want) {
		t.Errorf("TranslationStatuses.ListColors returned %+v, want %+v", r.Colors, want)
	}
}

func TestTranslationStatusService_Retrieve(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(
		fmt.Sprintf("/projects/%s/custom_translation_statuses/%d", testProjectID, 124),
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			testMethod(t, r, "GET")
			testHeader(t, r, apiTokenHeader, testApiToken)

			_, _ = fmt.Fprint(w, `{
				"project_id": "`+testProjectID+`",
				"custom_translation_status": {
					"status_id": 124
				}
			}`)
		})

	r, err := client.TranslationStatuses().Retrieve(testProjectID, 124)
	if err != nil {
		t.Errorf("TranslationStatuses.Retrieve returned error: %v", err)
	}

	want := TranslationStatus{
		StatusID: 124,
	}

	if !reflect.DeepEqual(r.TranslationStatus, want) {
		t.Errorf("TranslationStatuses.Retrieve returned %+v, want %+v", r.TranslationStatus, want)
	}
}

func TestTranslationStatusService_Update(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(
		fmt.Sprintf("/projects/%s/custom_translation_statuses/%d", testProjectID, 403021),
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			testMethod(t, r, "PUT")
			testHeader(t, r, apiTokenHeader, testApiToken)
			data := `{
				"title": "Reviewed by staff"
			}`

			req := new(bytes.Buffer)
			_ = json.Compact(req, []byte(data))

			testBody(t, r, req.String())

			_, _ = fmt.Fprint(w, `{
				"project_id": "`+testProjectID+`",
				"custom_translation_status": {
					"status_id": 124
				}
			}`)
		})

	r, err := client.TranslationStatuses().Update(testProjectID, 403021, UpdateTranslationStatus{
		Title: "Reviewed by staff",
	})
	if err != nil {
		t.Errorf("TranslationStatuses.Update returned error: %v", err)
	}

	want := TranslationStatus{
		StatusID: 124,
	}

	if !reflect.DeepEqual(r.TranslationStatus, want) {
		t.Errorf("TranslationStatuses.Update returned %+v, want %+v", r.TranslationStatus, want)
	}
}
