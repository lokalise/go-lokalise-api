package lokalise

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestTaskService_Create(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(
		fmt.Sprintf("/projects/%s/tasks", testProjectID),
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			testMethod(t, r, "POST")
			testHeader(t, r, apiTokenHeader, testApiToken)
			data := `{
				"title": "Voicemail messages",
				"languages": [
					{
						"language_iso": "fi",
						"users": [
							421
						]
					}
				],
				"description": "Need your help with some voicemail message translation. Thanks!",
				"due_date": "2018-12-31 12:00:00",
				"keys": [
					11212, 11241, 11245
				],
				"auto_close_languages": true,
				"auto_close_task": true,
				"task_type": "translation",
				"parent_task_id": 12345,
				"closing_tags": ["tag_one", "tag_two"],
				"do_lock_translations": true,
				"custom_translation_status_ids": [77, 85, 86]
			}`

			req := new(bytes.Buffer)
			_ = json.Compact(req, []byte(data))

			testBody(t, r, req.String())

			_, _ = fmt.Fprint(w, `{
				"project_id": "`+testProjectID+`",
				"task": {
					"task_id": 55392,
					"title": "Voicemail messages",
					"description": "Need your help with some voicemail message translation. Thanks!",
					"status": "in progress",
					"progress": 0,
					"can_be_parent": true,
					"task_type": "review",
					"parent_task_id": 12345,
					"closing_tags": ["tag_one", "tag_two"],
					"do_lock_translations": true,
					"due_date": "2018-12-31 12:00:00 (Etc/UTC)",
					"keys_count": 3,
					"words_count": 91,
					"created_at": "2018-12-31 12:00:00 (Etc/UTC)",
					"created_at_timestamp": 1546257600,
					"created_by": 420,
					"created_by_email": "manager@yourcompany.com",
					"languages": [
						{
							"language_iso": "fi",
							"users": [
								{
									"user_id": 421,
									"email": "jdoe@mycompany.com",
									"fullname": "John Doe"
								}
							],
							"groups": [],
							"keys": [
								11212, 11241, 11245
							],
							"status": "created",
							"progress": 0,
							"initial_tm_leverage": {},
							"keys_count": 3,
							"words_count": 91,
							"completed_at": null,
							"completed_at_timestamp": null,
							"completed_by": null,
							"completed_by_email": null
						}                   
					],
					"auto_close_languages": true,
					"auto_close_task": true,
					"completed_at": null,
					"completed_at_timestamp": null,
					"completed_by": null,
					"completed_by_email": null,
					"custom_translation_status_ids": [77, 85, 86]
				}
			}`)
		})

	r, err := client.Tasks().Create(testProjectID, CreateTask{
		Title: "Voicemail messages",
		Languages: []CreateTaskLang{{
			LanguageISO: "fi",
			Users:       []int64{421},
		}},
		Description:                "Need your help with some voicemail message translation. Thanks!",
		DueDate:                    "2018-12-31 12:00:00",
		Keys:                       []int64{11212, 11241, 11245},
		AutoCloseLanguages:         Bool(true),
		AutoCloseTask:              Bool(true),
		TaskType:                   "translation",
		ParentTaskID:               12345,
		ClosingTags:                []string{"tag_one", "tag_two"},
		LockTranslations:           true,
		CustomTranslationStatusIDs: []int64{77, 85, 86},
	})
	if err != nil {
		t.Errorf("Tasks.Create returned error: %v", err)
	}

	want := Task{
		WithCreationTime: WithCreationTime{
			CreatedAt:   "2018-12-31 12:00:00 (Etc/UTC)",
			CreatedAtTs: 1546257600,
		},
		WithCreationUser: WithCreationUser{
			CreatedBy:      420,
			CreatedByEmail: "manager@yourcompany.com",
		},
		TaskID:           55392,
		Title:            "Voicemail messages",
		Description:      "Need your help with some voicemail message translation. Thanks!",
		Status:           "in progress",
		Progress:         0,
		DueDate:          "2018-12-31 12:00:00 (Etc/UTC)",
		KeysCount:        3,
		WordsCount:       91,
		CanBeParent:      true,
		TaskType:         "review",
		ParentTaskID:     12345,
		ClosingTags:      []string{"tag_one", "tag_two"},
		LockTranslations: true,
		Languages: []TaskLanguage{
			{
				LanguageISO: "fi",
				Users: []TaskUser{
					{
						WithUserID: WithUserID{
							UserID: 421,
						},
						Email:    "jdoe@mycompany.com",
						Fullname: "John Doe",
					},
				},
				Groups:            []TaskGroup{},
				Keys:              []int64{11212, 11241, 11245},
				Status:            "created",
				Progress:          0,
				InitialTMLeverage: map[string]int64{},
				KeysCount:         3,
				WordsCount:        91,
				CompletedAt:       "",
				CompletedAtTs:     0,
				CompletedBy:       0,
				CompletedByEmail:  "",
			},
		},
		AutoCloseLanguages:         true,
		AutoCloseTask:              true,
		CompletedAt:                "",
		CompletedAtTs:              0,
		CompletedBy:                0,
		CompletedByEmail:           "",
		CustomTranslationStatusIDs: []int64{77, 85, 86},
	}

	if !reflect.DeepEqual(r.Task, want) {
		t.Errorf("Tasks.Create returned %+v, want %+v", r.Task, want)
	}
}

