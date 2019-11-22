package lokalise

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestTeamUserGroupService_AddMembers(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(
		fmt.Sprintf("/teams/%d/groups/%d/members/add", 444, 50031),
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			testMethod(t, r, "PUT")
			testHeader(t, r, apiTokenHeader, testApiToken)
			data := `{
				"users": [
					22212
				]
			}`

			req := new(bytes.Buffer)
			_ = json.Compact(req, []byte(data))

			testBody(t, r, req.String())

			_, _ = fmt.Fprint(w, `{
				"team_id": 444,
				"group": {
					"group_id": 50031,
					"name": "Proofreading admins",
					"permissions": {
						"is_admin": true,
						"is_reviewer": true,
						"admin_rights": [
							"upload",
							"download",
							"tasks",
							"contributors",
							"screenshots",
							"keys",
							"languages",
							"settings",
							"activity",
							"statistics"
						],
						"languages": []
					},
					"created_at": "2019-01-07 14:21:44 (Etc/UTC)",
					"created_at_timestamp": 1546257600,
					"team_id": 444,
					"projects": [],
					"members": [
						22212
					]
				}
			}`)
		})

	r, err := client.TeamUserGroups().AddMembers(444, 50031, []int64{22212})
	if err != nil {
		t.Errorf("TeamUserGroups.AddMembers returned error: %v", err)
	}

	want := TeamUserGroup{
		WithCreationTime: WithCreationTime{
			CreatedAt:   "2019-01-07 14:21:44 (Etc/UTC)",
			CreatedAtTs: 1546257600,
		},
		WithTeamID: WithTeamID{
			TeamID: 444,
		},
		GroupID: 50031,
		Name:    "Proofreading admins",
		Permissions: &Permission{
			IsAdmin:    true,
			IsReviewer: true,
			Languages:  []Language{},
			AdminRights: []string{
				"upload",
				"download",
				"tasks",
				"contributors",
				"screenshots",
				"keys",
				"languages",
				"settings",
				"activity",
				"statistics",
			},
		},
		Projects: []string{},
		Members:  []int64{22212},
	}

	if !reflect.DeepEqual(r.Group, want) {
		t.Errorf("TeamUserGroups.AddMembers returned %+v, want %+v", r.Group, want)
	}
}

func TestTeamUserGroupService_AddProjects(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(
		fmt.Sprintf("/teams/%d/groups/%d/projects/add", 444, 50031),
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			testMethod(t, r, "PUT")
			testHeader(t, r, apiTokenHeader, testApiToken)
			data := `{
				"projects": [
					"` + testProjectID + `"
				]
			}`

			req := new(bytes.Buffer)
			_ = json.Compact(req, []byte(data))

			testBody(t, r, req.String())

			_, _ = fmt.Fprint(w, `{
				"team_id": 444,
				"group": {
					"group_id": 50031
				}
			}`)
		})

	r, err := client.TeamUserGroups().AddProjects(444, 50031, []string{testProjectID})
	if err != nil {
		t.Errorf("TeamUserGroups.AddProjects returned error: %v", err)
	}

	want := TeamUserGroup{
		GroupID: 50031,
	}

	if !reflect.DeepEqual(r.Group, want) {
		t.Errorf("TeamUserGroups.AddProjects returned %+v, want %+v", r.Group, want)
	}
}

func TestTeamUserGroupService_Create(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(
		fmt.Sprintf("/teams/%d/groups", 444),
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			testMethod(t, r, "POST")
			testHeader(t, r, apiTokenHeader, testApiToken)
			data := `{
				"name": "Proofreaders",
				"is_admin": false,
				"is_reviewer": true,
				"languages": {
					"reference": [],
					"contributable": [640]
				}
			}`

			req := new(bytes.Buffer)
			_ = json.Compact(req, []byte(data))

			testBody(t, r, req.String())

			_, _ = fmt.Fprint(w, `{
				"team_id": 444,
				"group": {
					"group_id": 50031
				}
			}`)
		})

	r, err := client.TeamUserGroups().Create(444, NewGroup{
		Name:        "Proofreaders",
		IsAdmin:     false,
		IsReviewer:  true,
		AdminRights: nil,
		Languages: NewGroupLanguages{
			Reference:     []int64{},
			Contributable: []int64{640},
		},
	})
	if err != nil {
		t.Errorf("TeamUserGroups.Create returned error: %v", err)
	}

	want := TeamUserGroup{
		GroupID: 50031,
	}

	if !reflect.DeepEqual(r.Group, want) {
		t.Errorf("TeamUserGroups.Create returned %+v, want %+v", r.Group, want)
	}
}

func TestTeamUserGroupService_Delete(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(
		fmt.Sprintf("/teams/%d/groups/%d", 444, 50031),
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			testMethod(t, r, "DELETE")
			testHeader(t, r, apiTokenHeader, testApiToken)

			_, _ = fmt.Fprint(w, `{
				"team_id": 444,
				"group_deleted": true
			}`)
		},
	)

	r, err := client.TeamUserGroups().Delete(444, 50031)
	if err != nil {
		t.Errorf("TeamUserGroups.Delete returned error: %v", err)
	}

	want := DeleteGroupResponse{
		WithTeamID: WithTeamID{
			TeamID: 444,
		},
		IsDeleted: true,
	}

	if !reflect.DeepEqual(r, want) {
		t.Errorf("TeamUserGroups.Delete returned %+v, want %+v", r, want)
	}
}

