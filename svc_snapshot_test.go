package lokalise

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestSnapshotService_Create(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(
		fmt.Sprintf("/projects/%s/snapshots", testProjectID),
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			testMethod(t, r, "POST")
			testHeader(t, r, apiTokenHeader, testApiToken)
			data := `{
				"title": "API snapshot"
			}
			`

			req := new(bytes.Buffer)
			_ = json.Compact(req, []byte(data))

			testBody(t, r, req.String())

			_, _ = fmt.Fprint(w, `{
				"project_id": "`+testProjectID+`",
				"snapshot": {
					"snapshot_id": 1523966599,
					"title": "API snapshot",
					"created_at": "2018-12-31 12:00:00 (Etc/UTC)",
					"created_at_timestamp": 1546257600,
					"created_by": 420,
					"created_by_email": "user@mycompany.com"
				}
			}`)
		})

	r, err := client.Snapshots().Create(testProjectID, "API snapshot")
	if err != nil {
		t.Errorf("Snapshots.Create returned error: %v", err)
	}

	want := Snapshot{
		WithCreationTime: WithCreationTime{
			CreatedAt:   "2018-12-31 12:00:00 (Etc/UTC)",
			CreatedAtTs: 1546257600,
		},
		WithCreationUser: WithCreationUser{
			CreatedBy:      420,
			CreatedByEmail: "user@mycompany.com",
		},
		SnapshotID: 1523966599,
		Title:      "API snapshot",
	}

	if !reflect.DeepEqual(r.Snapshot, want) {
		t.Errorf("Snapshots.Create returned %+v, want %+v", r.Snapshot, want)
	}
}

func TestSnapshotService_Delete(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(
		fmt.Sprintf("/projects/%s/snapshots/%d", testProjectID, 1523966599),
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			testMethod(t, r, "DELETE")
			testHeader(t, r, apiTokenHeader, testApiToken)

			_, _ = fmt.Fprint(w, `{
				"project_id": "`+testProjectID+`",
				"snapshot_deleted": true
			}`)
		},
	)

	r, err := client.Snapshots().Delete(testProjectID, 1523966599)
	if err != nil {
		t.Errorf("Screenshots.Delete returned error: %v", err)
	}

	want := DeleteSnapshotResponse{
		WithProjectID: WithProjectID{
			ProjectID: testProjectID,
		},
		SnapshotDeleted: true,
	}

	if !reflect.DeepEqual(r, want) {
		t.Errorf("Screenshots.Delete returned %+v, want %+v", r, want)
	}
}

func TestSnapshotService_List(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(
		fmt.Sprintf("/projects/%s/snapshots", testProjectID),
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			testMethod(t, r, "GET")
			testHeader(t, r, apiTokenHeader, testApiToken)

			_, _ = fmt.Fprint(w, `{
				"project_id": "`+testProjectID+`",
				"snapshots": [
					{
						"snapshot_id": 1523966589
					}
				]
			}`)
		})

	r, err := client.Snapshots().List(testProjectID)
	if err != nil {
		t.Errorf("Snapshots.List returned error: %v", err)
	}

	want := []Snapshot{
		{
			SnapshotID: 1523966589,
		},
	}

	if !reflect.DeepEqual(r.Snapshots, want) {
		t.Errorf("Screenshots.List returned %+v, want %+v", r.Snapshots, want)
	}
}

func TestSnapshotService_Restore(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(
		fmt.Sprintf("/projects/%s/snapshots/%d", testProjectID, 1523966599),
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			testMethod(t, r, "POST")
			testHeader(t, r, apiTokenHeader, testApiToken)

			_, _ = fmt.Fprint(w, `{
				"project_id": "`+testProjectID+`"
			}`)
		})

	r, err := client.Snapshots().Restore(testProjectID, 1523966599)
	if err != nil {
		t.Errorf("Snapshots.Restore returned error: %v", err)
	}

	want := Project{
		WithProjectID: WithProjectID{
			ProjectID: testProjectID,
		},
	}

	if !reflect.DeepEqual(r, want) {
		t.Errorf("Screenshots.Restore returned %+v, want %+v", r, want)
	}
}
