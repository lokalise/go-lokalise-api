package lokalise

import (
	"context"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
)

const (
	apiTokenHeader    = "X-Api-Token"
	defaultBaseURL    = "https://api.lokalise.co/api2"
	defaultRetryCount = 3
)

type restClient struct {
	*resty.Client

	timeout    time.Duration
	baseURL    string
	apiToken   string
	retryCount int
}

func newClient(apiToken string) *restClient {
	c := restClient{
		apiToken:   apiToken,
		retryCount: defaultRetryCount,
		baseURL:    defaultBaseURL,
	}

	c.Client = resty.New().
		SetHostURL(c.baseURL).
		SetRetryCount(c.retryCount).
		SetHeader(apiTokenHeader, c.apiToken).
		SetError(errorResponse{}).
		AddRetryCondition(requestRetryCondition())

	return &c
}

// requestRetryCondition indicates a retry if the HTTP status code of the response
// is >= 500.
// failing requests due to network conditions, eg. "no such host", are handled by resty internally
func requestRetryCondition() resty.RetryConditionFunc {
	return func(res *resty.Response, err error) bool {
		if res == nil || err != nil {
			return true
		}
		if res.StatusCode() >= http.StatusInternalServerError {
			return true
		}
		return false
	}
}

func (c *restClient) get(ctx context.Context, path string, res interface{}) (*resty.Response, error) {
	return c.req(ctx, path, res).Get(path)
}

func (c *restClient) getList(ctx context.Context, path string, res interface{}, options OptionsApplier) (*resty.Response, error) {
	req := c.req(ctx, path, res)
	options.Apply(req)
	return req.Get(path)
}

func (c *restClient) post(ctx context.Context, path string, res, body interface{}) (*resty.Response, error) {
	return c.reqWithBody(ctx, path, res, body).Post(path)
}

func (c *restClient) put(ctx context.Context, path string, res, body interface{}) (*resty.Response, error) {
	return c.reqWithBody(ctx, path, res, body).Put(path)
}

func (c *restClient) delete(ctx context.Context, path string, res interface{}) (*resty.Response, error) {
	return c.req(ctx, path, res).Delete(path)
}
func (c *restClient) deleteWithBody(ctx context.Context, path string, res, body interface{}) (*resty.Response, error) {
	return c.reqWithBody(ctx, path, res, body).Delete(path)
}

func (c *restClient) req(ctx context.Context, path string, res interface{}) *resty.Request {
	return c.R().
		SetResult(&res).
		SetContext(ctx)
}

func (c *restClient) reqWithBody(ctx context.Context, path string, res, body interface{}) *resty.Request {
	return c.req(ctx, path, res).SetBody(body)
}
