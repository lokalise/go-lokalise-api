// Package lokalise provides functions to access the Lokalise web API.
package lokalise

import (
	"errors"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/go-resty/resty"
)

const (
	apiTokenHeader = "X-Api-Token"
)

type client struct {
	timeout    time.Duration
	baseURL    string
	apiToken   string
	retryCount int
	client     *resty.Client
	logger     io.Writer
}

type option func(*client) error

func newClient(apiToken string, options ...option) (*client, error) {
	c := client{
		apiToken:   apiToken,
		retryCount: 3,
	}
	for _, o := range options {
		err := o(&c)
		if err != nil {
			return nil, err
		}
	}
	if c.logger == nil {
		c.logger = os.Stderr
	}
	c.client = resty.New().
		SetRetryCount(c.retryCount).
		SetHeader(apiTokenHeader, c.apiToken).
		SetLogger(c.logger).
		SetError(RequestError{}).
		AddRetryCondition(requestRetryCondition())

	return &c, nil
}

// withRetryCount returns a client option setting the retry count of outgoing requests.
// if count is zero retries are disabled.
func withRetryCount(count int) option {
	return func(c *client) error {
		if count < 0 {
			return errors.New("lokalise: retry count must be positive")
		}
		c.retryCount = count
		return nil
	}
}

// requestRetryCondition indicates a retry if the HTTP status code of the response
// is >= 500.
// failing requests due to network conditions, eg. "no such host", are handled by resty internally
func requestRetryCondition() resty.RetryConditionFunc {
	return func(res *resty.Response) (bool, error) {
		if res == nil {
			return true, nil
		}
		if res.StatusCode() >= http.StatusInternalServerError {
			return true, nil
		}
		return false, nil
	}
}

// WithLoggerFunc returns a RestyOption configurring the logging function for internal Resty logs.
func withLogger(l io.Writer) option {
	return func(c *client) error {
		if l == nil {
			return errors.New("lokalise: logger value required")
		}
		c.logger = l
		return nil
	}
}

// RequestError is the API error model.
type RequestError struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}