func TestTaskService_Delete(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(
		fmt.Sprintf("/projects/%s/tasks/%d", testProjectID, 421),
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			testMethod(t, r, "DELETE")
			testHeader(t, r, apiTokenHeader, testApiToken)

			_, _ = fmt.Fprint(w, `{
				"project_id": "`+testProjectID+`",
				"task_deleted": true
			}`)
		},
	)

	r, err := client.Tasks().Delete(testProjectID, 421)
	if err != nil {
		t.Errorf("Tasks.Delete returned error: %v", err)
	}

	want := DeleteTaskResponse{
		WithProjectID: WithProjectID{
			ProjectID: testProjectID,
		},
		Deleted: true,
	}

	if !reflect.DeepEqual(r, want) {
		t.Errorf("Tasks.Delete returned %+v, want %+v", r, want)
	}
}

func TestTaskService_List(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(
		fmt.Sprintf("/projects/%s/tasks", testProjectID),
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			testMethod(t, r, "GET")
			testHeader(t, r, apiTokenHeader, testApiToken)

			_, _ = fmt.Fprint(w, `{
			"project_id": "`+testProjectID+`",
			"tasks": [
					{
						"task_id": 403021
					}
				]
			}`)
		})

	r, err := client.Tasks().List(testProjectID)
	if err != nil {
		t.Errorf("Tasks.List returned error: %v", err)
	}

	want := []Task{
		{
			TaskID: 403021,
		},
	}

	if !reflect.DeepEqual(r.Tasks, want) {
		t.Errorf("Tasks.List returned %+v, want %+v", r.Tasks, want)
	}
}

func TestTaskService_Retrieve(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(
		fmt.Sprintf("/projects/%s/tasks/%d", testProjectID, 403021),
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			testMethod(t, r, "GET")
			testHeader(t, r, apiTokenHeader, testApiToken)

			_, _ = fmt.Fprint(w, `{
				"project_id": "`+testProjectID+`",
				"task": {
					"task_id": 403021
				}
			}`)
		})

	r, err := client.Tasks().Retrieve(testProjectID, 403021)
	if err != nil {
		t.Errorf("Tasks.Retrieve returned error: %v", err)
	}

	want := Task{
		TaskID: 403021,
	}

	if !reflect.DeepEqual(r.Task, want) {
		t.Errorf("Tasks.Retrieve returned %+v, want %+v", r.Task, want)
	}
}

func TestTaskService_Update(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(
		fmt.Sprintf("/projects/%s/tasks/%d", testProjectID, 403021),
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			testMethod(t, r, "PUT")
			testHeader(t, r, apiTokenHeader, testApiToken)
			data := `{
				"due_date": "2019-12-31 12:00:00",
				"auto_close_languages": false,
				"auto_close_task": false
			}`

			req := new(bytes.Buffer)
			_ = json.Compact(req, []byte(data))

			testBody(t, r, req.String())

			_, _ = fmt.Fprint(w, `{
				"project_id": "`+testProjectID+`",
				"task": {
					"task_id": 403021
				}
			}`)
		})

	r, err := client.Tasks().Update(testProjectID, 403021, UpdateTask{
		DueDate:            "2019-12-31 12:00:00",
		AutoCloseLanguages: Bool(false),
		AutoCloseTask:      Bool(false),
	})
	if err != nil {
		t.Errorf("Tasks.Update returned error: %v", err)
	}

	want := Task{
		TaskID: 403021,
	}

	if !reflect.DeepEqual(r.Task, want) {
		t.Errorf("Tasks.Update returned %+v, want %+v", r.Task, want)
	}
}
