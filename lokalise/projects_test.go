package lokalise_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/17media/go-lokalise-api/lokalise"
	"github.com/17media/go-lokalise-api/model"
	"github.com/stretchr/testify/assert"
)

func TestClient_Projects_List(t *testing.T) {
	type input struct {
		options lokalise.ProjectsOptions
	}
	type output struct {
		expectedOutgoingRequest outgoingRequest
		response                model.ProjectsResponse
		err                     error
	}
	tt := []struct {
		name           string
		input          input
		serverResponse serverListResponse
		output         output
	}{
		{
			name:  "no pagination succesful json response",
			input: input{},
			serverResponse: serverListResponse{
				statusCode: http.StatusOK,
				body: `{
					"projects": [{
						"project_id": "20008339586cded200e0d8.29879849",
						"name": "SuperApp (iOS + Android)",
						"description": "",
						"created_at": "2018-12-31 12:00:00",
						"created_by": 420,
						"created_by_email": "user@mycompany.com",
						"team_id": 12345
					}]
				}`,
			},
			output: output{
				expectedOutgoingRequest: outgoingRequest{
					path:   "/projects",
					method: http.MethodGet,
				},
				response: model.ProjectsResponse{
					Paged: model.Paged{
						TotalCount: 1,
						PageCount:  2,
						Limit:      3,
						Page:       4,
					},
					Projects: []model.Project{
						{
							ProjectID:      "20008339586cded200e0d8.29879849",
							Name:           "SuperApp (iOS + Android)",
							Description:    "",
							CreatedAt:      "2018-12-31 12:00:00",
							CreatedBy:      420,
							CreatedByEmail: "user@mycompany.com",
							TeamID:         12345,
						},
					},
				},
				err: nil,
			},
		},
		{
			name: "with options succesful json response",
			input: input{
				options: lokalise.ProjectsOptions{
					TeamID: 3,
					PageOptions: lokalise.PageOptions{
						Limit: 1,
						Page:  2,
					},
				},
			},
			serverResponse: serverListResponse{
				statusCode:            http.StatusOK,
				pagedTotalCountHeader: "1",
				pagedPageCountHeader:  "2",
				pagedLimitHeader:      "3",
				pagedPageHeader:       "4",
				body: `{
					"projects": [{
						"project_id": "20008339586cded200e0d8.29879849",
						"name": "SuperApp (iOS + Android)",
						"description": "",
						"created_at": "2018-12-31 12:00:00",
						"created_by": 420,
						"created_by_email": "user@mycompany.com",
						"team_id": 12345
					}]
				}`,
			},
			output: output{
				expectedOutgoingRequest: outgoingRequest{
					path:   "/projects?filter_team_id=3&limit=1&page=2",
					method: http.MethodGet,
				},
				response: model.ProjectsResponse{
					Paged: model.Paged{
						TotalCount: 1,
						PageCount:  2,
						Limit:      3,
						Page:       4,
					},
					Projects: []model.Project{
						{
							ProjectID:      "20008339586cded200e0d8.29879849",
							Name:           "SuperApp (iOS + Android)",
							Description:    "",
							CreatedAt:      "2018-12-31 12:00:00",
							CreatedBy:      420,
							CreatedByEmail: "user@mycompany.com",
							TeamID:         12345,
						},
					},
				},
				err: nil,
			},
		},
		{
			name:  "404 error response",
			input: input{},
			serverResponse: serverListResponse{
				statusCode: http.StatusNotFound,
				body:       notFoundResponseBody("project not found"),
			},
			output: output{
				expectedOutgoingRequest: outgoingRequest{
					path:   "/projects",
					method: http.MethodGet,
				},
				response: model.ProjectsResponse{
					Paged: model.Paged{
						TotalCount: -1,
						PageCount:  -1,
						Limit:      -1,
						Page:       -1,
					},
				},
				err: &model.Error{
					Code:    404,
					Message: "project not found",
				},
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)
			client, fixture, close := setupListClient(t, tc.serverResponse)
			defer close()

			resp, err := client.Projects.List(context.Background(), tc.input.options)

			if tc.output.err != nil {
				assert.EqualError(err, tc.output.err.Error(), "output error not as expected")
			} else {
				assert.NoError(err, "output error not expected")
			}
			tc.output.expectedOutgoingRequest.Assert(t, fixture)
			assert.Equal(tc.output.response.Projects, resp.Projects, "response not as expected")
		})
	}
}

