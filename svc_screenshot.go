package lokalise

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

const (
	pathScreenshots = "screenshots"
)

type Screenshot struct {
	WithCreationTime
	ScreenshotID   int64    `json:"screenshot_id,omitempty"`
	KeyIDs         []int64  `json:"key_ids,omitempty"`
	URL            string   `json:"url,omitempty"`
	Title          string   `json:"title,omitempty"`
	Description    string   `json:"description,omitempty"`
	ScreenshotTags []string `json:"screenshot_tags,omitempty"`
	Width          int64    `json:"width,omitempty"`
	Height         int64    `json:"height,omitempty"`
}

type ScreenshotResponse struct {
	Paged
	Screenshots []Screenshot `json:"screenshots,omitempty"`
}

type ScreenshotDeleteResponse struct {
	WithProjectID
	Deleted bool `json:"screenshot_deleted"`
}

type ScreenshotsService struct {
	BaseService
}

type ScreenshotsOptions struct {
	PageOptions
	IncludeTags bool
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

type UpdateScreenshotOptions struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
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

func (c *ScreenshotsService) List(pageOptions ScreenshotsOptions) (r ScreenshotResponse, err error) {
	resp, err := c.getList(c.Ctx(), pathScreenshots, &r, pageOptions)

	if err != nil {
		return
	}
	applyPaged(resp, &r.Paged)
	return r, apiError(resp)
}

func (c *ScreenshotsService) Create(projectID string, options CreateScreenshotOptions) (r ScreenshotResponse, err error) {
	resp, err := c.post(c.Ctx(), fmt.Sprintf("%s/%s/%s", pathProjects, projectID, pathScreenshots), &r, options)

	if err != nil {
		return
	}
	return r, apiError(resp)
}

func (c *ScreenshotsService) Retrieve(projectID string, screenshotID int64) (r ScreenshotResponse, err error) {
	resp, err := c.get(c.Ctx(), fmt.Sprintf("%s/%s/%s/%d", pathProjects, projectID, pathScreenshots, screenshotID), &r)

	if err != nil {
		return
	}
	return r, apiError(resp)
}

func (c *ScreenshotsService) Update(projectID string, screenshotID int64, opts UpdateScreenshotOptions) (r ScreenshotResponse, err error) {
	resp, err := c.put(c.Ctx(), fmt.Sprintf("%s/%s/%s/%d", pathProjects, projectID, pathScreenshots, screenshotID), &r, opts)

	if err != nil {
		return
	}
	return r, apiError(resp)
}

func (c *ScreenshotsService) Delete(projectID string, screenshotID int64) (r ScreenshotDeleteResponse, err error) {
	resp, err := c.delete(c.Ctx(), fmt.Sprintf("%s/%s/%s/%d", pathProjects, projectID, pathScreenshots, screenshotID), &r)

	if err != nil {
		return
	}
	return r, apiError(resp)
}
