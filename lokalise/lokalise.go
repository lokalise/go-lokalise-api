// Package lokalise provides functions to access the Lokalise web API.
package lokalise

import (
	"io"
	"time"

	"github.com/go-resty/resty"
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
	c := client{}
	for _, o := range options {
		err := o(&c)
		if err != nil {
			return nil, err
		}
	}
	return &c, nil
}