func TestClient_Projects_Create(t *testing.T) {
	inputName := "TheApp Project"
	inputDescription := "iOS + Android strings of TheApp. https://theapp.com"
	mockedServerResponseBody := `{
		"project_id": "3002780358964f9bab5a92.87762498",
		"name": "TheApp Project",
		"description": "iOS + Android strings of TheApp. https://theapp.com",
		"created_at": "2018-12-31 12:00:00",
		"created_by": 420,
		"created_by_email": "user@mycompany.com",
		"team_id": 12345
	}`
	expectedOutgoingRequest := outgoingRequest{
		path:   "/projects",
		method: http.MethodPost,
		body:   `{"description":"iOS + Android strings of TheApp. https://theapp.com","name":"TheApp Project"}`,
	}
	expectedResult := model.Project{
		ProjectID:      "3002780358964f9bab5a92.87762498",
		Name:           "TheApp Project",
		Description:    "iOS + Android strings of TheApp. https://theapp.com",
		CreatedAt:      "2018-12-31 12:00:00",
		CreatedBy:      420,
		CreatedByEmail: "user@mycompany.com",
		TeamID:         12345,
	}
	assert := assert.New(t)
	client, fixture, close := setupClient(t, mockedServerResponseBody)
	defer close()

	resp, err := client.Projects.Create(context.Background(), inputName, inputDescription)

	assert.NoError(err, "unexpeted output error")
	expectedOutgoingRequest.Assert(t, fixture)
	assert.Equal(expectedResult.ProjectID, resp.ProjectID, "response project id not as expected")
	assert.Equal(expectedResult.Name, resp.Name, "response name not as expected")
	assert.Equal(expectedResult.Description, resp.Description, "response description not as expected")
	assert.Equal(expectedResult.CreatedAt, resp.CreatedAt, "response created at not as expected")
	assert.Equal(expectedResult.CreatedBy, resp.CreatedBy, "response created by not as expected")
	assert.Equal(expectedResult.CreatedByEmail, resp.CreatedByEmail, "response created by email not as expected")
	assert.Equal(expectedResult.TeamID, resp.TeamID, "response team id not as expected")
}

func TestClient_Projects_CreateForTeam(t *testing.T) {
	inputName := "TheApp Project"
	inputDescription := "iOS + Android strings of TheApp. https://theapp.com"
	inputTeamID := int64(12345)
	mockedServerResponseBody := `{
		"project_id": "3002780358964f9bab5a92.87762498",
		"name": "TheApp Project",
		"description": "iOS + Android strings of TheApp. https://theapp.com",
		"created_at": "2018-12-31 12:00:00",
		"created_by": 420,
		"created_by_email": "user@mycompany.com",
		"team_id": 12345
	}`
	expectedOutgoingRequest := outgoingRequest{
		path:   "/projects",
		method: http.MethodPost,
		body:   `{"description":"iOS + Android strings of TheApp. https://theapp.com","name":"TheApp Project","team_id":12345}`,
	}
	expectedResult := model.Project{
		ProjectID:      "3002780358964f9bab5a92.87762498",
		Name:           "TheApp Project",
		Description:    "iOS + Android strings of TheApp. https://theapp.com",
		CreatedAt:      "2018-12-31 12:00:00",
		CreatedBy:      420,
		CreatedByEmail: "user@mycompany.com",
		TeamID:         12345,
	}
	assert := assert.New(t)
	client, fixture, close := setupClient(t, mockedServerResponseBody)
	defer close()

	resp, err := client.Projects.CreateForTeam(context.Background(), inputName, inputDescription, inputTeamID)

	assert.NoError(err, "output error not expected")
	expectedOutgoingRequest.Assert(t, fixture)
	assert.Equal(expectedResult.ProjectID, resp.ProjectID, "response project id not as expected")
	assert.Equal(expectedResult.Name, resp.Name, "response name not as expected")
	assert.Equal(expectedResult.Description, resp.Description, "response description not as expected")
	assert.Equal(expectedResult.CreatedAt, resp.CreatedAt, "response created at not as expected")
	assert.Equal(expectedResult.CreatedBy, resp.CreatedBy, "response created by not as expected")
	assert.Equal(expectedResult.CreatedByEmail, resp.CreatedByEmail, "response created by email not as expected")
	assert.Equal(expectedResult.TeamID, resp.TeamID, "response team id not as expected")
}

func TestClient_Projects_Retrieve(t *testing.T) {
	inputProjectID := "1"
	mockedServerResponseBody := `{
		"project_id": "3002780358964f9bab5a92.87762498",
		"name": "TheApp Project",
		"description": "iOS + Android strings of TheApp. https://theapp.com",
		"created_at": "2018-12-31 12:00:00",
		"created_by": 420,
		"created_by_email": "user@mycompany.com",
		"team_id": 12345
	}`
	expectedOutgoingRequest := outgoingRequest{
		path:   "/projects/1",
		method: http.MethodGet,
	}
	expectedResult := model.Project{
		ProjectID:      "3002780358964f9bab5a92.87762498",
		Name:           "TheApp Project",
		Description:    "iOS + Android strings of TheApp. https://theapp.com",
		CreatedAt:      "2018-12-31 12:00:00",
		CreatedBy:      420,
		CreatedByEmail: "user@mycompany.com",
		TeamID:         12345,
	}
	assert := assert.New(t)
	client, fixture, close := setupClient(t, mockedServerResponseBody)
	defer close()

	resp, err := client.Projects.Retrieve(context.Background(), inputProjectID)

	assert.NoError(err, "output error not expected")
	expectedOutgoingRequest.Assert(t, fixture)
	assert.Equal(expectedResult.TeamID, resp.TeamID, "response team id not as expected")
	assert.Equal(expectedResult.Name, resp.Name, "response name not as expected")
}

