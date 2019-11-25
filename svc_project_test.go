package lokalise

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestProjectService_Create(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(
		"/projects",
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			testMethod(t, r, "POST")
			testHeader(t, r, apiTokenHeader, testApiToken)
			data := `{
				"team_id": 12345,
				"name": "TheApp Project",
				"description": "iOS + Android strings of TheApp. https://theapp.com",
				"languages": [
					{
						"lang_iso": "en",
						"custom_iso": "en-us"
					}
				],
				"base_lang_iso": "en-us"
			}`

			req := new(bytes.Buffer)
			_ = json.Compact(req, []byte(data))

			testBody(t, r, req.String())

			_, _ = fmt.Fprint(w, `{
				"project_id": "3002780358964f9bab5a92.87762498",
				"project_type": "localization_files",
				"name": "TheApp Project",
				"description": "iOS + Android strings of TheApp. https://theapp.com",
				"created_at": "2018-12-31 12:00:00 (Etc/UTC)",
				"created_at_timestamp": 1546257600,
				"created_by": 420,
				"created_by_email": "user@mycompany.com",
				"team_id": 12345,
				"base_language_id": 640,
				"base_language_iso": "en-us",
				"settings": {
					"per_platform_key_names": false,
					"reviewing": false,
					"upvoting": false,
					"auto_toggle_unverified": true,
					"offline_translation": false,
					"key_editing": true,
					"inline_machine_translations": true
				},
				"statistics": {
					"progress_total": 0,
					"keys_total": 0,
					"team": 1,
					"base_words": 0,
					"qa_issues_total": 0,
					"qa_issues": {
						"not_reviewed": 0,
						"unverified": 0,
						"spelling_grammar": 0,
						"inconsistent_placeholders": 0,
						"inconsistent_html": 0,
						"different_number_of_urls": 0,
						"different_urls": 0,
						"leading_whitespace": 0,
						"trailing_whitespace": 0,
						"different_number_of_email_address": 0,
						"different_email_address": 0,
						"different_brackets": 0,
						"different_numbers": 0,
						"double_space": 1,
						"special_placeholder": 0
					},
					"languages": [
						{
							"language_id": 640,
							"language_iso": "en-us",
							"progress": 0,
							"words_to_do": 0
						}
					]
				}
			}`)
		},
	)

	r, err := client.Projects().Create(NewProject{
		WithTeamID:  WithTeamID{TeamID: 12345},
		Name:        "TheApp Project",
		Description: "iOS + Android strings of TheApp. https://theapp.com",
		Languages: []NewLanguage{
			{
				LangISO:   "en",
				CustomISO: "en-us",
			},
		},
		BaseLangISO: "en-us",
	})
	if err != nil {
		t.Errorf("PaymentCards.Create returned error: %v", err)
	}

	want := Project{
		WithProjectID: WithProjectID{
			ProjectID: "3002780358964f9bab5a92.87762498",
		},
		WithTeamID: WithTeamID{
			TeamID: 12345,
		},
		WithCreationTime: WithCreationTime{
			CreatedAt:   "2018-12-31 12:00:00 (Etc/UTC)",
			CreatedAtTs: 1546257600,
		},
		WithCreationUser: WithCreationUser{
			CreatedBy:      420,
			CreatedByEmail: "user@mycompany.com",
		},
		Name:        "TheApp Project",
		Description: "iOS + Android strings of TheApp. https://theapp.com",
		Type:        "localization_files",
		BaseLangID:  640,
		BaseLangISO: "en-us",
		Settings: &ProjectSettings{
			PerPlatformKeyNames:       false,
			Reviewing:                 false,
			Upvoting:                  false,
			AutoToggleUnverified:      true,
			OfflineTranslation:        false,
			KeyEditing:                true,
			InlineMachineTranslations: true,
		},
		Statistics: &ProjectStatistics{
			ProgressTotal: 0,
			KeysTotal:     0,
			Team:          1,
			BaseWords:     0,
			QAIssuesTotal: 0,
			QAIssues: QAIssues{
				NotReviewed:                   0,
				Unverified:                    0,
				SpellingGrammar:               0,
				InconsistentPlaceholders:      0,
				InconsistentHtml:              0,
				DifferentNumberOfUrls:         0,
				DifferentUrls:                 0,
				LeadingWhitespace:             0,
				TrailingWhitespace:            0,
				DifferentNumberOfEmailAddress: 0,
				DifferentEmailAddress:         0,
				DifferentBrackets:             0,
				DifferentNumbers:              0,
				DoubleSpace:                   1,
				SpecialPlaceholder:            0,
			},
			Languages: []LanguageStatistics{
				{
					LanguageID:  640,
					LanguageISO: "en-us",
					Progress:    0,
					WordsToDo:   0,
				},
			},
		},
	}

	if !reflect.DeepEqual(r, want) {
		t.Errorf("PaymentCards.Create returned %+v, want %+v", r, want)
	}
}

