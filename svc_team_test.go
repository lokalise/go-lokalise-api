package lokalise

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestTeamService_List(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(
		"/teams",
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			testMethod(t, r, "GET")
			testHeader(t, r, apiTokenHeader, testApiToken)

			_, _ = fmt.Fprint(w, `{
				"teams": [
					{
						"team_id": 18821,
						"name": "MyCompany, Ltd.",
						"created_at": "2018-12-31 12:00:00 (Etc/UTC)",
						"created_at_timestamp": 1546257600,
						"plan": "Essential",
						"quota_usage": {
							"users": 14,
							"keys": 8125,
							"projects": 4,
							"mau": 119337
						},
						"quota_allowed": {
							"users": 40,
							"keys": 10000,
							"projects": 99999999,
							"mau": 200000
						}
					}
				]
			}`)
		})

	r, err := client.Teams().List()
	if err != nil {
		t.Errorf("Teams.List returned error: %v", err)
	}

	want := []Team{
		{
			WithCreationTime: WithCreationTime{
				CreatedAt:   "2018-12-31 12:00:00 (Etc/UTC)",
				CreatedAtTs: 1546257600,
			},
			WithTeamID: WithTeamID{
				TeamID: 18821,
			},
			Name: "MyCompany, Ltd.",
			Plan: "Essential",
			QuotaUsage: Quota{
				Users:    14,
				Keys:     8125,
				Projects: 4,
				MAU:      119337,
			},
			QuotaAllowed: Quota{
				Users:    40,
				Keys:     10000,
				Projects: 99999999,
				MAU:      200000,
			},
		},
	}

	if !reflect.DeepEqual(r.Teams, want) {
		t.Errorf("Screenshots.List returned %+v, want %+v", r.Teams, want)
	}
}
