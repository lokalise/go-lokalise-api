package lokalise

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestTeamPermissionTemplates_List(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(
		"/teams/1/roles",
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			testMethod(t, r, "GET")
			testHeader(t, r, apiTokenHeader, testApiToken)

			_, _ = fmt.Fprint(w, `{
				"roles": [
					{
						"id": 1,
						"role": "Localisation management",
						"permissions": [
							"activity",
							"branches_main_modify"
						],
						"description": "Manage project settings, contributors and tasks",
						"tag": "Full access",
						"tagColor": "green",
						"doesEnableAllReadOnlyLanguages": true
					},
					{
						"id": 2,
						"role": "Developer",
						"permissions": [
							"download",
							"upload"
						],
						"description": "Create keys, upload and download content",
						"tag": "Advanced",
						"tagColor": "cyan",
						"doesEnableAllReadOnlyLanguages": true
					}
				]
			}`)
		})

	r, err := client.PermissionTemplates().ListPermissionRoles(1)
	if err != nil {
		t.Errorf("Teams.List returned error: %v", err)
	}

	want := PermissionRoleResponse{
		Roles: []PermissionTemplate{
			{
				ID:   1,
				Role: "Localisation management",
				Permissions: []string{
					"activity",
					"branches_main_modify",
				},
				Description:                    "Manage project settings, contributors and tasks",
				Tag:                            "Full access",
				TagColor:                       "green",
				DoesEnableAllReadOnlyLanguages: true,
			},
			{
				ID:   2,
				Role: "Developer",
				Permissions: []string{
					"download",
					"upload",
				},
				Description:                    "Create keys, upload and download content",
				Tag:                            "Advanced",
				TagColor:                       "cyan",
				DoesEnableAllReadOnlyLanguages: true,
			},
		},
	}

	if !reflect.DeepEqual(r, want) {
		t.Errorf("Team.PermissionRoles.List returned %+v, want %+v", r, want)
	}
}