func TestTeamUserGroupService_List(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(
		fmt.Sprintf("/teams/%d/groups", 444),
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			testMethod(t, r, "GET")
			testHeader(t, r, apiTokenHeader, testApiToken)

			_, _ = fmt.Fprint(w, `{
				"team_id": 444,
				"user_groups": [
					{
						"group_id": 50031
					}
				]
			}`)
		})

	r, err := client.TeamUserGroups().List(444)
	if err != nil {
		t.Errorf("TeamUserGroups.List returned error: %v", err)
	}

	want := []TeamUserGroup{
		{
			GroupID: 50031,
		},
	}

	if !reflect.DeepEqual(r.UserGroups, want) {
		t.Errorf("Screenshots.List returned %+v, want %+v", r.UserGroups, want)
	}
}

func TestTeamUserGroupService_RemoveMembers(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(
		fmt.Sprintf("/teams/%d/groups/%d/members/remove", 444, 50031),
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			testMethod(t, r, "PUT")
			testHeader(t, r, apiTokenHeader, testApiToken)
			data := `{
				"users": [
					22212
				]
			}`

			req := new(bytes.Buffer)
			_ = json.Compact(req, []byte(data))

			testBody(t, r, req.String())
			_, _ = fmt.Fprint(w, `{
				"team_id": 444,
				"group":{
					"group_id": 50031
				}
			}`)
		})

	r, err := client.TeamUserGroups().RemoveMembers(444, 50031, []int64{22212})
	if err != nil {
		t.Errorf("TeamUserGroups.RemoveMembers returned error: %v", err)
	}

	want := TeamUserGroup{
		GroupID: 50031,
	}

	if !reflect.DeepEqual(r.Group, want) {
		t.Errorf("Screenshots.RemoveMembers returned %+v, want %+v", r.Group, want)
	}
}

func TestTeamUserGroupService_RemoveProjects(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(
		fmt.Sprintf("/teams/%d/groups/%d/projects/remove", 444, 50031),
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			testMethod(t, r, "PUT")
			testHeader(t, r, apiTokenHeader, testApiToken)
			data := `{
				"projects": [
					"` + testProjectID + `"
				]
			}`

			req := new(bytes.Buffer)
			_ = json.Compact(req, []byte(data))

			testBody(t, r, req.String())
			_, _ = fmt.Fprint(w, `{
				"team_id": 444,
				"group":{
					"group_id": 50031
				}
			}`)
		})

	r, err := client.TeamUserGroups().RemoveProjects(444, 50031, []string{testProjectID})
	if err != nil {
		t.Errorf("TeamUserGroups.RemoveProjects returned error: %v", err)
	}

	want := TeamUserGroup{
		GroupID: 50031,
	}

	if !reflect.DeepEqual(r.Group, want) {
		t.Errorf("Screenshots.RemoveProjects returned %+v, want %+v", r.Group, want)
	}
}

func TestTeamUserGroupService_Retrieve(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(
		fmt.Sprintf("/teams/%d/groups/%d", 444, 50031),
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			testMethod(t, r, "GET")
			testHeader(t, r, apiTokenHeader, testApiToken)
			_, _ = fmt.Fprint(w, `{
				"group_id": 50031
			}`)
		})

	r, err := client.TeamUserGroups().Retrieve(444, 50031)
	if err != nil {
		t.Errorf("TeamUserGroups.Retrieve returned error: %v", err)
	}

	want := TeamUserGroup{
		GroupID: 50031,
	}

	if !reflect.DeepEqual(r, want) {
		t.Errorf("Screenshots.Retrieve returned %+v, want %+v", r, want)
	}
}

func TestTeamUserGroupService_Update(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(
		fmt.Sprintf("/teams/%d/groups/%d", 444, 50031),
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			testMethod(t, r, "PUT")
			testHeader(t, r, apiTokenHeader, testApiToken)
			data := `{
				"name": "Proofreading admins",
				"is_admin": true,
				"is_reviewer": true,
				"admin_rights": [
					"upload"
				],
				"languages": {
					"reference": [],
					"contributable": []
				}
			}`

			req := new(bytes.Buffer)
			_ = json.Compact(req, []byte(data))

			testBody(t, r, req.String())
			_, _ = fmt.Fprint(w, `{
				"team_id": 444,
				"group":{
					"group_id": 50031
				}
			}`)
		})

	r, err := client.TeamUserGroups().Update(444, 50031, NewGroup{
		Name:        "Proofreading admins",
		IsReviewer:  true,
		IsAdmin:     true,
		AdminRights: []string{"upload"},
		Languages: NewGroupLanguages{
			Reference:     []int64{},
			Contributable: []int64{},
		},
	})
	if err != nil {
		t.Errorf("TeamUserGroups.Update returned error: %v", err)
	}

	want := TeamUserGroup{
		GroupID: 50031,
	}

	if !reflect.DeepEqual(r.Group, want) {
		t.Errorf("TeamUserGroups.Update returned %+v, want %+v", r.Group, want)
	}
}