func TestClient_Projects_Update(t *testing.T) {
	inputProjectID := "3002780358964f9bab5a92"
	inputName := "TheZapp Project"
	inputDescription := "iOS + Android strings of TheZapp. https://thezapp.com"
	mockedServerResponseBody := `{
		"project_id": "3002780358964f9bab5a92.87762498",
		"name": "TheApp Project",
		"description": "iOS + Android strings of TheApp. https://theapp.com",
		"created_at": "2018-12-31 12:00:00 (Etc/UTC)",
		"created_by": 420,
		"created_by_email": "user@mycompany.com",
		"team_id": 12345
	}`
	expectedOutgoingRequest := outgoingRequest{
		path:   "/projects/3002780358964f9bab5a92",
		method: http.MethodPut,
		body:   `{"description":"iOS + Android strings of TheZapp. https://thezapp.com","name":"TheZapp Project"}`,
	}
	expectedResult := model.Project{
		ProjectID:      "3002780358964f9bab5a92.87762498",
		Name:           "TheApp Project",
		Description:    "iOS + Android strings of TheApp. https://theapp.com",
		CreatedAt:      "2018-12-31 12:00:00 (Etc/UTC)",
		CreatedBy:      420,
		CreatedByEmail: "user@mycompany.com",
		TeamID:         12345,
	}
	assert := assert.New(t)
	client, fixture, close := setupClient(t, mockedServerResponseBody)
	defer close()

	resp, err := client.Projects.Update(context.Background(), inputProjectID, inputName, inputDescription)

	assert.NoError(err, "output error not expected")
	expectedOutgoingRequest.Assert(t, fixture)
	assert.Equal(expectedResult.ProjectID, resp.ProjectID, "response project id not as expected")
	assert.Equal(expectedResult.Name, resp.Name, "response name not as expected")
	assert.Equal(expectedResult.Description, resp.Description, "response description not as expected")
	assert.Equal(expectedResult.CreatedAt, resp.CreatedAt, "response created at not as expected")
	assert.Equal(expectedResult.CreatedBy, resp.CreatedBy, "response created by not as expected")
	assert.Equal(expectedResult.CreatedByEmail, resp.CreatedByEmail, "response created by email not as expected")
	assert.Equal(expectedResult.TeamID, resp.TeamID, "response team id not as expected")
}

func TestClient_Projects_Empty(t *testing.T) {
	inputProjectID := "3002780358964f9bab5a92.87762498"
	mockedServerResponseBody := `{
		"project_id": "3002780358964f9bab5a92.87762498",
		"keys_deleted": true
	}`
	expectedOutgoingRequest := outgoingRequest{
		path:   "/projects/3002780358964f9bab5a92.87762498/empty",
		method: http.MethodPut,
	}
	expectedResult := model.ProjectEmptyResponse{
		ProjectID:   "3002780358964f9bab5a92.87762498",
		KeysDeleted: true,
	}
	assert := assert.New(t)
	client, fixture, close := setupClient(t, mockedServerResponseBody)
	defer close()

	resp, err := client.Projects.Empty(context.Background(), inputProjectID)

	assert.NoError(err, "output error not expected")
	expectedOutgoingRequest.Assert(t, fixture)
	assert.Equal(expectedResult.ProjectID, resp.ProjectID, "response project id not as expected")
	assert.Equal(expectedResult.KeysDeleted, resp.KeysDeleted, "response keys deleted not as expected")
}

func TestClient_Projects_Delete(t *testing.T) {
	inputProjectID := "3002780358964f9bab5a92.87762498"
	mockedServerResponseBody := `{
		"project_id": "3002780358964f9bab5a92.87762498",
		"project_deleted": true
	}`
	expectedOutgoingRequest :=
		outgoingRequest{
			path:   "/projects/3002780358964f9bab5a92.87762498",
			method: http.MethodDelete,
		}
	expectedResult := model.ProjectDeleteResponse{
		ProjectID: "3002780358964f9bab5a92.87762498",
		Deleted:   true,
	}
	assert := assert.New(t)
	client, fixture, close := setupClient(t, mockedServerResponseBody)
	defer close()

	resp, err := client.Projects.Delete(context.Background(), inputProjectID)

	assert.NoError(err, "output error not expected")
	expectedOutgoingRequest.Assert(t, fixture)
	assert.Equal(expectedResult.ProjectID, resp.ProjectID, "response project id not as expected")
	assert.Equal(expectedResult.Deleted, resp.Deleted, "response project deleted not as expected")
}
