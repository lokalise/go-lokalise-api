package lokalise

import (
	"fmt"
	"github.com/google/go-querystring/query"

	"github.com/go-resty/resty/v2"
)

const (
	pathScreenshots = "screenshots"
)

type ScreenshotService struct {
	BaseService

	listOpts ScreenshotListOptions
}

// ‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// Service entity objects
// _____________________________________________________________________________________________________________________

type Screenshot struct {
	WithCreationTime

	ScreenshotID   int64    `json:"screenshot_id"`
	KeyIDs         []int64  `json:"key_ids"`
	URL            string   `json:"url"`
	Title          string   `json:"title"`
	Description    string   `json:"description"`
	ScreenshotTags []string `json:"tags"`
	Width          int64    `json:"width"`
	Height         int64    `json:"height"`
}

// ‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// Service request/response objects
// _____________________________________________________________________________________________________________________

type ScreenshotResponse struct {
	WithProjectID
	Screenshot Screenshot `json:"screenshot"`
}

type ScreenshotsResponse struct {
	Paged
	WithProjectID
	Screenshots []Screenshot `json:"screenshots"`
}

type DeleteScreenshotResponse struct {
	WithProjectID
	Deleted bool `json:"screenshot_deleted"`
}

type NewScreenshot struct {
	// The screenshot, base64 encoded (with leading image type `data:image/jpeg;base64,`).
	// Supported file formats are JPG and PNG.
	Body        string   `json:"data"` // maybe []byte?
	Title       string   `json:"title,omitempty"`
	Description string   `json:"description,omitempty"`
	Ocr         *bool    `json:"ocr,omitempty"`
	KeyIDs      []int64  `json:"key_ids,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}

type UpdateScreenshot struct {
	Title       string   `json:"title,omitempty"`
	Description string   `json:"description,omitempty"`
	KeyIDs      []int64  `json:"key_ids,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}

func (c *ScreenshotService) List(projectID string) (r ScreenshotsResponse, err error) {
	resp, err := c.getWithOptions(c.Ctx(), fmt.Sprintf("%s/%s/%s", pathProjects, projectID, pathScreenshots), &r, c.ListOpts())

	if err != nil {
		return
	}
	applyPaged(resp, &r.Paged)
	return r, apiError(resp)
}

func (c *ScreenshotService) Create(projectID string, screenshots []NewScreenshot) (r ScreenshotsResponse, err error) {
	resp, err := c.post(c.Ctx(), fmt.Sprintf("%s/%s/%s", pathProjects, projectID, pathScreenshots), &r,
		map[string]interface{}{
			"screenshots": screenshots,
		},
	)

	if err != nil {
		return
	}
	return r, apiError(resp)
}

func (c *ScreenshotService) Retrieve(projectID string, screenshotID int64) (r ScreenshotResponse, err error) {
	resp, err := c.get(c.Ctx(), fmt.Sprintf("%s/%s/%s/%d", pathProjects, projectID, pathScreenshots, screenshotID), &r)

	if err != nil {
		return
	}
	return r, apiError(resp)
}

func (c *ScreenshotService) Update(projectID string, screenshotID int64, opts UpdateScreenshot) (r ScreenshotResponse, err error) {
	resp, err := c.put(c.Ctx(), fmt.Sprintf("%s/%s/%s/%d", pathProjects, projectID, pathScreenshots, screenshotID), &r, opts)

	if err != nil {
		return
	}
	return r, apiError(resp)
}

func (c *ScreenshotService) Delete(projectID string, screenshotID int64) (r DeleteScreenshotResponse, err error) {
	resp, err := c.delete(c.Ctx(), fmt.Sprintf("%s/%s/%s/%d", pathProjects, projectID, pathScreenshots, screenshotID), &r)

	if err != nil {
		return
	}
	return r, apiError(resp)
}

// ‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// Additional methods
// _____________________________________________________________________________________________________________________

type ScreenshotListOptions struct {
	// page options
	Page  uint `url:"page,omitempty"`
	Limit uint `url:"limit,omitempty"`

	IncludeTags uint8 `json:"include_tags,omitempty"`
	ListOnly    uint8 `json:"list_only,omitempty"`
}

func (options ScreenshotListOptions) Apply(req *resty.Request) {
	v, _ := query.Values(options)
	req.SetQueryString(v.Encode())
}

func (c *ScreenshotService) ListOpts() ScreenshotListOptions        { return c.listOpts }
func (c *ScreenshotService) SetListOptions(o ScreenshotListOptions) { c.listOpts = o }
func (c *ScreenshotService) WithListOptions(o ScreenshotListOptions) *ScreenshotService {
	c.listOpts = o
	return c
}
