package lokalise

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestBranchService_Create(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	branchName := "hotfix/really-important"

	mux.HandleFunc(fmt.Sprintf("/projects/%s/branches", testProjectID), func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		testMethod(t, r, "POST")
		testHeader(t, r, apiTokenHeader, testApiToken)
		data := `{
			"name": "` + branchName + `"
		}`

		req := new(bytes.Buffer)
		_ = json.Compact(req, []byte(data))

		testBody(t, r, req.String())

		body := `{
			"project_id": "` + testProjectID + `",
			"branch": {
				"branch_id": 995991,
				"name": "` + branchName + `",
				"created_at": "2019-10-03 14:15:50 (Etc/UTC)",
				"created_at_timestamp": 1567001750,
				"created_by": 1123,
				"created_by_email": "john@lokalise.com"
			}
		}`
		_, _ = fmt.Fprint(w, string(body))
	})

	r, err := client.Branches().Create(testProjectID, "hotfix/really-important")
	if err != nil {
		t.Errorf("Branches.Create returned error: %v", err)
	}

	want := Branch{
		WithCreationTime: WithCreationTime{
			CreatedAt:   "2019-10-03 14:15:50 (Etc/UTC)",
			CreatedAtTs: 1567001750,
		},
		WithCreationUser: WithCreationUser{
			CreatedBy:      1123,
			CreatedByEmail: "john@lokalise.com",
		},
		BranchID: 995991,
		Name:     branchName,
	}

	if !reflect.DeepEqual(r.Branch, want) {
		t.Errorf("Branches.Create returned %+v, want %+v", r.Branch, want)
	}
}

func TestBranchService_List(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/projects/%s/branches", testProjectID), func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		testMethod(t, r, "GET")
		testHeader(t, r, apiTokenHeader, testApiToken)
		body := `{
			"project_id": "` + testProjectID + `",
			"branches": [
				{
					"branch_id": 1234
				}
			]
		}`
		_, _ = fmt.Fprint(w, string(body))
	})

	r, err := client.Branches().List(testProjectID)
	if err != nil {
		t.Errorf("Branches.List returned error: %v", err)
	}

	want := []Branch{
		{
			BranchID: 1234,
		},
	}

	if !reflect.DeepEqual(r.Branches, want) {
		t.Errorf("Branches.List returned %+v, want %+v", r.Branches, want)
	}
}

func TestBranchService_Delete(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/projects/%s/branches/1", testProjectID), func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		testMethod(t, r, "DELETE")
		testHeader(t, r, apiTokenHeader, testApiToken)
		body := `{
			"project_id": "` + testProjectID + `",
			"branch_deleted": true
		}`
		_, _ = fmt.Fprint(w, string(body))
	})

	r, err := client.Branches().Delete(testProjectID, 1)
	if err != nil {
		t.Errorf("Branches.Delete returned error: %v", err)
	}

	want := DeleteBranchResponse{
		WithProjectID: WithProjectID{ProjectID: testProjectID},
		BranchDeleted: true,
	}

	if !reflect.DeepEqual(r, want) {
		t.Errorf("Branches.Delete returned %+v, want %+v", r, want)
	}
}