func TestProjectService_Delete(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(
		fmt.Sprintf("/projects/%s", testProjectID),
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			testMethod(t, r, "DELETE")
			testHeader(t, r, apiTokenHeader, testApiToken)

			_, _ = fmt.Fprint(w, `{
				"project_id": "`+testProjectID+`",
				"project_deleted": true
			}`)
		})

	r, err := client.Projects().Delete(testProjectID)
	if err != nil {
		t.Errorf("Projects.Delete returned error: %v", err)
	}

	want := DeleteProjectResponse{
		WithProjectID: WithProjectID{
			ProjectID: testProjectID,
		},
		Deleted: true,
	}

	if !reflect.DeepEqual(r, want) {
		t.Errorf("Projects.Create returned %+v, want %+v", r, want)
	}
}

func TestProjectService_List(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(
		"/projects",
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			testMethod(t, r, "GET")
			testHeader(t, r, apiTokenHeader, testApiToken)

			_, _ = fmt.Fprint(w, `{
			"projects": [
				{
					"project_id": "`+testProjectID+`"
				}
			]
		}`)
		})

	r, err := client.Projects().List()
	if err != nil {
		t.Errorf("Projects.List returned error: %v", err)
	}

	want := []Project{
		{
			WithProjectID: WithProjectID{
				ProjectID: testProjectID,
			},
		},
	}

	if !reflect.DeepEqual(r.Projects, want) {
		t.Errorf("Projects.List returned %+v, want %+v", r.Projects, want)
	}
}

func TestProjectService_Retrieve(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(
		fmt.Sprintf("/projects/%s", testProjectID),
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			testMethod(t, r, "GET")
			testHeader(t, r, apiTokenHeader, testApiToken)

			_, _ = fmt.Fprint(w, `{
				"project_id": "`+testProjectID+`"  
			}`)
		})

	r, err := client.Projects().Retrieve(testProjectID)
	if err != nil {
		t.Errorf("Projects.Retrieve returned error: %v", err)
	}

	want := Project{
		WithProjectID: WithProjectID{
			ProjectID: testProjectID,
		},
	}

	if !reflect.DeepEqual(r, want) {
		t.Errorf("Projects.Retrieve returned %+v, want %+v", r, want)
	}
}

func TestProjectService_Truncate(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(
		fmt.Sprintf("/projects/%s/empty", testProjectID),
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			testMethod(t, r, "PUT")
			testHeader(t, r, apiTokenHeader, testApiToken)

			_, _ = fmt.Fprint(w, `{
				"project_id": "`+testProjectID+`",
				"keys_deleted": true
			}`)
		})

	r, err := client.Projects().Truncate(testProjectID)
	if err != nil {
		t.Errorf("Languages.Truncate returned error: %v", err)
	}

	want := TruncateProjectResponse{
		WithProjectID: WithProjectID{
			ProjectID: testProjectID,
		},
		KeysDeleted: true,
	}

	if !reflect.DeepEqual(r, want) {
		t.Errorf("Languages.Truncate returned %+v, want %+v", r, want)
	}
}

func TestProjectService_Update(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(
		fmt.Sprintf("/projects/%s", testProjectID),
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			testMethod(t, r, "PUT")
			testHeader(t, r, apiTokenHeader, testApiToken)
			data := `{
				"name": "TheZapp Project",
				"description": "iOS + Android strings of TheZapp. https://thezapp.com"
			}`

			req := new(bytes.Buffer)
			_ = json.Compact(req, []byte(data))

			testBody(t, r, req.String())

			_, _ = fmt.Fprint(w, `{
				"project_id": "`+testProjectID+`"   
			}`)
		})

	r, err := client.Projects().Update(testProjectID, UpdateProject{
		Name:        "TheZapp Project",
		Description: "iOS + Android strings of TheZapp. https://thezapp.com",
	})
	if err != nil {
		t.Errorf("Languages.Update returned error: %v", err)
	}

	want := Project{
		WithProjectID: WithProjectID{
			ProjectID: testProjectID,
		},
	}

	if !reflect.DeepEqual(r, want) {
		t.Errorf("Languages.Update returned %+v, want %+v", r, want)
	}
}
