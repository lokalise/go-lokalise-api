package lokalise

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/lokalise/go-lokalise-api/model"
)

const (
	pathScreenshots = "screenshots"
)

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

func (c *ScreenshotsService) List(ctx context.Context, pageOptions ScreenshotsOptions) (model.ScreenshotResponse, error) {
	var result model.ScreenshotResponse
	resp, err := c.client.getList(ctx, pathScreenshots, &result, pageOptions)
	if err != nil {
		return result, err
	}

	applyPaged(resp, &result.Paged)
	return result, apiError(resp)
}

func (c *ScreenshotsService) Create(ctx context.Context, projectID string, options CreateScreenshotOptions) (model.ScreenshotResponse, error) {
	var res model.ScreenshotResponse
	body, err := json.Marshal(options)
	resp, err := c.client.post(ctx, fmt.Sprintf("%s/%s", projectID, pathScreenshots), &res, string(body))
	if err != nil {
		return res, err
	}

	return res, apiError(resp)
}
