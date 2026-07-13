package lokalise

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestQueuedProcessService_List(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(
		fmt.Sprintf("/projects/%s/processes", testProjectID),
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			testMethod(t, r, "GET")
			testHeader(t, r, apiTokenHeader, testApiToken)

			_, _ = fmt.Fprint(w, `{
				"project_id": "`+testProjectID+`",
				"processes": [
					{
						"process_id": "2e0559e60e856555fbc15bdf78ab2b0ca3406e8f",
						"type": "file-import",
						"status": "finished",
						"message": "",
						"created_by": 1234,
						"created_by_email": "example@example.com",
						"created_at": "2020-04-20 13:43:43 (Etc/UTC)",
						"created_at_timestamp": 1587390223
					}
				]
			}`)
		})

	r, err := client.QueuedProcesses().List(testProjectID)
	if err != nil {
		t.Errorf("QueuedProcesses.List returned error: %v", err)
	}

	want := []QueuedProcess{
		{
			ID:      "2e0559e60e856555fbc15bdf78ab2b0ca3406e8f",
			Type:    "file-import",
			Status:  "finished",
			Message: "",
			WithCreationUser: WithCreationUser{
				CreatedBy:      1234,
				CreatedByEmail: "example@example.com",
			},
			WithCreationTime: WithCreationTime{
				CreatedAt:   "2020-04-20 13:43:43 (Etc/UTC)",
				CreatedAtTs: 1587390223,
			},
		},
	}

	if !reflect.DeepEqual(r.Processes, want) {
		t.Errorf("QueuedProcesses.List returned %+v, want %+v", r.Processes, want)
	}
}

func TestQueuedProcessService_Retrieve(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	processId := "2e0559e60e856555fbc15bdf78ab2b0ca3406e8f"

	mux.HandleFunc(
		fmt.Sprintf("/projects/%s/processes/%s", testProjectID, processId),
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			testMethod(t, r, "GET")
			testHeader(t, r, apiTokenHeader, testApiToken)

			_, _ = fmt.Fprint(w, `{
				"project_id": "`+testProjectID+`",
				"process": {
					"process_id": "`+processId+`",
					"type": "file-import",
					"status": "finished",
					"message": "",
					"created_by": 1234,
					"created_by_email": "example@example.com",
					"created_at": "2020-04-20 13:43:43 (Etc/UTC)",
					"created_at_timestamp": 1587390223,
					"details": {
						"files": [
							{
								"status": "finished",
								"message": "",
								"name_original": "index.json",
								"name_custom": "index.json",
								"word_count_total": 2,
								"key_count_total": 1,
								"key_count_inserted": 0,
								"key_count_updated": 0,
								"key_count_skipped": 1
							}
						]
					}
				}
			}`)
		})

	r, err := client.QueuedProcesses().Retrieve(testProjectID, processId)
	if err != nil {
		t.Errorf("QueuedProcesses.Retrieve returned error: %v", err)
	}

	want := QueuedProcess{
		ID:      processId,
		Type:    "file-import",
		Status:  "finished",
		Message: "",
		Details: ProcessDetails{
			Files: []ProcessDetailsFile{
				{
					Status:           "finished",
					Message:          "",
					NameOriginal:     "index.json",
					NameCustom:       "index.json",
					WordCountTotal:   2,
					KeyCountTotal:    1,
					KeyCountInserted: 0,
					KeyCountUpdated:  0,
					KeyCountSkipped:  1,
				},
			},
		},
		WithCreationUser: WithCreationUser{
			CreatedBy:      1234,
			CreatedByEmail: "example@example.com",
		},
		WithCreationTime: WithCreationTime{
			CreatedAt:   "2020-04-20 13:43:43 (Etc/UTC)",
			CreatedAtTs: 1587390223,
		},
	}

	if !reflect.DeepEqual(r.Process, want) {
		t.Errorf("QueuedProcesses.Retrieve returned %+v, want %+v", r.Process, want)
	}
}
