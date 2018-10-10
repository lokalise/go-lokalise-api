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
	defaultBaseURL = "https://api.lokalise.co/api2"
)

type Client struct {
	timeout    time.Duration
	baseURL    string
	apiToken   string
	retryCount int
	httpClient *resty.Client
	logger     io.Writer

	Teams TeamsService
}

type ClientOption func(*Client) error

func NewClient(apiToken string, options ...ClientOption) (*Client, error) {
	c := Client{
		apiToken:   apiToken,
		retryCount: 3,
		baseURL:    defaultBaseURL,
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
	c.httpClient = resty.New().
		SetHostURL(c.baseURL).
		SetRetryCount(c.retryCount).
		SetHeader(apiTokenHeader, c.apiToken).
		SetLogger(c.logger).
		SetError(RequestError{}).
		AddRetryCondition(requestRetryCondition())

	c.Teams = TeamsService{httpClient: c.httpClient}
	return &c, nil
}

// WithRetryCount returns a client ClientOption setting the retry count of outgoing requests.
// if count is zero retries are disabled.
func WithRetryCount(count int) ClientOption {
	return func(c *Client) error {
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

// WithLogger returns a ClientOption setting the output destination of any log messages.
func WithLogger(l io.Writer) ClientOption {
	return func(c *Client) error {
		if l == nil {
			return errors.New("lokalise: logger value required")
		}
		c.logger = l
		return nil
	}
}

// WithBaseURL returns a ClientOption setting the base URL of the client.
//
// This should only be used for testing different API versions or for using a mocked
// backend in tests.
func WithBaseURL(url string) ClientOption {
	return func(c *Client) error {
		c.baseURL = url
		return nil
	}
}
