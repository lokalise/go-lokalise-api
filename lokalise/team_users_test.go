package lokalise_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/17media/go-lokalise-api/lokalise"
	"github.com/17media/go-lokalise-api/model"
	"github.com/stretchr/testify/assert"
)

func TestClient_TeamUsers_List(t *testing.T) {
	type input struct {
		teamID  int64
		options lokalise.PageOptions
	}
	type output struct {
		expectedOutgoingRequest outgoingRequest
		response                model.TeamUsersResponse
		err                     error
	}
	tt := []struct {
		name           string
		input          input
		serverResponse serverListResponse
		output         output
	}{
		{
			name: "no pagination succesful json response",
			input: input{
				teamID: 1,
			},
			serverResponse: serverListResponse{
				statusCode: http.StatusOK,
				body: `{
					"team_id": 18821,
					"team_users": [{
						"user_id": 420,
						"email": "jdoe@mycompany.com"
					}]
				}`,
			},
			output: output{
				expectedOutgoingRequest: outgoingRequest{
					path:   "/teams/1/users",
					method: http.MethodGet,
				},
				response: model.TeamUsersResponse{
					TeamID: 18821,
					Paged: model.Paged{
						TotalCount: -1,
						PageCount:  -1,
						Limit:      -1,
						Page:       -1,
					},
					TeamUsers: []model.TeamUser{
						{
							UserID: 420,
							Email:  "jdoe@mycompany.com",
						},
					},
				},
				err: nil,
			},
		},
		{
			name: "with pagination succesful json response",
			input: input{
				teamID: 1,
				options: lokalise.PageOptions{
					Limit: 1,
					Page:  2,
				},
			},
			serverResponse: serverListResponse{
				statusCode:            http.StatusOK,
				pagedTotalCountHeader: "1",
				pagedPageCountHeader:  "2",
				pagedLimitHeader:      "3",
				pagedPageHeader:       "4",
				body: `{
					"team_id": 18821,
					"team_users": [{
						"user_id": 420,
						"email": "jdoe@mycompany.com"
					}]
				}`,
			},
			output: output{
				expectedOutgoingRequest: outgoingRequest{
					path:   "/teams/1/users?limit=1&page=2",
					method: http.MethodGet,
				},
				response: model.TeamUsersResponse{
					Paged: model.Paged{
						TotalCount: 1,
						PageCount:  2,
						Limit:      3,
						Page:       4,
					},
					TeamID: 18821,
					TeamUsers: []model.TeamUser{
						{
							UserID: 420,
							Email:  "jdoe@mycompany.com",
						},
					},
				},
				err: nil,
			},
		},
		{
			name: "invalid integers in pagination headers",
			input: input{
				teamID: 1,
			},
			serverResponse: serverListResponse{
				statusCode:            http.StatusOK,
				pagedTotalCountHeader: "a",
				pagedPageCountHeader:  "b",
				pagedLimitHeader:      "c",
				pagedPageHeader:       "d",
				body: `{
					"team_id": 18821,
					"team_users": [{
						"user_id": 420,
						"email": "jdoe@mycompany.com"
					}]
				}`,
			},
			output: output{
				expectedOutgoingRequest: outgoingRequest{
					path:   "/teams/1/users",
					method: http.MethodGet,
				},
				response: model.TeamUsersResponse{
					Paged: model.Paged{
						TotalCount: -1,
						PageCount:  -1,
						Limit:      -1,
						Page:       -1,
					},
					TeamID: 18821,
					TeamUsers: []model.TeamUser{
						{
							UserID: 420,
							Email:  "jdoe@mycompany.com",
						},
					},
				},
				err: nil,
			},
		},
		{
			name: "404 error response",
			input: input{
				teamID: 1,
			},
			serverResponse: serverListResponse{
				statusCode: http.StatusNotFound,
				body:       notFoundResponseBody("team not found"),
			},
			output: output{
				expectedOutgoingRequest: outgoingRequest{
					path:   "/teams/1/users",
					method: http.MethodGet,
				},
				response: model.TeamUsersResponse{
					Paged: model.Paged{
						TotalCount: -1,
						PageCount:  -1,
						Limit:      -1,
						Page:       -1,
					},
				},
				err: &model.Error{
					Code:    404,
					Message: "team not found",
				},
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)
			client, fixture, close := setupListClient(t, tc.serverResponse)
			defer close()

			resp, err := client.TeamUsers.List(context.Background(), tc.input.teamID, tc.input.options)

			if tc.output.err != nil {
				assert.EqualError(err, tc.output.err.Error(), "output error not as expected")
			} else {
				assert.NoError(err, "output error not expected")
			}
			tc.output.expectedOutgoingRequest.Assert(t, fixture)
			assert.Equal(tc.output.response.TeamID, resp.TeamID, "response team id not as expected")
			assert.Equal(tc.output.response.Paged, resp.Paged, "paged response values not as expected")
			assert.Equal(tc.output.response.TeamUsers, resp.TeamUsers, "response team user id not as expected")
		})
	}
}

