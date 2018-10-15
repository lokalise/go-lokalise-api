package lokalise_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lokalise/lokalise-go-sdk/lokalise"
	"github.com/lokalise/lokalise-go-sdk/model"
	"github.com/stretchr/testify/assert"
)

func TestClient_TeamUsers_List(t *testing.T) {
	type input struct {
		teamID  int64
		options lokalise.PageOptions
	}
	type output struct {
		calledPath string
		response   model.TeamUsersResponse
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
			name: "no pagination succesful json response",
			input: input{
				teamID: 1,
			},
			serverResponse: serverResponse{
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
				calledPath: "/teams/1/users",
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
			serverResponse: serverResponse{
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
				calledPath: "/teams/1/users?limit=1&page=2",
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
			serverResponse: serverResponse{
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
				calledPath: "/teams/1/users",
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
				calledPath: "/teams/1/users",
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

			resp, err := client.TeamUsers.List(context.Background(), tc.input.teamID, tc.input.options)

			if tc.output.err != nil {
				assert.EqualError(err, tc.output.err.Error(), "output error not as expected")
			} else {
				assert.NoError(err, "output error not expected")
			}
			assert.Equal(tc.output.calledPath, calledPath, "called path not as expected")
			assert.Equal(tc.output.response.TeamID, resp.TeamID, "response team id not as expected")
			assert.Equal(tc.output.response.Paged, resp.Paged, "paged response values not as expected")
			assert.Equal(tc.output.response.TeamUsers, resp.TeamUsers, "response team user id not as expected")
		})
	}
}

func TestClient_TeamUsers_Retrieve(t *testing.T) {
	type input struct {
		teamID int64
		userID int64
	}
	type output struct {
		calledPath string
		response   model.TeamUserResponse
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
				teamID: 1,
				userID: 2,
			},
			serverResponse: serverResponse{
				statusCode: http.StatusOK,
				body: `{
					"team_id": 18821,
					"team_user": {
						"user_id": 420,
						"email": "jdoe@mycompany.com",
						"fullname": "John Doe",
						"created_at": "2018-12-31 12:00:00",
						"role": "owner"
					}
				}`,
			},
			output: output{
				calledPath: "/teams/1/users/2",
				response: model.TeamUserResponse{
					TeamID: 18821,
					TeamUser: model.TeamUser{
						UserID: 420,
						Email:  "jdoe@mycompany.com",
					},
				},
				err: nil,
			},
		},
		{
			name: "404 error response",
			input: input{
				teamID: 1,
				userID: 2,
			},
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
				calledPath: "/teams/1/users/2",
				response:   model.TeamUserResponse{},
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
				calledPath = req.URL.Path
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

			resp, err := client.TeamUsers.Retrieve(context.Background(), tc.input.teamID, tc.input.userID)

			if tc.output.err != nil {
				assert.EqualError(err, tc.output.err.Error(), "output error not as expected")
			} else {
				assert.NoError(err, "output error not expected")
			}
			assert.Equal(tc.output.calledPath, calledPath, "called path not as expected")
			assert.Equal(tc.output.response.TeamID, resp.TeamID, "response team id not as expected")
			assert.Equal(tc.output.response.TeamUser.UserID, resp.TeamUser.UserID, "response team user id not as expected")
		})
	}
}

func TestClient_TeamUsers_UpdateRole(t *testing.T) {
	type input struct {
		teamID int64
		userID int64
		role   model.TeamUserRole
	}
	type output struct {
		calledPath  string
		requestbody string
		response    model.TeamUserResponse
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
				teamID: 1,
				userID: 2,
				role:   model.TeamUserRoleAdmin,
			},
			serverResponse: serverResponse{
				statusCode: http.StatusOK,
				body: `{
					"team_id": 18821,
					"team_user": {
						"user_id": 420,
						"email": "jdoe@mycompany.com",
						"fullname": "John Doe",
						"created_at": "2018-12-31 12:00:00",
						"role": "owner"
					}
				}`,
			},
			output: output{
				calledPath:  "/teams/1/users/2",
				requestbody: `{"role":"admin"}`,
				response: model.TeamUserResponse{
					TeamID: 18821,
					TeamUser: model.TeamUser{
						UserID: 420,
						Email:  "jdoe@mycompany.com",
					},
				},
				err: nil,
			},
		},
		{
			name: "404 error response",
			input: input{
				teamID: 1,
				userID: 2,
				role:   model.TeamUserRoleAdmin,
			},
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
				calledPath:  "/teams/1/users/2",
				requestbody: `{"role":"admin"}`,
				response:    model.TeamUserResponse{},
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
			var requestbody []byte
			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				calledPath = req.URL.Path
				var err error
				requestbody, err = ioutil.ReadAll(req.Body)
				assert.NoError(err, "read request body failed")
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

			resp, err := client.TeamUsers.UpdateRole(context.Background(), tc.input.teamID, tc.input.userID, tc.input.role)

			if tc.output.err != nil {
				assert.EqualError(err, tc.output.err.Error(), "output error not as expected")
			} else {
				assert.NoError(err, "output error not expected")
			}
			assert.Equal(tc.output.calledPath, calledPath, "called path not as expected")
			assert.Equal(tc.output.requestbody, string(requestbody), "call body no as expected")
			assert.Equal(tc.output.response.TeamID, resp.TeamID, "response team id not as expected")
			assert.Equal(tc.output.response.TeamUser.UserID, resp.TeamUser.UserID, "response team user id not as expected")
		})
	}
}

func TestClient_TeamUsers_Delete(t *testing.T) {
	type input struct {
		teamID int64
		userID int64
	}
	type output struct {
		calledPath string
		response   model.TeamUserDeleteResponse
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
				teamID: 1,
				userID: 2,
			},
			serverResponse: serverResponse{
				statusCode: http.StatusOK,
				body: `{
					"team_id": 18821,
					"team_user_deleted": true
				}`,
			},
			output: output{
				calledPath: "/teams/1/users/2",
				response: model.TeamUserDeleteResponse{
					TeamID:  18821,
					Deleted: true,
				},
				err: nil,
			},
		},
		{
			name: "404 error response",
			input: input{
				teamID: 1,
				userID: 2,
			},
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
				calledPath: "/teams/1/users/2",
				response:   model.TeamUserDeleteResponse{},
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
				calledPath = req.URL.Path
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

			resp, err := client.TeamUsers.Delete(context.Background(), tc.input.teamID, tc.input.userID)

			if tc.output.err != nil {
				assert.EqualError(err, tc.output.err.Error(), "output error not as expected")
			} else {
				assert.NoError(err, "output error not expected")
			}
			assert.Equal(tc.output.calledPath, calledPath, "called path not as expected")
			assert.Equal(tc.output.response.TeamID, resp.TeamID, "response team id not as expected")
			assert.Equal(tc.output.response.Deleted, resp.Deleted, "response team user id not as expected")
		})
	}
}

type testLogger struct {
	T *testing.T
}

func (l *testLogger) Write(p []byte) (n int, err error) {
	l.T.Log(string(p))
	return len(p), nil
}
