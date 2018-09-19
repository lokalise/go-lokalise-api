package lokalise

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-resty/resty"
)

func TestClient_retryLogic(t *testing.T) {
	type input struct {
		retryCount  int
		okAtAttempt int
	}
	type output struct {
		attempts   int
		statusCode int
		err        error
	}
	tt := []struct {
		name   string
		input  input
		output output
	}{
		{
			name: "internal error on all attempts",
			input: input{
				retryCount:  2,
				okAtAttempt: 5,
			},
			output: output{
				attempts:   2,
				statusCode: http.StatusInternalServerError,
				err:        nil,
			},
		},
		{
			name: "ok at last attempts",
			input: input{
				retryCount:  2,
				okAtAttempt: 2,
			},
			output: output{
				attempts:   2,
				statusCode: http.StatusOK,
				err:        nil,
			},
		},
		{
			name: "zero retries",
			input: input{
				retryCount:  0,
				okAtAttempt: 2,
			},
			output: output{
				attempts:   1,
				statusCode: http.StatusInternalServerError,
				err:        nil,
			},
		},
		{
			name: "default to three retries",
			input: input{
				retryCount:  -1, // disables the retry option in test
				okAtAttempt: 3,
			},
			output: output{
				attempts:   3,
				statusCode: http.StatusOK,
				err:        nil,
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			attempts := 0
			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				attempts++
				if attempts >= tc.input.okAtAttempt {
					rw.WriteHeader(http.StatusOK)
					return
				}
				rw.WriteHeader(http.StatusInternalServerError)
			}))
			var opts []option
			if tc.input.retryCount >= 0 {
				opts = append(opts, withRetryCount(tc.input.retryCount))
			}
			c, err := newClient("token", opts...)
			if err != nil {
				t.Fatalf("client instantiation error: %v", err)
			}
			response, err := c.client.R().Get(server.URL)
			if tc.output.err != nil {
				if err == nil {
					t.Fatalf("expected error %s but got nil", tc.output.err)
				}
				if tc.output.err.Error() != err.Error() {
					t.Fatalf("wrong error: expected '%s': got '%s'", tc.output.err.Error(), err.Error())
				}
				return
			}
			if err != nil {
				t.Fatalf("expected no error but got '%s'", err.Error())
			}
			if response == nil {
				t.Fatal("expected a response but got nil")
			}
			if response.StatusCode() != tc.output.statusCode {
				t.Fatalf("wrong response status code: expected %d: got %d", tc.output.statusCode, response.StatusCode())
			}
			if tc.output.attempts != attempts {
				t.Fatalf("expected %d request attempts: got %d", tc.output.attempts, attempts)
			}
		})
	}
}

func TestClient_contextCancelation(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		attempts++
		time.Sleep(10 * time.Millisecond)
		rw.WriteHeader(http.StatusOK)
	}))
	c, err := newClient("token")
	if err != nil {
		t.Fatalf("client instantiation error: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err = c.client.R().SetContext(ctx).Get(server.URL)
	if err == nil {
		t.Fatal("expected an error but got nil")
	}
	if ctx.Err() == nil {
		t.Fatal("expected context to be canceled")
	}
	if attempts != 0 {
		t.Fatalf("expected context cancelation not to be retried but got %d attempts", attempts)
	}
}

func TestClient_apiTokenHeader(t *testing.T) {
	apiToken := "api-token"
	var requestHeaderContent string
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		requestHeaderContent = req.Header.Get(apiTokenHeader)
		rw.WriteHeader(http.StatusOK)
	}))
	c, err := newClient(apiToken)
	if err != nil {
		t.Fatalf("client instantiation error: %v", err)
	}

	response, err := c.client.R().Get(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if response.StatusCode() != http.StatusOK {
		t.Fatalf("expected status code 200 OK: got %d", response.StatusCode())
	}
	if requestHeaderContent != apiToken {
		t.Fatalf("expected header %s to contain '%s': got '%s'", apiTokenHeader, apiToken, requestHeaderContent)
	}
}

func TestClient_negativeRetries(t *testing.T) {
	_, err := newClient("token", withRetryCount(-1))
	if err == nil {
		t.Fatal("expected an error but got nil")
	}
}

func TestRequestRetryCondition(t *testing.T) {
	tt := []struct {
		name        string
		response    *resty.Response
		shouldRetry bool
		err         error
	}{
		{
			name:        "nil response",
			response:    nil,
			shouldRetry: true,
			err:         nil,
		},
		{
			name: "status code 200",
			response: &resty.Response{
				Request: &resty.Request{},
				RawResponse: &http.Response{
					StatusCode: 200,
				},
			},
			shouldRetry: false,
			err:         nil,
		},
		{
			name: "status code 499",
			response: &resty.Response{
				Request: &resty.Request{},
				RawResponse: &http.Response{
					StatusCode: 499,
				},
			},
			shouldRetry: false,
			err:         nil,
		},
		{
			name: "status code 500",
			response: &resty.Response{
				Request: &resty.Request{},
				RawResponse: &http.Response{
					StatusCode: 500,
				},
			},
			shouldRetry: true,
			err:         nil,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			retry, err := requestRetryCondition()(tc.response)
			if tc.err != nil {
				if err == nil {
					t.Fatalf("expected error %s but got nil", tc.err)
				}
				if tc.err.Error() != err.Error() {
					t.Fatalf("wrong error: expected '%s': got '%s'", tc.err.Error(), err.Error())
				}
				return
			}
			if err != nil {
				t.Fatalf("expected no error but got '%s'", err.Error())
			}
			if tc.shouldRetry != retry {
				t.Fatalf("expected retry to be %t but got %t", tc.shouldRetry, retry)
			}
		})
	}
}
