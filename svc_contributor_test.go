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
						"user_id": 420
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
					"user_id": 421
				}
			}`)
		})

	r, err := client.Contributors().Retrieve(testProjectID, 421)
	if err != nil {
		t.Errorf("Contributors.Retrieve returned error: %v", err)
	}

	want := Contributor{
		WithUserID: WithUserID{UserID: 421},
	}

	if !reflect.DeepEqual(r.Contributor, want) {
		t.Errorf("Contributors.Retrieve returned %+v, want %+v", r.Contributor, want)
	}
}

func TestContributorService_Update(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(
		fmt.Sprintf("/projects/%s/contributors/%d", testProjectID, 421),
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			testMethod(t, r, "PUT")
			testHeader(t, r, apiTokenHeader, testApiToken)
			data := `{
				"is_admin": true,
				"is_reviewer":false
			}`

			req := new(bytes.Buffer)
			_ = json.Compact(req, []byte(data))

			testBody(t, r, req.String())

			_, _ = fmt.Fprint(w, `{
				"project_id": "`+testProjectID+`",
				"contributor": {
					"user_id": 421
				}
			}`)
		})

	r, err := client.Contributors().Update(testProjectID, 421, Permission{IsAdmin: true})
	if err != nil {
		t.Errorf("Contributors.Update returned error: %v", err)
	}

	want := Contributor{
		WithUserID: WithUserID{UserID: 421},
	}

	if !reflect.DeepEqual(r.Contributor, want) {
		t.Errorf("Contributors.Update returned %+v, want %+v", r.Contributor, want)
	}
}