func TestClient_TeamUsers_Retrieve(t *testing.T) {
	inputTeamID := int64(1)
	inputUserID := int64(2)
	mockedServerResponseBody := `{
		"team_id": 18821,
		"team_user": {
			"user_id": 420,
			"email": "jdoe@mycompany.com",
			"fullname": "John Doe",
			"created_at": "2018-12-31 12:00:00",
			"role": "owner"
		}
	}`
	expectedOutgoingRequest := outgoingRequest{
		path:   "/teams/1/users/2",
		method: http.MethodGet,
	}
	expectedResult := model.TeamUserResponse{
		TeamID: 18821,
		TeamUser: model.TeamUser{
			UserID: 420,
			Email:  "jdoe@mycompany.com",
		},
	}
	assert := assert.New(t)
	client, fixture, close := setupClient(t, mockedServerResponseBody)
	defer close()

	resp, err := client.TeamUsers.Retrieve(context.Background(), inputTeamID, inputUserID)

	assert.NoError(err, "output error not expected")
	expectedOutgoingRequest.Assert(t, fixture)
	assert.Equal(expectedResult.TeamID, resp.TeamID, "response team id not as expected")
	assert.Equal(expectedResult.TeamUser.UserID, resp.TeamUser.UserID, "response team user id not as expected")
}

func TestClient_TeamUsers_UpdateRole(t *testing.T) {
	inputTeamID := int64(1)
	inputUserID := int64(2)
	inputRole := model.TeamUserRoleAdmin
	mockedServerResponseBody := `{
		"team_id": 18821,
		"team_user": {
			"user_id": 420,
			"email": "jdoe@mycompany.com",
			"fullname": "John Doe",
			"created_at": "2018-12-31 12:00:00",
			"role": "owner"
		}
	}`
	expectedOutgoingRequest := outgoingRequest{
		method: http.MethodPut,
		path:   "/teams/1/users/2",
		body:   `{"role":"admin"}`,
	}
	expectedResult := model.TeamUserResponse{
		TeamID: 18821,
		TeamUser: model.TeamUser{
			UserID: 420,
			Email:  "jdoe@mycompany.com",
		},
	}
	assert := assert.New(t)
	client, fixture, close := setupClient(t, mockedServerResponseBody)
	defer close()

	resp, err := client.TeamUsers.UpdateRole(context.Background(), inputTeamID, inputUserID, inputRole)

	assert.NoError(err, "output error not expected")
	expectedOutgoingRequest.Assert(t, fixture)
	assert.Equal(expectedResult.TeamID, resp.TeamID, "response team id not as expected")
	assert.Equal(expectedResult.TeamUser.UserID, resp.TeamUser.UserID, "response team user id not as expected")
}

func TestClient_TeamUsers_Delete(t *testing.T) {
	inputTeamID := int64(1)
	inputUserID := int64(2)
	mockedServerResponseBody := `{
		"team_id": 18821,
		"team_user_deleted": true
	}`
	expectedOutgoingRequest := outgoingRequest{
		path:   "/teams/1/users/2",
		method: http.MethodDelete,
	}
	expectedResult := model.TeamUserDeleteResponse{
		TeamID:  18821,
		Deleted: true,
	}
	assert := assert.New(t)
	client, fixture, close := setupClient(t, mockedServerResponseBody)
	defer close()

	resp, err := client.TeamUsers.Delete(context.Background(), inputTeamID, inputUserID)

	assert.NoError(err, "output error not expected")
	expectedOutgoingRequest.Assert(t, fixture)
	assert.Equal(expectedResult.TeamID, resp.TeamID, "response team id not as expected")
	assert.Equal(expectedResult.Deleted, resp.Deleted, "response team user id not as expected")
}

type testLogger struct {
	T *testing.T
}

func (l *testLogger) Write(p []byte) (n int, err error) {
	l.T.Log(string(p))
	return len(p), nil
}
