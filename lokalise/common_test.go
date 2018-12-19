package lokalise_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/17media/go-lokalise-api/lokalise"
	"github.com/stretchr/testify/assert"
)

func notFoundResponseBody(message string) string {
	return fmt.Sprintf(`{
		"error": {
			"code": 404,
			"message": "%s"
		}
	}`, message)
}

type fixture struct {
	calledPath   string
	calledMethod string
	requestBody  string
}

type outgoingRequest struct {
	path   string
	method string
	body   string
}

// Assert asserts that provided fixture matched the expectations of the outgoing request.
func (req *outgoingRequest) Assert(t *testing.T, fixture *fixture) {
	assert.Equal(t, req.path, fixture.calledPath, "called path not as expected")
	assert.Equal(t, req.method, fixture.calledMethod, "called path not as expected")
	assert.Equal(t, req.body, fixture.requestBody, "call body no as expected")
}

// setupClient returns a lokalise.Client that is configured to hit a test HTTP server returning provided
// body on any requests. The returned fixture holds recorded values of the outgoing request.
//
// The function return value is a close function to stop the test server. Ensure to call this
// between tests.
//
//    client, fixture, close := setup(t, "")
//    defer close()
func setupClient(t *testing.T, serverResponseBody string) (*lokalise.Client, *fixture, func()) {
	return setupServerAndClient(t, func(f *fixture) http.HandlerFunc {
		return func(rw http.ResponseWriter, req *http.Request) {
			f.calledMethod = req.Method
			f.calledPath = req.URL.Path
			requestBody, err := ioutil.ReadAll(req.Body)
			if err != nil {
				t.Fatalf("read request body failed: %v", err)
			}
			f.requestBody = string(requestBody)
			rw.Header().Set("Content-Type", "application/json")
			rw.WriteHeader(http.StatusOK)
			fmt.Fprint(rw, serverResponseBody)
		}
	})
}

type serverListResponse struct {
	body                  string
	statusCode            int
	pagedTotalCountHeader string
	pagedPageCountHeader  string
	pagedLimitHeader      string
	pagedPageHeader       string
}

// setupListClient returns a lokalise.Client that is configured to hit a test HTTP server returning provided
// body and list paramters on any requests. The returned fixture holds recorded values of the outgoing
// request.
//
// The function return value is a close function to stop the test server. Ensure to call this
// between tests.
//
//    client, fixture, close := setup(t, "")
//    defer close()
func setupListClient(t *testing.T, serverResponse serverListResponse) (*lokalise.Client, *fixture, func()) {
	return setupServerAndClient(t, func(f *fixture) http.HandlerFunc {
		return func(rw http.ResponseWriter, req *http.Request) {
			f.calledMethod = req.Method
			f.calledPath = req.URL.String()
			rw.Header().Set("Content-Type", "application/json")
			if serverResponse.pagedTotalCountHeader != "" {
				rw.Header().Set("X-Pagination-Total-Count", serverResponse.pagedTotalCountHeader)
			}
			if serverResponse.pagedPageCountHeader != "" {
				rw.Header().Set("X-Pagination-Page-Count", serverResponse.pagedPageCountHeader)
			}
			if serverResponse.pagedLimitHeader != "" {
				rw.Header().Set("X-Pagination-Limit", serverResponse.pagedLimitHeader)
			}
			if serverResponse.pagedPageHeader != "" {
				rw.Header().Set("X-Pagination-Page", serverResponse.pagedPageHeader)
			}
			rw.WriteHeader(serverResponse.statusCode)
			fmt.Fprint(rw, serverResponse.body)
		}
	})
}

// setupServerAndClient returns a lokalise.Client that is configured to hit a test HTTP server
// handled with the provided http.HandlerFunc.
//
// The HandlerFunc has access to a fixture pointer to report information to the tests about
// what endpoint was called, the status code and payloads.
func setupServerAndClient(t *testing.T, handler func(*fixture) http.HandlerFunc) (*lokalise.Client, *fixture, func()) {
	f := fixture{}
	server := httptest.NewServer(handler(&f))
	t.Logf("http server: %s", server.URL)
	client, err := lokalise.NewClient("token",
		lokalise.WithBaseURL(server.URL),
		lokalise.WithLogger(&testLogger{T: t}),
	)
	if err != nil {
		server.Close()
		t.Fatalf("unexpected client instantiation error: %v", err)
	}
	return client, &f, server.Close
}
