package lokalise

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

const (
	pathScreenshots = "screenshots"
)

type Screenshot struct {
	ScreenshotID   int64     `json:"screenshot_id,omitempty"`
	KeyIDs         []int64   `json:"key_ids,omitempty"`
	URL            string    `json:"url,omitempty"`
	Title          string    `json:"title,omitempty"`
	Description    string    `json:"description,omitempty"`
	ScreenshotTags []string  `json:"screenshot_tags,omitempty"`
	Width          int64     `json:"width,omitempty"`
	Height         int64     `json:"height,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	CreatedAtTs    int64     `json:"created_at_timestamp"`
}

type ScreenshotResponse struct {
	Paged
	Screenshots []Screenshot `json:"screenshots,omitempty"`
}

type ScreenshotsService struct {
	client *Client
}

type ScreenshotsOptions struct {
	PageOptions
	IncludeTags bool // todo use `query:"include_tags"`
	ListOnly    bool
}

type CreateScreenshotOptions struct {
	Body        string   `json:"data"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Ocr         bool     `json:"ocr"`
	KeyIds      []int64  `json:"key_ids"`
	Tags        []string `json:"tags"`
}

func (options ScreenshotsOptions) Apply(req *resty.Request) {
	options.PageOptions.Apply(req)
	if options.IncludeTags {
		req.SetQueryParam("include_tags", fmt.Sprintf("%v", options.IncludeTags))
	}
	if options.ListOnly {
		req.SetQueryParam("list_only", fmt.Sprintf("%v", options.ListOnly))
	}
}

func (c *ScreenshotsService) List(ctx context.Context, pageOptions ScreenshotsOptions) (ScreenshotResponse, error) {
	var result ScreenshotResponse
	resp, err := c.client.getList(ctx, pathScreenshots, &result, pageOptions)
	if err != nil {
		return result, err
	}

	applyPaged(resp, &result.Paged)
	return result, apiError(resp)
}

func (c *ScreenshotsService) Create(ctx context.Context, projectID string, options CreateScreenshotOptions) (ScreenshotResponse, error) {
	var res ScreenshotResponse
	body, err := json.Marshal(options)
	if err != nil {
		return res, err
	}
	resp, err := c.client.post(ctx, fmt.Sprintf("%s/%s", projectID, pathScreenshots), &res, string(body))
	if err != nil {
		return res, err
	}

	return res, apiError(resp)
}
