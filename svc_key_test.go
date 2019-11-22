package lokalise

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestKeyService_BulkDelete(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(
		fmt.Sprintf("/projects/%s/keys", testProjectID),
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			testMethod(t, r, "DELETE")
			testHeader(t, r, apiTokenHeader, testApiToken)
			data := `{
				"keys": [
					12345, 12346
				]
			}`

			req := new(bytes.Buffer)
			_ = json.Compact(req, []byte(data))

			testBody(t, r, req.String())
			_, _ = fmt.Fprint(w, `{
				"project_id": "3002780358964f9bab5a92.87762498",
				"keys_removed": true,
				"keys_locked": 0
			}`)
		})

	r, err := client.Keys().BulkDelete(testProjectID, []int64{12345, 12346})
	if err != nil {
		t.Errorf("Languages.BulkDelete returned error: %v", err)
	}

	want := DeleteKeysResponse{
		WithProjectID: WithProjectID{
			ProjectID: testProjectID,
		},
		AreRemoved:     true,
		NumberOfLocked: 0,
	}

	if !reflect.DeepEqual(r, want) {
		t.Errorf("Keys.BulkDelete returned %+v, want %+v", r, want)
	}
}

func TestKeyService_BulkUpdate(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(
		fmt.Sprintf("/projects/%s/keys", testProjectID),
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			testMethod(t, r, "PUT")
			testHeader(t, r, apiTokenHeader, testApiToken)
			data := `{
				"keys": [
					{
						"key_id": 331223,
						"key_name": "index.welcome",
						"description": "Index app welcome",
						"platforms": [
							"web"
					   ]
					}
				]
			}`

			req := new(bytes.Buffer)
			_ = json.Compact(req, []byte(data))

			testBody(t, r, req.String())
			_, _ = fmt.Fprint(w, `{
				"project_id": "`+testProjectID+`",
				"keys": [
					{
						"key_id": 331223
					}
				],
				"errors": []
			}`)
		})

	r, err := client.Keys().BulkUpdate(testProjectID, []BulkUpdateKey{
		{
			KeyID: 331223,
			NewKey: NewKey{
				KeyName:     "index.welcome",
				Description: "Index app welcome",
				Platforms: []string{
					"web",
				},
			},
		},
	})
	if err != nil {
		t.Errorf("Keys.BulkUpdate returned error: %v", err)
	}

	want := []Key{
		{
			KeyID: 331223,
		},
	}

	if !reflect.DeepEqual(r.Keys, want) {
		t.Errorf("Keys.BulkUpdate returned %+v, want %+v", r.Keys, want)
	}
}

func TestKeyService_Create(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(
		fmt.Sprintf("/projects/%s/keys", testProjectID),
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			testMethod(t, r, "POST")
			testHeader(t, r, apiTokenHeader, testApiToken)
			data := `{
				"keys": [
					{
						"key_name": "index.welcome",
						"description": "Index app welcome",
						"platforms": [
							"web"
						],
						"translations": [
							{
								"language_iso": "en",
								"translation": "Welcome"
							}
						]
					}
				]
			}`

			req := new(bytes.Buffer)
			_ = json.Compact(req, []byte(data))

			testBody(t, r, req.String())

			_, _ = fmt.Fprint(w, `{
				"project_id": "3002780358964f9bab5a92.87762498",
				"keys": [
					{
						"key_id": 331223,
						"created_at": "2018-12-31 12:00:00 (Etc/UTC)",
						"created_at_timestamp": 1546257600,
						"key_name": {
							"ios": "index.welcome",
							"android": "index.welcome",
							"web": "index.welcome",
							"other": "index.welcome"
						},
						"filenames": {
							"ios": "",
							"android": "",
							"web": "",
							"other": ""
						},
						"description": "Index app welcome",
						"platforms": [
							"web"
						],
						"tags": [],
						"comments": [],
						"screenshots": [],
						"translations": [
							{
								"translation_id": 444921,
								"key_id": 331223,
								"language_iso": "en",
								"translation": "Welcome",
								"modified_by": 420,
								"modified_by_email": "user@mycompany.com",
								"modified_at": "2018-12-31 12:00:00 (Etc/UTC)",
								"modified_at_timestamp": 1546257600,
								"is_reviewed": false,
								"reviewed_by": 0,
								"words": 0
							}
						]
					}
				],
				"errors": [
					{
						"message": "This key name is already taken",
						"code": 400,
						"key": {
							"key_name": "index.hello"
						}
					}
				]
			}`)
		})

	r, err := client.Keys().Create(testProjectID, []NewKey{
		{
			KeyName:     "index.welcome",
			Description: "Index app welcome",
			Platforms:   []string{"web"},
			Translations: []NewTranslation{
				{
					LanguageISO: "en",
					Translation: "Welcome",
				},
			},
		},
	})
	if err != nil {
		t.Errorf("Keys.Create returned error: %v", err)
	}

	want := []Key{
		{
			KeyID: 331223,
			WithCreationTime: WithCreationTime{
				CreatedAt:   "2018-12-31 12:00:00 (Etc/UTC)",
				CreatedAtTs: 1546257600,
			},
			KeyName: PlatformStrings{
				Ios:     "index.welcome",
				Android: "index.welcome",
				Web:     "index.welcome",
				Other:   "index.welcome",
			},
			Filenames: PlatformStrings{
				Ios:     "",
				Android: "",
				Web:     "",
				Other:   "",
			},
			Description: "Index app welcome",
			Platforms: []string{
				"web",
			},
			Tags:        []string{},
			Comments:    []Comment{},
			Screenshots: []Screenshot{},
			Translations: []Translation{
				{
					TranslationID:   444921,
					KeyID:           331223,
					LanguageISO:     "en",
					Translation:     "Welcome",
					ModifiedBy:      420,
					ModifiedByEmail: "user@mycompany.com",
					ModifiedAt:      "2018-12-31 12:00:00 (Etc/UTC)",
					ModifiedAtTs:    1546257600,
					IsReviewed:      false,
					ReviewedBy:      0,
					Words:           0,
				},
			},
		},
	}

	if !reflect.DeepEqual(r.Keys, want) {
		t.Errorf("Keys.Create returned %+v, want %+v", r.Keys, want)
	}
}

