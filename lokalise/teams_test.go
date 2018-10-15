package lokalise_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lokalise/lokalise-go-sdk/lokalise"
	"github.com/lokalise/lokalise-go-sdk/model"
	"github.com/stretchr/testify/assert"
)

func TestClient_Teams_List(t *testing.T) {
	type input struct {
		options lokalise.PageOptions
	}
	type output struct {
		calledPath string
		response   model.TeamsResponse
		err        error
	}
	type serverResponse struct {
		body                  string
		statusCode            int
		pagedTotalCountHeader string
		pagedPageCountHeader  string
		pagedLimitHeader      string
		pagedPageHeader       string
	}
	tt := []struct {
		name           string
		input          input
		serverResponse serverResponse
		output         output
	}{
		{
			name:  "no pagination succesful json response",
			input: input{},
			serverResponse: serverResponse{
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
				calledPath: "/teams",
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
			serverResponse: serverResponse{
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
				calledPath: "/teams?limit=1&page=2",
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
			serverResponse: serverResponse{
				statusCode: http.StatusNotFound,
				body: `{
					"error": {
						"code": 404,
						"message": "team not found"
					}
				}`,
			},
			output: output{
				calledPath: "/teams",
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

			var calledPath string
			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				calledPath = req.URL.String()
				rw.Header().Set("Content-Type", "application/json")
				if tc.serverResponse.pagedTotalCountHeader != "" {
					rw.Header().Set("X-Pagination-Total-Count", tc.serverResponse.pagedTotalCountHeader)
				}
				if tc.serverResponse.pagedPageCountHeader != "" {
					rw.Header().Set("X-Pagination-Page-Count", tc.serverResponse.pagedPageCountHeader)
				}
				if tc.serverResponse.pagedLimitHeader != "" {
					rw.Header().Set("X-Pagination-Limit", tc.serverResponse.pagedLimitHeader)
				}
				if tc.serverResponse.pagedPageHeader != "" {
					rw.Header().Set("X-Pagination-Page", tc.serverResponse.pagedPageHeader)
				}
				rw.WriteHeader(tc.serverResponse.statusCode)
				fmt.Fprintf(rw, tc.serverResponse.body)
			}))
			defer server.Close()
			t.Logf("http server: %s", server.URL)
			client, err := lokalise.NewClient("token",
				lokalise.WithBaseURL(server.URL),
				lokalise.WithLogger(&testLogger{T: t}),
			)
			if !assert.NoError(err, "unexpected client instantiation error") {
				return
			}

			resp, err := client.Teams.List(context.Background(), tc.input.options)

			if tc.output.err != nil {
				assert.EqualError(err, tc.output.err.Error(), "output error not as expected")
			} else {
				assert.NoError(err, "output error not expected")
			}
			assert.Equal(tc.output.calledPath, calledPath, "called path not as expected")
			assert.Equal(tc.output.response.Teams, resp.Teams, "response not as expected")
		})
	}
}
