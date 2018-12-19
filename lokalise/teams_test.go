package lokalise_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/17media/go-lokalise-api/lokalise"
	"github.com/17media/go-lokalise-api/model"
	"github.com/stretchr/testify/assert"
)

func TestClient_Teams_List(t *testing.T) {
	type input struct {
		options lokalise.PageOptions
	}
	type output struct {
		expectedOutgoingRequest outgoingRequest
		response                model.TeamsResponse
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
					"teams": [
						{
							"team_id": 178017,
							"name": "test",
							"plan": "Free",
							"created_at": "2018-10-09 21:08:05 (Etc/UTC)",
							"quota_usage": {
								"users": 1,
								"keys": 0,
								"projects": 0,
								"mau": 0
							},
							"quota_allowed": {
								"users": 999999999,
								"keys": 999999999,
								"projects": 999999999,
								"mau": 999999999
							}
						}
					]
				}`,
			},
			output: output{
				expectedOutgoingRequest: outgoingRequest{
					path:   "/teams",
					method: http.MethodGet,
				},
				response: model.TeamsResponse{
					Paged: model.Paged{
						TotalCount: 1,
						PageCount:  2,
						Limit:      3,
						Page:       4,
					},
					Teams: []model.Team{
						{
							TeamID:    178017,
							Name:      "test",
							Plan:      "Free",
							CreatedAt: "2018-10-09 21:08:05 (Etc/UTC)",
							QuotaUsage: model.Quota{
								Users:    1,
								Keys:     0,
								Projects: 0,
								MAU:      0,
							},
							QuotaAllowed: model.Quota{
								Users:    999999999,
								Keys:     999999999,
								Projects: 999999999,
								MAU:      999999999,
							},
						},
					},
				},
				err: nil,
			},
		},
		{
			name: "with pagination succesful json response",
			input: input{
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
					"teams": [
						{
							"team_id": 178017,
							"name": "test",
							"plan": "Free",
							"created_at": "2018-10-09 21:08:05 (Etc/UTC)",
							"quota_usage": {
								"users": 1,
								"keys": 0,
								"projects": 0,
								"mau": 0
							},
							"quota_allowed": {
								"users": 999999999,
								"keys": 999999999,
								"projects": 999999999,
								"mau": 999999999
							}
						}
					]
				}`,
			},
			output: output{
				expectedOutgoingRequest: outgoingRequest{
					path:   "/teams?limit=1&page=2",
					method: http.MethodGet,
				},
				response: model.TeamsResponse{
					Paged: model.Paged{
						TotalCount: 1,
						PageCount:  2,
						Limit:      3,
						Page:       4,
					},
					Teams: []model.Team{
						{
							TeamID:    178017,
							Name:      "test",
							Plan:      "Free",
							CreatedAt: "2018-10-09 21:08:05 (Etc/UTC)",
							QuotaUsage: model.Quota{
								Users:    1,
								Keys:     0,
								Projects: 0,
								MAU:      0,
							},
							QuotaAllowed: model.Quota{
								Users:    999999999,
								Keys:     999999999,
								Projects: 999999999,
								MAU:      999999999,
							},
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
				body:       notFoundResponseBody("team not found"),
			},
			output: output{
				expectedOutgoingRequest: outgoingRequest{
					path:   "/teams",
					method: http.MethodGet,
				},
				response: model.TeamsResponse{
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

			resp, err := client.Teams.List(context.Background(), tc.input.options)

			if tc.output.err != nil {
				assert.EqualError(err, tc.output.err.Error(), "output error not as expected")
			} else {
				assert.NoError(err, "output error not expected")
			}
			tc.output.expectedOutgoingRequest.Assert(t, fixture)
			assert.Equal(tc.output.response.Teams, resp.Teams, "response not as expected")
		})
	}
}
