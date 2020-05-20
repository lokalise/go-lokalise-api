package lokalise

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestFileService_Download(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(
		fmt.Sprintf("/projects/%s/files/download", testProjectID),
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			testMethod(t, r, "POST")
			testHeader(t, r, apiTokenHeader, testApiToken)
			data := `{
				"format": "json",
				"original_filenames": true
			}`

			req := new(bytes.Buffer)
			_ = json.Compact(req, []byte(data))

			testBody(t, r, req.String())

			_, _ = fmt.Fprint(w, `{
				"project_id": "`+testProjectID+`",
				"bundle_url": "https://s3-eu-west-1.amazonaws.com/lokalise-assets/export/MyApp-locale.zip"
			}`)
		})

	r, err := client.Files().Download(testProjectID, FileDownload{
		Format:            "json",
		OriginalFilenames: Bool(true),
	})
	if err != nil {
		t.Errorf("Files.Download returned error: %v", err)
	}

	want := FileDownloadResponse{
		WithProjectID: WithProjectID{ProjectID: testProjectID},
		BundleURL:     "https://s3-eu-west-1.amazonaws.com/lokalise-assets/export/MyApp-locale.zip",
	}

	if !reflect.DeepEqual(r, want) {
		t.Errorf("Files.Download returned %+v, want %+v", r, want)
	}
}

func TestFileService_List(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(
		fmt.Sprintf("/projects/%s/files", testProjectID),
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			testMethod(t, r, "GET")
			testHeader(t, r, apiTokenHeader, testApiToken)

			_, _ = fmt.Fprint(w, `{
				"project_id": "`+testProjectID+`",
				"files": [
					{
						"filename": "index.json",
						"key_count": 441
					}
				]
			}`)
		})

	r, err := client.Files().List(testProjectID)
	if err != nil {
		t.Errorf("Files.List returned error: %v", err)
	}

	want := []File{
		{
			Filename: "index.json",
			KeyCount: 441,
		},
	}

	if !reflect.DeepEqual(r.Files, want) {
		t.Errorf("Files.List returned %+v, want %+v", r.Files, want)
	}
}

func TestFileService_Upload(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(
		fmt.Sprintf("/projects/%s/files/upload", testProjectID),
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			testMethod(t, r, "POST")
			testHeader(t, r, apiTokenHeader, testApiToken)
			data := `{
				"data": "D94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGL.....",
				"filename": "index.json",
				"lang_iso": "en",
				"tags": [
					"index", "admin", "v2.0"
				],
				"convert_placeholders": true,
				"custom_translation_status_ids": [1, 2, 3],
				"custom_translation_status_inserted_keys": false,
				"custom_translation_status_updated_keys": true,
				"custom_translation_status_skipped_keys": true,
				"queue": true
			}`

			req := new(bytes.Buffer)
			_ = json.Compact(req, []byte(data))

			testBody(t, r, req.String())

			_, _ = fmt.Fprint(w, `{
				"project_id": "`+testProjectID+`",
				"process": {
					"process_id": "2e0559e60e856555fbc15bdf78ab2b0ca3406e8f",
					"type": "file-import",
					"status": "queued",
					"message": "",
					"created_by": 1234,
					"created_by_email": "example@example.com",
					"created_at": "2020-04-20 13:43:43 (Etc/UTC)",
					"created_at_timestamp": 1587390223
				}
			}`)
		})

	r, err := client.Files().Upload(testProjectID, FileUpload{
		Data:                                "D94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGL.....",
		Filename:                            "index.json",
		LangISO:                             "en",
		Tags:                                []string{"index", "admin", "v2.0"},
		ConvertPlaceholders:                 Bool(true),
		CustomTranslationStatusIds:          []int64{1, 2, 3},
		CustomTranslationStatusInsertedKeys: Bool(false),
		CustomTranslationStatusUpdatedKeys:  Bool(true),
		CustomTranslationStatusSkippedKeys:  Bool(true),
	})
	if err != nil {
		t.Errorf("Files.Upload returned error: %v", err)
	}

	want := FileUploadResponse{
		WithProjectID: WithProjectID{
			ProjectID: testProjectID,
		},
		Process: QueuedProcess{
			ID:      "2e0559e60e856555fbc15bdf78ab2b0ca3406e8f",
			Type:    "file-import",
			Status:  "queued",
			Message: "",
			WithCreationUser: WithCreationUser{
				CreatedBy:      1234,
				CreatedByEmail: "example@example.com",
			},
			WithCreationTime: WithCreationTime{
				CreatedAt:   "2020-04-20 13:43:43 (Etc/UTC)",
				CreatedAtTs: 1587390223,
			},
		},
	}

	if !reflect.DeepEqual(r, want) {
		t.Errorf("Files.Upload returned %+v, want %+v", r, want)
	}
}

func TestFileServiceDefaults_Upload(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(
		fmt.Sprintf("/projects/%s/files/upload", testProjectID),
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			testMethod(t, r, "POST")
			testHeader(t, r, apiTokenHeader, testApiToken)
			data := `{
				"data": "D94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGL.....",
				"filename": "index.json",
				"lang_iso": "en",
				"tags": [
					"index", "admin", "v2.0"
				],
				"convert_placeholders": true,
				"custom_translation_status_ids": [1, 2, 3],
				"custom_translation_status_inserted_keys": true,
				"custom_translation_status_updated_keys": true,
				"custom_translation_status_skipped_keys": false,
				"queue": true
			}`

			req := new(bytes.Buffer)
			_ = json.Compact(req, []byte(data))

			testBody(t, r, req.String())

			_, _ = fmt.Fprint(w, `{
				"project_id": "`+testProjectID+`",
				"process": {
					"process_id": "2e0559e60e856555fbc15bdf78ab2b0ca3406e8f",
					"type": "file-import",
					"status": "queued",
					"message": "",
					"created_by": 1234,
					"created_by_email": "example@example.com",
					"created_at": "2020-04-20 13:43:43 (Etc/UTC)",
					"created_at_timestamp": 1587390223
				}
			}`)
		})

	r, err := client.Files().Upload(testProjectID, FileUpload{
		Data:                       "D94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGL.....",
		Filename:                   "index.json",
		LangISO:                    "en",
		Tags:                       []string{"index", "admin", "v2.0"},
		ConvertPlaceholders:        Bool(true),
		CustomTranslationStatusIds: []int64{1, 2, 3},
	})
	if err != nil {
		t.Errorf("Files.Upload returned error: %v", err)
	}

	want := FileUploadResponse{
		WithProjectID: WithProjectID{
			ProjectID: testProjectID,
		},
		Process: QueuedProcess{
			ID:      "2e0559e60e856555fbc15bdf78ab2b0ca3406e8f",
			Type:    "file-import",
			Status:  "queued",
			Message: "",
			WithCreationUser: WithCreationUser{
				CreatedBy:      1234,
				CreatedByEmail: "example@example.com",
			},
			WithCreationTime: WithCreationTime{
				CreatedAt:   "2020-04-20 13:43:43 (Etc/UTC)",
				CreatedAtTs: 1587390223,
			},
		},
	}

	if !reflect.DeepEqual(r, want) {
		t.Errorf("Files.Upload returned %+v, want %+v", r, want)
	}
}
