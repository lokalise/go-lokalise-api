package lokalise_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/lokalise/go-lokalise-api"
	"github.com/stretchr/testify/assert"
)

func TestClient_Snapshots_List(t *testing.T) {
	t.Run("no pagination", func(t *testing.T) {
		assert := assert.New(t)
		const projectID = "20008339586cded200e0d8.29879849"
		client, fixture, stop := setupListClient(t, serverListResponse{
			statusCode: http.StatusOK,
			body: ` {
				    "project_id": "20008339586cded200e0d8.29879849",
				    "snapshots": [
					{
					    "snapshot_id": 1523966589,
					    "title": "API snapshot",
					    "created_at": "2018-12-31 12:00:00 (Etc\/UTC)",
					    "created_at_timestamp": 1546257600,
					    "created_by": 420,
					    "created_by_email": "user@mycompany.com"
					},
					{
					    "snapshot_id": 1523966599,
					    "title": "API snapshot",
					    "created_at": "2018-12-31 12:00:00 (Etc\/UTC)",
					    "created_at_timestamp": 1546257600,
					    "created_by": 420,
					    "created_by_email": "user@mycompany.com"
					}
				    ]
				}`,
		})
		defer stop()
		expectedRequest := &outgoingRequest{
			path:   "/projects/20008339586cded200e0d8.29879849/snapshots",
			method: http.MethodGet,
		}
		expectedResponse := lokalise.ListSnapshotsResponse{
			Paged: lokalise.Paged{
				TotalCount: -1,
				PageCount:  -1,
				Limit:      -1,
				Page:       -1,
			},
			ProjectID: projectID,
			Snapshots: []lokalise.Snapshot{{
				SnapshotID:     1523966589,
				Title:          "API snapshot",
				CreatedAt:      "2018-12-31 12:00:00 (Etc/UTC)",
				CreatedAtTs:    1546257600,
				CreatedBy:      420,
				CreatedByEmail: "user@mycompany.com",
			}, {
				SnapshotID:     1523966599,
				Title:          "API snapshot",
				CreatedAt:      "2018-12-31 12:00:00 (Etc/UTC)",
				CreatedAtTs:    1546257600,
				CreatedBy:      420,
				CreatedByEmail: "user@mycompany.com",
			}},
		}

		resp, err := client.Snapshots.List(context.Background(), projectID, lokalise.SnapshotsOptions{})
		assert.NoError(err)
		expectedRequest.Assert(t, fixture)
		assert.Equal(expectedResponse, resp)
	})
	t.Run("with pagination", func(t *testing.T) {
		assert := assert.New(t)
		const projectID = "20008339586cded200e0d8.29879849"
		client, fixture, stop := setupListClient(t, serverListResponse{
			statusCode:            http.StatusOK,
			pagedTotalCountHeader: "2",
			pagedPageCountHeader:  "2",
			pagedLimitHeader:      "1",
			pagedPageHeader:       "2",
			body: ` {
				    "project_id": "20008339586cded200e0d8.29879849",
				    "snapshots": [
					{
					    "snapshot_id": 1523966589,
					    "title": "API snapshot",
					    "created_at": "2018-12-31 12:00:00 (Etc\/UTC)",
					    "created_at_timestamp": 1546257600,
					    "created_by": 420,
					    "created_by_email": "user@mycompany.com"
					},
					{
					    "snapshot_id": 1523966599,
					    "title": "API snapshot",
					    "created_at": "2018-12-31 12:00:00 (Etc\/UTC)",
					    "created_at_timestamp": 1546257600,
					    "created_by": 420,
					    "created_by_email": "user@mycompany.com"
					}
				    ]
				}`,
		})
		defer stop()
		expectedRequest := &outgoingRequest{
			path:   "/projects/20008339586cded200e0d8.29879849/snapshots?limit=1&page=2",
			method: http.MethodGet,
		}
		expectedResponse := lokalise.ListSnapshotsResponse{
			Paged: lokalise.Paged{
				TotalCount: 2,
				PageCount:  2,
				Limit:      1,
				Page:       2,
			},
			ProjectID: projectID,
			Snapshots: []lokalise.Snapshot{{
				SnapshotID:     1523966589,
				Title:          "API snapshot",
				CreatedAt:      "2018-12-31 12:00:00 (Etc/UTC)",
				CreatedAtTs:    1546257600,
				CreatedBy:      420,
				CreatedByEmail: "user@mycompany.com",
			}, {
				SnapshotID:     1523966599,
				Title:          "API snapshot",
				CreatedAt:      "2018-12-31 12:00:00 (Etc/UTC)",
				CreatedAtTs:    1546257600,
				CreatedBy:      420,
				CreatedByEmail: "user@mycompany.com",
			}},
		}

		resp, err := client.Snapshots.List(context.Background(), projectID, lokalise.SnapshotsOptions{
			PageOptions: lokalise.PageOptions{
				Limit: 1,
				Page:  2,
			},
		})
		assert.NoError(err)
		expectedRequest.Assert(t, fixture)
		assert.Equal(expectedResponse, resp)
	})
	t.Run("http error", func(t *testing.T) {
		assert := assert.New(t)
		const projectID = "20008339586cded200e0d8.29879849"
		client, fixture, stop := setupListClient(t, serverListResponse{
			statusCode: http.StatusNotFound,
			body:       notFoundResponseBody("Snapshots not found"),
		})
		defer stop()
		expectedRequest := &outgoingRequest{
			path:   "/projects/20008339586cded200e0d8.29879849/snapshots",
			method: http.MethodGet,
		}
		expectedError := &lokalise.Error{Code: 404, Message: "Snapshots not found"}

		_, err := client.Snapshots.List(context.Background(), projectID, lokalise.SnapshotsOptions{})
		assert.EqualError(err, expectedError.Error())
		expectedRequest.Assert(t, fixture)
	})
}

