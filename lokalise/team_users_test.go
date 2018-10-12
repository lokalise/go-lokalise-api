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
					"code": 404,
					"message": "team not found"
				}`,
			},
			output: output{
				calledPath: "/teams/1/users/2",
				response:   model.TeamUserResponse{},
				err: &lokalise.RequestError{
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

type testLogger struct {
	T *testing.T
}

func (l *testLogger) Write(p []byte) (n int, err error) {
	l.T.Log(string(p))
	return len(p), nil
}
