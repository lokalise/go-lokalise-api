package lokalise_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lokalise/go-lokalise-api/lokalise"
	"github.com/lokalise/go-lokalise-api/model"
	"github.com/stretchr/testify/assert"
)

func TestClient_Projects_List(t *testing.T) {
	type input struct {
		options lokalise.ProjectsOptions
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
				calledPath: "/projects?filter_team_id=3&limit=1&page=2",
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
				body:       notFoundResponseBody("project not found"),
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
					Message: "project not found",
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

func TestClient_Projects_Create(t *testing.T) {
	type input struct {
		name        string
		description string
	}
	type output struct {
		calledPath  string
		requestbody string
		response    model.Project
		err         error
	}
	type serverResponse struct {
		body       string
		statusCode int
	}
	tt := []struct {
		name           string
		input          input
		serverResponse serverResponse
		output         output
	}{
		{
			name: "succesful json response",
			input: input{
				name:        "TheApp Project",
				description: "iOS + Android strings of TheApp. https://theapp.com",
			},
			serverResponse: serverResponse{
				statusCode: http.StatusOK,
				body: `{
					"project_id": "3002780358964f9bab5a92.87762498",
					"name": "TheApp Project",
					"description": "iOS + Android strings of TheApp. https://theapp.com",
					"created_at": "2018-12-31 12:00:00",
					"created_by": 420,
					"created_by_email": "user@mycompany.com",
					"team_id": 12345
				}`,
			},
			output: output{
				calledPath:  "/projects",
				requestbody: `{"description":"iOS + Android strings of TheApp. https://theapp.com","name":"TheApp Project"}`,
				response: model.Project{
					ProjectID:      "3002780358964f9bab5a92.87762498",
					Name:           "TheApp Project",
					Description:    "iOS + Android strings of TheApp. https://theapp.com",
					CreatedAt:      "2018-12-31 12:00:00",
					CreatedBy:      420,
					CreatedByEmail: "user@mycompany.com",
					TeamID:         12345,
				},
				err: nil,
			},
		},
		{
			name: "404 error response",
			input: input{
				name:        "name",
				description: "description",
			},
			serverResponse: serverResponse{
				statusCode: http.StatusNotFound,
				body:       notFoundResponseBody("project not found"),
			},
			output: output{
				calledPath:  "/projects",
				requestbody: `{"description":"description","name":"name"}`,
				response:    model.Project{},
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

			var calledPath string
			var requestbody []byte
			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				calledPath = req.URL.Path
				var err error
				requestbody, err = ioutil.ReadAll(req.Body)
				assert.NoError(err, "read request body failed")
				assert.Equal(req.Method, "POST", "wrong HTTP request verb")
				rw.Header().Set("Content-Type", "application/json")
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

			resp, err := client.Projects.Create(context.Background(), tc.input.name, tc.input.description)

			if tc.output.err != nil {
				assert.EqualError(err, tc.output.err.Error(), "output error not as expected")
			} else {
				assert.NoError(err, "output error not expected")
			}
			assert.Equal(tc.output.calledPath, calledPath, "called path not as expected")
			assert.Equal(tc.output.requestbody, string(requestbody), "call body no as expected")
			assert.Equal(tc.output.response.ProjectID, resp.ProjectID, "response project id not as expected")
			assert.Equal(tc.output.response.Name, resp.Name, "response name not as expected")
			assert.Equal(tc.output.response.Description, resp.Description, "response description not as expected")
			assert.Equal(tc.output.response.CreatedAt, resp.CreatedAt, "response created at not as expected")
			assert.Equal(tc.output.response.CreatedBy, resp.CreatedBy, "response created by not as expected")
			assert.Equal(tc.output.response.CreatedByEmail, resp.CreatedByEmail, "response created by email not as expected")
			assert.Equal(tc.output.response.TeamID, resp.TeamID, "response team id not as expected")
		})
	}
}

func TestClient_Projects_CreateForTeam(t *testing.T) {
	type input struct {
		name        string
		description string
		teamID      int64
	}
	type output struct {
		calledPath  string
		requestbody string
		response    model.Project
		err         error
	}
	type serverResponse struct {
		body       string
		statusCode int
	}
	tt := []struct {
		name           string
		input          input
		serverResponse serverResponse
		output         output
	}{
		{
			name: "succesful json response",
			input: input{
				name:        "TheApp Project",
				description: "iOS + Android strings of TheApp. https://theapp.com",
				teamID:      12345,
			},
			serverResponse: serverResponse{
				statusCode: http.StatusOK,
				body: `{
					"project_id": "3002780358964f9bab5a92.87762498",
					"name": "TheApp Project",
					"description": "iOS + Android strings of TheApp. https://theapp.com",
					"created_at": "2018-12-31 12:00:00",
					"created_by": 420,
					"created_by_email": "user@mycompany.com",
					"team_id": 12345
				}`,
			},
			output: output{
				calledPath:  "/projects",
				requestbody: `{"description":"iOS + Android strings of TheApp. https://theapp.com","name":"TheApp Project","team_id":12345}`,
				response: model.Project{
					ProjectID:      "3002780358964f9bab5a92.87762498",
					Name:           "TheApp Project",
					Description:    "iOS + Android strings of TheApp. https://theapp.com",
					CreatedAt:      "2018-12-31 12:00:00",
					CreatedBy:      420,
					CreatedByEmail: "user@mycompany.com",
					TeamID:         12345,
				},
				err: nil,
			},
		},
		{
			name: "404 error response",
			input: input{
				name:        "name",
				description: "description",
				teamID:      12345,
			},
			serverResponse: serverResponse{
				statusCode: http.StatusNotFound,
				body:       notFoundResponseBody("project not found"),
			},
			output: output{
				calledPath:  "/projects",
				requestbody: `{"description":"description","name":"name","team_id":12345}`,
				response:    model.Project{},
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

			var calledPath string
			var requestbody []byte
			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				calledPath = req.URL.Path
				var err error
				requestbody, err = ioutil.ReadAll(req.Body)
				assert.NoError(err, "read request body failed")
				assert.Equal(req.Method, "POST", "wrong HTTP request verb")
				rw.Header().Set("Content-Type", "application/json")
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

			resp, err := client.Projects.CreateForTeam(context.Background(), tc.input.name, tc.input.description, tc.input.teamID)

			if tc.output.err != nil {
				assert.EqualError(err, tc.output.err.Error(), "output error not as expected")
			} else {
				assert.NoError(err, "output error not expected")
			}
			assert.Equal(tc.output.calledPath, calledPath, "called path not as expected")
			assert.Equal(tc.output.requestbody, string(requestbody), "call body no as expected")
			assert.Equal(tc.output.response.ProjectID, resp.ProjectID, "response project id not as expected")
			assert.Equal(tc.output.response.Name, resp.Name, "response name not as expected")
			assert.Equal(tc.output.response.Description, resp.Description, "response description not as expected")
			assert.Equal(tc.output.response.CreatedAt, resp.CreatedAt, "response created at not as expected")
			assert.Equal(tc.output.response.CreatedBy, resp.CreatedBy, "response created by not as expected")
			assert.Equal(tc.output.response.CreatedByEmail, resp.CreatedByEmail, "response created by email not as expected")
			assert.Equal(tc.output.response.TeamID, resp.TeamID, "response team id not as expected")
		})
	}
}

func TestClient_Projects_Retrieve(t *testing.T) {
	type input struct {
		projectID string
	}
	type output struct {
		calledPath string
		response   model.Project
		err        error
	}
	type serverResponse struct {
		body       string
		statusCode int
	}
	tt := []struct {
		name           string
		input          input
		serverResponse serverResponse
		output         output
	}{
		{
			name: "succesful json response",
			input: input{
				projectID: "1",
			},
			serverResponse: serverResponse{
				statusCode: http.StatusOK,
				body: `{
					"project_id": "3002780358964f9bab5a92.87762498",
					"name": "TheApp Project",
					"description": "iOS + Android strings of TheApp. https://theapp.com",
					"created_at": "2018-12-31 12:00:00",
					"created_by": 420,
					"created_by_email": "user@mycompany.com",
					"team_id": 12345
				}`,
			},
			output: output{
				calledPath: "/projects/1",
				response: model.Project{
					ProjectID:      "3002780358964f9bab5a92.87762498",
					Name:           "TheApp Project",
					Description:    "iOS + Android strings of TheApp. https://theapp.com",
					CreatedAt:      "2018-12-31 12:00:00",
					CreatedBy:      420,
					CreatedByEmail: "user@mycompany.com",
					TeamID:         12345,
				},
				err: nil,
			},
		},
		{
			name: "404 error response",
			input: input{
				projectID: "1",
			},
			serverResponse: serverResponse{
				statusCode: http.StatusNotFound,
				body:       notFoundResponseBody("project not found"),
			},
			output: output{
				calledPath: "/projects/1",
				response:   model.Project{},
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

			var calledPath string
			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				calledPath = req.URL.Path
				rw.Header().Set("Content-Type", "application/json")
				assert.Equal("GET", req.Method, "request HTTP verb not as expected")
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

			resp, err := client.Projects.Retrieve(context.Background(), tc.input.projectID)

			if tc.output.err != nil {
				assert.EqualError(err, tc.output.err.Error(), "output error not as expected")
			} else {
				assert.NoError(err, "output error not expected")
			}
			assert.Equal(tc.output.calledPath, calledPath, "called path not as expected")
			assert.Equal(tc.output.response.TeamID, resp.TeamID, "response team id not as expected")
			assert.Equal(tc.output.response.Name, resp.Name, "response name not as expected")
		})
	}
}

func TestClient_Projects_Update(t *testing.T) {
	type input struct {
		projectID   string
		name        string
		description string
	}
	type output struct {
		calledPath  string
		requestbody string
		response    model.Project
		err         error
	}
	type serverResponse struct {
		body       string
		statusCode int
	}
	tt := []struct {
		name           string
		input          input
		serverResponse serverResponse
		output         output
	}{
		{
			name: "succesful json response",
			input: input{
				projectID:   "3002780358964f9bab5a92",
				name:        "TheZapp Project",
				description: "iOS + Android strings of TheZapp. https://thezapp.com",
			},
			serverResponse: serverResponse{
				statusCode: http.StatusOK,
				body: `{
					"project_id": "3002780358964f9bab5a92.87762498",
					"name": "TheApp Project",
					"description": "iOS + Android strings of TheApp. https://theapp.com",
					"created_at": "2018-12-31 12:00:00 (Etc/UTC)",
					"created_by": 420,
					"created_by_email": "user@mycompany.com",
					"team_id": 12345
				}`,
			},
			output: output{
				calledPath:  "/projects/3002780358964f9bab5a92",
				requestbody: `{"description":"iOS + Android strings of TheZapp. https://thezapp.com","name":"TheZapp Project"}`,
				response: model.Project{
					ProjectID:      "3002780358964f9bab5a92.87762498",
					Name:           "TheApp Project",
					Description:    "iOS + Android strings of TheApp. https://theapp.com",
					CreatedAt:      "2018-12-31 12:00:00 (Etc/UTC)",
					CreatedBy:      420,
					CreatedByEmail: "user@mycompany.com",
					TeamID:         12345,
				},
				err: nil,
			},
		},
		{
			name: "404 error response",
			input: input{
				projectID:   "12345",
				name:        "name",
				description: "description",
			},
			serverResponse: serverResponse{
				statusCode: http.StatusNotFound,
				body:       notFoundResponseBody("project not found"),
			},
			output: output{
				calledPath:  "/projects/12345",
				requestbody: `{"description":"description","name":"name"}`,
				response:    model.Project{},
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

			var calledPath string
			var requestbody []byte
			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				calledPath = req.URL.Path
				var err error
				requestbody, err = ioutil.ReadAll(req.Body)
				assert.NoError(err, "read request body failed")
				assert.Equal(req.Method, "PUT", "wrong HTTP request verb")
				rw.Header().Set("Content-Type", "application/json")
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

			resp, err := client.Projects.Update(context.Background(), tc.input.projectID, tc.input.name, tc.input.description)

			if tc.output.err != nil {
				assert.EqualError(err, tc.output.err.Error(), "output error not as expected")
			} else {
				assert.NoError(err, "output error not expected")
			}
			assert.Equal(tc.output.calledPath, calledPath, "called path not as expected")
			assert.Equal(tc.output.requestbody, string(requestbody), "call body no as expected")
			assert.Equal(tc.output.response.ProjectID, resp.ProjectID, "response project id not as expected")
			assert.Equal(tc.output.response.Name, resp.Name, "response name not as expected")
			assert.Equal(tc.output.response.Description, resp.Description, "response description not as expected")
			assert.Equal(tc.output.response.CreatedAt, resp.CreatedAt, "response created at not as expected")
			assert.Equal(tc.output.response.CreatedBy, resp.CreatedBy, "response created by not as expected")
			assert.Equal(tc.output.response.CreatedByEmail, resp.CreatedByEmail, "response created by email not as expected")
			assert.Equal(tc.output.response.TeamID, resp.TeamID, "response team id not as expected")
		})
	}
}

func TestClient_Projects_Empty(t *testing.T) {
	type input struct {
		projectID string
	}
	type output struct {
		calledPath  string
		requestbody string
		response    model.ProjectEmptyResponse
		err         error
	}
	type serverResponse struct {
		body       string
		statusCode int
	}
	tt := []struct {
		name           string
		input          input
		serverResponse serverResponse
		output         output
	}{
		{
			name: "succesful json response",
			input: input{
				projectID: "3002780358964f9bab5a92.87762498",
			},
			serverResponse: serverResponse{
				statusCode: http.StatusOK,
				body: `{
					"project_id": "3002780358964f9bab5a92.87762498",
					"keys_deleted": true
				}`,
			},
			output: output{
				calledPath:  "/projects/3002780358964f9bab5a92.87762498/empty",
				requestbody: ``,
				response: model.ProjectEmptyResponse{
					ProjectID:   "3002780358964f9bab5a92.87762498",
					KeysDeleted: true,
				},
				err: nil,
			},
		},
		{
			name: "404 error response",
			input: input{
				projectID: "12345",
			},
			serverResponse: serverResponse{
				statusCode: http.StatusNotFound,
				body:       notFoundResponseBody("project not found"),
			},
			output: output{
				calledPath:  "/projects/12345/empty",
				requestbody: ``,
				response:    model.ProjectEmptyResponse{},
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

			var calledPath string
			var requestbody []byte
			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				calledPath = req.URL.Path
				var err error
				requestbody, err = ioutil.ReadAll(req.Body)
				assert.NoError(err, "read request body failed")
				assert.Equal(req.Method, "PUT", "wrong HTTP request verb")
				rw.Header().Set("Content-Type", "application/json")
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

			resp, err := client.Projects.Empty(context.Background(), tc.input.projectID)

			if tc.output.err != nil {
				assert.EqualError(err, tc.output.err.Error(), "output error not as expected")
			} else {
				assert.NoError(err, "output error not expected")
			}
			assert.Equal(tc.output.calledPath, calledPath, "called path not as expected")
			assert.Equal(tc.output.requestbody, string(requestbody), "call body no as expected")
			assert.Equal(tc.output.response.ProjectID, resp.ProjectID, "response project id not as expected")
			assert.Equal(tc.output.response.KeysDeleted, resp.KeysDeleted, "response keys deleted not as expected")
		})
	}
}