func TestClient_Snapshots_Create(t *testing.T) {
	assert := assert.New(t)
	const projectID = "20008339586cded200e0d8.29879849"
	client, fixture, stop := setupClient(t, `
				{
				    "project_id": "20008339586cded200e0d8.29879849",
				    "snapshot": {
					    "snapshot_id": 1523966589,
					    "title": "testname",
					    "created_at": "2018-12-31 12:00:00 (Etc\/UTC)",
					    "created_at_timestamp": 1546257600,
					    "created_by": 420,
					    "created_by_email": "user@mycompany.com"
					}
				}`,
	)
	defer stop()
	expectedRequest := &outgoingRequest{
		path:   "/projects/20008339586cded200e0d8.29879849/snapshots",
		method: http.MethodPost,
		body:   `{"title":"testname"}`,
	}
	expectedResponse := lokalise.CreateSnapshotResponse{
		ProjectID: projectID,
		Snapshot: lokalise.Snapshot{
			SnapshotID:     1523966589,
			Title:          "testname",
			CreatedAt:      "2018-12-31 12:00:00 (Etc/UTC)",
			CreatedAtTs:    1546257600,
			CreatedBy:      420,
			CreatedByEmail: "user@mycompany.com",
		},
	}

	resp, err := client.Snapshots.Create(context.Background(), projectID, "testname")
	assert.NoError(err)
	expectedRequest.Assert(t, fixture)
	assert.Equal(expectedResponse, resp)
}

func TestClient_Snapshots_Delete(t *testing.T) {
	assert := assert.New(t)
	const projectID = "20008339586cded200e0d8.29879849"
	client, fixture, stop := setupClient(t, `
				{
				    "project_id": "20008339586cded200e0d8.29879849",
				    "snapshot_deleted": true
				}`,
	)
	defer stop()
	expectedRequest := &outgoingRequest{
		path:   "/projects/20008339586cded200e0d8.29879849/snapshots/42",
		method: http.MethodDelete,
	}
	expectedResponse := lokalise.DeleteSnapshotResponse{
		ProjectID:       projectID,
		SnapshotDeleted: true,
	}

	resp, err := client.Snapshots.Delete(context.Background(), projectID, 42)
	assert.NoError(err)
	expectedRequest.Assert(t, fixture)
	assert.Equal(expectedResponse, resp)
}

func TestClient_Snapshots_Restore(t *testing.T) {
	assert := assert.New(t)
	const projectID = "20008339586cded200e0d8.29879849"
	client, fixture, stop := setupClient(t, `
	{
	    "project_id": "20008339586cded200e0d8.29879849",
	    "project_type": "localization_files",
	    "name": "TheApp Project Copy",
	    "description": "iOS + Android strings of TheApp. https://theapp.com",
	    "created_at": "2018-12-13 12:00:00 (Etc\/UTC)",
	    "created_at_timestamp": 1546257600,
	    "created_by": 420,
	    "created_by_email": "user@mycompany.com",
	    "team_id": 12345,
	    "base_language_id": 640,
	    "base_language_iso": "en",
	    "settings": {
		"per_platform_key_names": false,
		"reviewing": false,
		"upvoting": false,
		"auto_toggle_unverified": true,
		"offline_translation": false,
		"key_editing": true,
		"inline_machine_translations": true
	    },
	    "statistics": {
		"progress_total": 87,
		"keys_total": 13,
		"team": 2,
		"base_words": 22,
		"qa_issues_total": 65,
		"qa_issues": {
		    "not_reviewed": 39,
		    "unverified": 7,
		    "spelling_grammar": 4,
		    "inconsistent_placeholders": 5,
		    "inconsistent_html": 1,
		    "different_number_of_urls": 0,
		    "different_urls": 1,
		    "leading_whitespace": 1,
		    "trailing_whitespace": 0,
		    "different_number_of_email_address": 0,
		    "different_email_address": 0,
		    "different_brackets": 6,
		    "different_numbers": 1,
		    "double_space": 1,
		    "special_placeholder": 0
		},
		"languages": [
		    {
			"language_id": 640,
			"language_iso": "en",
			"progress": 100,
			"words_to_do": 0
		    },
		    {
			"language_id": 800,
			"language_iso": "lv_LV",
			"progress": 100,
			"words_to_do": 0
		    },
		    {
			"language_id": 673,
			"language_iso": "fr",
			"progress": 63,
			"words_to_do": 8
		    }
		]
	    }
	}`,
	)
	defer stop()
	expectedRequest := &outgoingRequest{
		path:   "/projects/20008339586cded200e0d8.29879849/snapshots/42",
		method: http.MethodPost,
	}
	expectedResponse := lokalise.Project{
		ProjectID:      "20008339586cded200e0d8.29879849",
		Name:           "TheApp Project Copy",
		Description:    "iOS + Android strings of TheApp. https://theapp.com",
		CreatedAt:      "2018-12-13 12:00:00 (Etc/UTC)",
		CreatedBy:      420,
		CreatedByEmail: "user@mycompany.com",
		TeamID:         12345,
	}

	resp, err := client.Snapshots.Restore(context.Background(), projectID, 42)
	assert.NoError(err)
	expectedRequest.Assert(t, fixture)
	assert.Equal(expectedResponse, resp)
}