func TestKeyService_Delete(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(
		fmt.Sprintf("/projects/%s/keys/%d", testProjectID, 640),
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			testMethod(t, r, "DELETE")
			testHeader(t, r, apiTokenHeader, testApiToken)

			_, _ = fmt.Fprint(w, `{
				"project_id": "`+testProjectID+`",
				"key_removed": false,
				"keys_locked": 1
			}`)
		})

	r, err := client.Keys().Delete(testProjectID, 640)
	if err != nil {
		t.Errorf("Keys.Delete returned error: %v", err)
	}

	want := DeleteKeyResponse{
		WithProjectID: WithProjectID{
			ProjectID: testProjectID,
		},
		IsRemoved:      false,
		NumberOfLocked: 1,
	}

	if !reflect.DeepEqual(r, want) {
		t.Errorf("Keys.Delete returned %+v, want %+v", r, want)
	}
}

func TestKeyService_List(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(
		fmt.Sprintf("/projects/%s/keys", testProjectID),
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			testMethod(t, r, "GET")
			testHeader(t, r, apiTokenHeader, testApiToken)

			_, _ = fmt.Fprint(w, `{
				"keys": [
					{
						"key_id": 640
					}
				]
			}`)
		})

	r, err := client.Keys().List(testProjectID)
	if err != nil {
		t.Errorf("Keys.List returned error: %v", err)
	}

	want := []Key{
		{
			KeyID: 640,
		},
	}

	if !reflect.DeepEqual(r.Keys, want) {
		t.Errorf("Keys.List returned %+v, want %+v", r.Keys, want)
	}
}

func TestKeyService_Retrieve(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(
		fmt.Sprintf("/projects/%s/keys/%d", testProjectID, 640),
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			testMethod(t, r, "GET")
			testHeader(t, r, apiTokenHeader, testApiToken)

			_, _ = fmt.Fprint(w, `{
				"project_id": "`+testProjectID+`",
				"key": {
					"key_id": 640
				}   
			}`)
		})

	r, err := client.Keys().Retrieve(testProjectID, 640)
	if err != nil {
		t.Errorf("Keys.Retrieve returned error: %v", err)
	}

	want := Key{
		KeyID: 640,
	}

	if !reflect.DeepEqual(r.Key, want) {
		t.Errorf("Keys.Retrieve returned %+v, want %+v", r.Key, want)
	}
}

func TestKeyService_Update(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(
		fmt.Sprintf("/projects/%s/keys/%d", testProjectID, 640),
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			testMethod(t, r, "PUT")
			testHeader(t, r, apiTokenHeader, testApiToken)
			data := `{
				"description": "Index app welcome",
				"platforms": [
					"web","other"
				]
			}`

			req := new(bytes.Buffer)
			_ = json.Compact(req, []byte(data))

			testBody(t, r, req.String())

			_, _ = fmt.Fprint(w, `{
				"project_id": "`+testProjectID+`",
				"key": {
					"key_id": 640
				}   
			}`)
		})

	r, err := client.Keys().Update(testProjectID, 640, NewKey{
		Platforms:   []string{"web", "other"},
		Description: "Index app welcome",
	})
	if err != nil {
		t.Errorf("Keys.Update returned error: %v", err)
	}

	want := Key{
		KeyID: 640,
	}

	if !reflect.DeepEqual(r.Key, want) {
		t.Errorf("Keys.Update returned %+v, want %+v", r.Key, want)
	}
}
