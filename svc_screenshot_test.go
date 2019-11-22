package lokalise

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestScreenshotService_Create(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(
		fmt.Sprintf("/projects/%s/screenshots", testProjectID),
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			testMethod(t, r, "POST")
			testHeader(t, r, apiTokenHeader, testApiToken)
			data := `{
				"screenshots": [
					{
						"data": "data:image/jpeg;base64,D94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGL.....",
						"ocr": false,
						"key_ids": [
							1132290, 1132292, 1132293
						],
						"tags": [
							"onboarding"
						]
					}
				]
			}`

			req := new(bytes.Buffer)
			_ = json.Compact(req, []byte(data))

			testBody(t, r, req.String())

			_, _ = fmt.Fprint(w, `{
				"project_id": "`+testProjectID+`",
				"screenshots": [
					{
						"screenshot_id": 403021,
						"url": "https://s3-eu-west-1.amazonaws.com/files/screenshots/welcome.jpg",
						"key_ids": [
							1132290, 1132292, 1132293
						],
						"title": "Welcome screen",
						"description" : "",
						"tags": [
							"onboarding"
						],
						"width": 1024,
						"height": 768,
						"created_at": "2019-02-01 12:00:00 (Etc/UTC)",
						"created_at_timestamp": 1546257600
					}
				]
			}`)
		})

	r, err := client.Screenshots().Create(testProjectID, []NewScreenshot{
		{
			Body:   "data:image/jpeg;base64,D94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGL.....",
			Ocr:    Bool(false),
			KeyIDs: []int64{1132290, 1132292, 1132293},
			Tags:   []string{"onboarding"},
		},
	})
	if err != nil {
		t.Errorf("Screenshots.Create returned error: %v", err)
	}

	want := []Screenshot{
		{
			WithCreationTime: WithCreationTime{
				CreatedAt:   "2019-02-01 12:00:00 (Etc/UTC)",
				CreatedAtTs: 1546257600,
			},
			ScreenshotID:   403021,
			KeyIDs:         []int64{1132290, 1132292, 1132293},
			URL:            "https://s3-eu-west-1.amazonaws.com/files/screenshots/welcome.jpg",
			Title:          "Welcome screen",
			Description:    "",
			ScreenshotTags: []string{"onboarding"},
			Width:          1024,
			Height:         768,
		},
	}

	if !reflect.DeepEqual(r.Screenshots, want) {
		t.Errorf("Screenshots.Create returned %+v, want %+v", r.Screenshots, want)
	}
}

func TestScreenshotService_Delete(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(
		fmt.Sprintf("/projects/%s/screenshots/%d", testProjectID, 421),
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			testMethod(t, r, "DELETE")
			testHeader(t, r, apiTokenHeader, testApiToken)

			_, _ = fmt.Fprint(w, `{
				"project_id": "`+testProjectID+`",
				"screenshot_deleted": true
			}`)
		},
	)

	r, err := client.Screenshots().Delete(testProjectID, 421)
	if err != nil {
		t.Errorf("Screenshots.Delete returned error: %v", err)
	}

	want := DeleteScreenshotResponse{
		WithProjectID: WithProjectID{
			ProjectID: testProjectID,
		},
		Deleted: true,
	}

	if !reflect.DeepEqual(r, want) {
		t.Errorf("Screenshots.Delete returned %+v, want %+v", r, want)
	}
}

func TestScreenshotService_List(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(
		fmt.Sprintf("/projects/%s/screenshots", testProjectID),
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			testMethod(t, r, "GET")
			testHeader(t, r, apiTokenHeader, testApiToken)

			_, _ = fmt.Fprint(w, `{
			"project_id": "`+testProjectID+`",
			"screenshots": [
					{
						"screenshot_id": 403021
					}
				]
			}`)
		})

	r, err := client.Screenshots().List(testProjectID)
	if err != nil {
		t.Errorf("Screenshots.List returned error: %v", err)
	}

	want := []Screenshot{
		{
			ScreenshotID: 403021,
		},
	}

	if !reflect.DeepEqual(r.Screenshots, want) {
		t.Errorf("Screenshots.List returned %+v, want %+v", r.Screenshots, want)
	}
}

func TestScreenshotService_Retrieve(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(
		fmt.Sprintf("/projects/%s/screenshots/%d", testProjectID, 403021),
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			testMethod(t, r, "GET")
			testHeader(t, r, apiTokenHeader, testApiToken)

			_, _ = fmt.Fprint(w, `{
				"project_id": "`+testProjectID+`",
				"screenshot": {
					"screenshot_id": 403021
				}
			}`)
		})

	r, err := client.Screenshots().Retrieve(testProjectID, 403021)
	if err != nil {
		t.Errorf("Screenshots.Retrieve returned error: %v", err)
	}

	want := Screenshot{
		ScreenshotID: 403021,
	}

	if !reflect.DeepEqual(r.Screenshot, want) {
		t.Errorf("Screenshots.Retrieve returned %+v, want %+v", r.Screenshot, want)
	}
}

func TestScreenshotService_Update(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(
		fmt.Sprintf("/projects/%s/screenshots/%d", testProjectID, 403021),
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			testMethod(t, r, "PUT")
			testHeader(t, r, apiTokenHeader, testApiToken)
			data := `{
				"key_ids": [
					1132290, 1132292
				],
				"tags": [
					"main"
				]
			}`

			req := new(bytes.Buffer)
			_ = json.Compact(req, []byte(data))

			testBody(t, r, req.String())

			_, _ = fmt.Fprint(w, `{
				"project_id": "`+testProjectID+`",
				"screenshot": {
					"screenshot_id": 403021
				}
			}`)
		})

	r, err := client.Screenshots().Update(testProjectID, 403021, UpdateScreenshot{
		KeyIDs: []int64{1132290, 1132292},
		Tags:   []string{"main"},
	})
	if err != nil {
		t.Errorf("Screenshots.Update returned error: %v", err)
	}

	want := Screenshot{
		ScreenshotID: 403021,
	}

	if !reflect.DeepEqual(r.Screenshot, want) {
		t.Errorf("Screenshots.Update returned %+v, want %+v", r.Screenshot, want)
	}
}
