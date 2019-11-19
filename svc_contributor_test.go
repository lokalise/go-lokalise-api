package lokalise

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestContributorService_Create(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(
		fmt.Sprintf("/projects/%s/contributors", testProjectID),
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			testMethod(t, r, "POST")
			testHeader(t, r, apiTokenHeader, testApiToken)
			data := `{
				"contributors": [
					{
						"email": "translator@mycompany.com",
						"fullname": "Mr. Translator",
						"is_admin": false,
						"is_reviewer": true,
						"languages": [
							{
								"lang_iso": "en",
								"is_writable": false
							}
						]
					}
				]
			}`

			req := new(bytes.Buffer)
			_ = json.Compact(req, []byte(data))

			testBody(t, r, req.String())

			_, _ = fmt.Fprint(w, `{
				"project_id": "`+testProjectID+`",
				"contributors": [{
					"user_id": 421,
					"email": "translator@mycompany.com",
					"fullname": "Mr. Translator",
					"is_admin": false,
					"is_reviewer": true,
					"languages": [
						{
							"lang_iso": "en",
							"is_writable": false
						}
					]
				}]
			}`)
		})

	r, err := client.Contributors().Create(testProjectID, []NewContributor{
		{
			Email:    "translator@mycompany.com",
			Fullname: "Mr. Translator",
			Permission: Permission{
				IsAdmin:    false,
				IsReviewer: true,
				Languages: []Language{{
					LangISO:    "en",
					IsWritable: false,
				}},
				AdminRights: nil,
			},
		},
	})
	if err != nil {
		t.Errorf("Contributors.Create returned error: %v", err)
	}

	want := []Contributor{
		{
			WithUserID: WithUserID{UserID: 421},
			Email:      "translator@mycompany.com",
			Fullname:   "Mr. Translator",
			Permission: Permission{
				IsAdmin:    false,
				IsReviewer: true,
				Languages: []Language{{
					LangISO:    "en",
					IsWritable: false,
				}},
				AdminRights: nil,
			},
		},
	}

	if !reflect.DeepEqual(r.Contributors, want) {
		t.Errorf("Contributors.Create returned %+v, want %+v", r.Contributors, want)
	}
}

func TestContributorService_Delete(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(
		fmt.Sprintf("/projects/%s/contributors/%d", testProjectID, 421),
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			testMethod(t, r, "DELETE")
			testHeader(t, r, apiTokenHeader, testApiToken)

			_, _ = fmt.Fprint(w, `{
				"project_id": "`+testProjectID+`",
				"contributor_deleted": true
			}`)
		},
	)

	r, err := client.Contributors().Delete(testProjectID, 421)
	if err != nil {
		t.Errorf("Contributors.Delete returned error: %v", err)
	}

	want := DeleteContributorResponse{
		WithProjectID: WithProjectID{
			ProjectID: testProjectID,
		},
		IsDeleted: true,
	}

	if !reflect.DeepEqual(r, want) {
		t.Errorf("Contributors.Delete returned %+v, want %+v", r, want)
	}
}

func TestContributorService_List(t *testing.T) {

	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(
		fmt.Sprintf("/projects/%s/contributors", testProjectID),
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			testMethod(t, r, "GET")
			testHeader(t, r, apiTokenHeader, testApiToken)

			_, _ = fmt.Fprint(w, `{
			"project_id": "`+testProjectID+`",
			"contributors": [
					{
						"user_id": 420,
						"email": "johndoe@mycompany.com",
						"fullname": "John Doe",
						"created_at": "2018-01-01 12:00:00 (Etc/UTC)",
						"created_at_timestamp": 1546257600,
						"is_admin": true,
						"is_reviewer": true,
						"languages": [
							{
								"lang_id": 640,
								"lang_iso": "en",
								"lang_name": "English",
								"is_writable": true
							}
						],
						"admin_rights": [
							"keys", "languages"
						]
					}
				]
			}`)
		})

	r, err := client.Contributors().List(testProjectID)
	if err != nil {
		t.Errorf("Contributors.List returned error: %v", err)
	}

	want := []Contributor{
		{
			WithUserID: WithUserID{UserID: 420},
			Email:      "johndoe@mycompany.com",
			Fullname:   "John Doe",
			WithCreationTime: WithCreationTime{
				CreatedAt:   "2018-01-01 12:00:00 (Etc/UTC)",
				CreatedAtTs: 1546257600,
			},

			Permission: Permission{
				IsAdmin:     true,
				IsReviewer:  true,
				AdminRights: []string{"keys", "languages"},
				Languages: []Language{
					{
						LangID:     640,
						LangISO:    "en",
						LangName:   "English",
						IsWritable: true,
					},
				},
			},
		},
	}

	if !reflect.DeepEqual(r.Contributors, want) {
		t.Errorf("Contributors.List returned %+v, want %+v", r.Contributors, want)
	}
}

func TestContributorService_Retrieve(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(
		fmt.Sprintf("/projects/%s/contributors/%d", testProjectID, 421),
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			testMethod(t, r, "GET")
			testHeader(t, r, apiTokenHeader, testApiToken)

			_, _ = fmt.Fprint(w, `{
				"project_id": "`+testProjectID+`",
				"contributor": {
					"user_id": 421,
					"email": "translator@mycompany.com",
					"fullname": "Mr. Translator",
					"created_at": "2018-01-01 12:00:00 (Etc/UTC)",
					"created_at_timestamp": 1546257600,
					"is_admin": false,
					"is_reviewer": true,
					"languages": [
						{
							"lang_id": 640,
							"lang_iso": "en",
							"lang_name": "English",
							"is_writable": true
						}
					],
					"admin_rights": [
						"keys", "languages"
					]
				}
			}`)
		})

	r, err := client.Contributors().Retrieve(testProjectID, 421)
	if err != nil {
		t.Errorf("Contributors.Retrieve returned error: %v", err)
	}

	want := Contributor{
		WithUserID: WithUserID{UserID: 421},
		Email:      "translator@mycompany.com",
		Fullname:   "Mr. Translator",
		WithCreationTime: WithCreationTime{
			CreatedAt:   "2018-01-01 12:00:00 (Etc/UTC)",
			CreatedAtTs: 1546257600,
		},

		Permission: Permission{
			IsAdmin:     false,
			IsReviewer:  true,
			Languages: []Language{
				{
					LangID:     640,
					LangISO:    "en",
					LangName:   "English",
					IsWritable: true,
				},
			},
			AdminRights: []string{"keys", "languages"},
		},
	}

	if !reflect.DeepEqual(r.Contributor, want) {
		t.Errorf("Contributors.Retrieve returned %+v, want %+v", r.Contributor, want)
	}
}

func TestContributorService_Update(t *testing.T) {

}
