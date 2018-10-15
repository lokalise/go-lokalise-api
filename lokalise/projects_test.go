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

func TestClient_Projects_List(t *testing.T) {
	type input struct {
		options lokalise.PageOptions
	}
	type output struct {
		calledPath string
		response   model.ProjectsResponse
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
				calledPath: "/projects",
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
				calledPath: "/projects?limit=1&page=2",
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
				calledPath: "/projects",
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

			resp, err := client.Projects.List(context.Background(), tc.input.options)

			if tc.output.err != nil {
				assert.EqualError(err, tc.output.err.Error(), "output error not as expected")
			} else {
				assert.NoError(err, "output error not expected")
			}
			assert.Equal(tc.output.calledPath, calledPath, "called path not as expected")
			assert.Equal(tc.output.response.Projects, resp.Projects, "response not as expected")
		})
	}
}
