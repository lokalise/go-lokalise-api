package lokalise

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/17media/go-lokalise-api/model"
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
			opts := []ClientOption{WithLogger(&testLogger{T: t})}
			if tc.input.retryCount >= 0 {
				opts = append(opts, WithRetryCount(tc.input.retryCount))
			}
			c, err := NewClient("token", opts...)
			if err != nil {
				t.Fatalf("client instantiation error: %v", err)
			}
			response, err := c.httpClient.R().Get(server.URL)
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
	c, err := NewClient("token", WithLogger(&testLogger{T: t}))
	if err != nil {
		t.Fatalf("client instantiation error: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err = c.httpClient.R().SetContext(ctx).Get(server.URL)
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
	c, err := NewClient(apiToken, WithLogger(&testLogger{T: t}))
	if err != nil {
		t.Fatalf("client instantiation error: %v", err)
	}

	response, err := c.httpClient.R().Get(server.URL)

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
	_, err := NewClient("token", WithLogger(&testLogger{T: t}), WithRetryCount(-1))
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

func TestClient_nilLogger(t *testing.T) {
	client, err := NewClient("", WithLogger(nil))
	if err == nil {
		t.Errorf("expected error but got nil")
	} else {
		if err.Error() != "lokalise: logger value required" {
			t.Errorf("expected error %s but got %s", "lokalise: logger value required", err.Error())
		}
	}
	if client != nil {
		t.Errorf("expected client to be nil but got %v", *client)
	}
}

type testWriter struct {
	lines []string
}

func (tw *testWriter) Write(p []byte) (n int, err error) {
	tw.lines = append(tw.lines, string(p))
	return len(p), nil
}

func TestClient_customLogger(t *testing.T) {
	client, err := NewClient("", WithLogger(&testWriter{}))
	if err != nil {
		t.Errorf("expected no error but got '%s'", err.Error())
	}
	if client == nil {
		t.Error("expected client to be instantiated but got nil")
	}
}

func TestClient_errorModel(t *testing.T) {
	serverResponse := errorResponse{
		Error: model.Error{
			Code:    http.StatusInternalServerError,
			Message: "some server error",
		},
	}
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		data, err := json.Marshal(serverResponse)
		if err != nil {
			t.Fatalf("failed to marshal server error: %v", err)
		}
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(serverResponse.Error.Code)
		rw.Write(data)
	}))
	c, err := NewClient("token", WithLogger(&testLogger{T: t}), WithRetryCount(0))
	if err != nil {
		t.Fatalf("client instantiation error: %v", err)
	}

	response, err := c.httpClient.R().
		SetResult(make(map[string]interface{})).
		Get(server.URL)

	if err != nil {
		t.Fatalf("expected no error but got '%s'", err.Error())
	}
	if response == nil {
		t.Fatal("expected a response but got nil")
	}
	if response.StatusCode() != serverResponse.Error.Code {
		t.Errorf("wrong response status code: expected %d: got %d", serverResponse.Error.Code, response.StatusCode())
	}
	requestErr := response.Error()
	if requestErr == nil {
		t.Fatal("expected request error but got nil")
	}
	requestErrModel, ok := requestErr.(*errorResponse)
	if !ok {
		t.Fatalf("expected request error to be type %T but got %T", &model.Error{}, requestErr)
	}
	if requestErrModel.Error.Code != http.StatusInternalServerError {
		t.Errorf("wrong error code: expected %d: got %d", http.StatusInternalServerError, requestErrModel.Error.Code)
	}
	if requestErrModel.Error.Message != "some server error" {
		t.Errorf("wrong message: expected '%s': got '%s'", "some server error", requestErrModel.Error.Message)
	}
}

type testLogger struct {
	T *testing.T
}

func (l *testLogger) Write(p []byte) (n int, err error) {
	l.T.Log(string(p))
	return len(p), nil
}
