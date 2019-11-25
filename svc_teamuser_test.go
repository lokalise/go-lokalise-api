package lokalise

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestTeamUserService_Delete(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(
		fmt.Sprintf("/teams/%d/users/%d", 18821, 421),
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			testMethod(t, r, "DELETE")
			testHeader(t, r, apiTokenHeader, testApiToken)

			_, _ = fmt.Fprint(w, `{
				"team_id": 18821,
				"team_user_deleted": true
			}`)
		},
	)

	r, err := client.TeamUsers().Delete(18821, 421)
	if err != nil {
		t.Errorf("TeamUsers.Delete returned error: %v", err)
	}

	want := DeleteTeamUserResponse{
		WithTeamID: WithTeamID{TeamID: 18821},
		Deleted:    true,
	}

	if !reflect.DeepEqual(r, want) {
		t.Errorf("TeamUsers.Delete returned %+v, want %+v", r, want)
	}
}

func TestTeamUserService_List(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(
		fmt.Sprintf("/teams/%d/users", 18821),
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			testMethod(t, r, "GET")
			testHeader(t, r, apiTokenHeader, testApiToken)

			_, _ = fmt.Fprint(w, `{
				"team_id": 18821,
				"team_users": [
					{
						"user_id": 420,
						"email": "jdoe@mycompany.com",
						"fullname": "John Doe",
						"created_at": "2018-12-31 12:00:00 (Etc/UTC)",
						"created_at_timestamp": 1546257600,
						"role": "owner"
					}
				]
			}`)
		})

	r, err := client.TeamUsers().List(18821)
	if err != nil {
		t.Errorf("TeamUsers.List returned error: %v", err)
	}

	want := []TeamUser{
		{
			WithCreationTime: WithCreationTime{
				CreatedAt:   "2018-12-31 12:00:00 (Etc/UTC)",
				CreatedAtTs: 1546257600,
			},
			WithUserID: WithUserID{
				UserID: 420,
			},
			Email:    "jdoe@mycompany.com",
			Fullname: "John Doe",
			Role:     "owner",
		},
	}

	if !reflect.DeepEqual(r.TeamUsers, want) {
		t.Errorf("TeamUsers.List returned %+v, want %+v", r.TeamUsers, want)
	}
}

func TestTeamUserService_Retrieve(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(
		fmt.Sprintf("/teams/%d/users/%d", 18821, 420),
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			testMethod(t, r, "GET")
			testHeader(t, r, apiTokenHeader, testApiToken)

			_, _ = fmt.Fprint(w, `{
				"team_id": 18821,
				"team_user": {
					"user_id": 420
				}
			}`)
		})

	r, err := client.TeamUsers().Retrieve(18821, 420)
	if err != nil {
		t.Errorf("TeamUsers.Retrieve returned error: %v", err)
	}

	want := TeamUser{
		WithUserID: WithUserID{
			UserID: 420,
		},
	}

	if !reflect.DeepEqual(r.TeamUser, want) {
		t.Errorf("TeamUsers.Retrieve returned %+v, want %+v", r.TeamUser, want)
	}
}

func TestTeamUserService_UpdateRole(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(
		fmt.Sprintf("/teams/%d/users/%d", 18821, 420),
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			testMethod(t, r, "PUT")
			testHeader(t, r, apiTokenHeader, testApiToken)
			data := `{
				"role": "admin"
			}`

			req := new(bytes.Buffer)
			_ = json.Compact(req, []byte(data))

			testBody(t, r, req.String())

			_, _ = fmt.Fprint(w, `{
				"project_id": "`+testProjectID+`",
				"team_user": {
					"user_id": 420
				}
			}`)
		})

	r, err := client.TeamUsers().UpdateRole(18821, 420, "admin")
	if err != nil {
		t.Errorf("TeamUsers.UpdateRole returned error: %v", err)
	}

	want := TeamUser{
		WithUserID: WithUserID{
			UserID: 420,
		},
	}

	if !reflect.DeepEqual(r.TeamUser, want) {
		t.Errorf("TeamUsers.UpdateRole returned %+v, want %+v", r.TeamUser, want)
	}
}
