package lokalise

import (
	"fmt"
)

const (
	pathTranslationStatuses = "custom_translation_statuses"
)

type TranslationStatusService struct {
	BaseService
}

// ‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// Service entity objects
// _____________________________________________________________________________________________________________________

type TranslationStatus struct {
	StatusID int64  `json:"status_id"` // todo check with `json:"id"` for svc_translation
	Title    string `json:"title"`
	Color    string `json:"color"`
}

// ‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// Service request/response objects
// _____________________________________________________________________________________________________________________

type TranslationStatusesResponse struct {
	Paged
	WithProjectID
	TranslationStatuses []TranslationStatus `json:"custom_translation_statuses"`
}

type ListColorsTranslationStatusResponse struct {
	Colors []string `json:"colors"`
}

type TranslationStatusResponse struct {
	WithProjectID
	TranslationStatus TranslationStatus `json:"custom_translation_status"`
}

type DeleteTranslationStatusResponse struct {
	WithProjectID
	Deleted bool `json:"deleted"`
}

type NewTranslationStatus struct {
	Title string `json:"title"`
	Color string `json:"color"`
}

type UpdateTranslationStatus struct {
	Title string `json:"title,omitempty"`
	Color string `json:"color,omitempty"`
}

// ‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// Service methods
// _____________________________________________________________________________________________________________________

func (c *TranslationStatusService) List(projectID string) (r TranslationStatusesResponse, err error) {
	resp, err := c.getList(c.Ctx(), fmt.Sprintf("%s/%s/%s", pathProjects, projectID, pathTranslationStatuses), &r, c.PageOpts())

	if err != nil {
		return
	}
	applyPaged(resp, &r.Paged)
	return r, apiError(resp)
}

func (c *TranslationStatusService) ListColors(projectID string) (r ListColorsTranslationStatusResponse, err error) {
	resp, err := c.get(c.Ctx(), fmt.Sprintf("%s/%s/%s/%s", pathProjects, projectID, pathTranslationStatuses, "colors"), &r)

	if err != nil {
		return
	}
	return r, apiError(resp)
}

func (c *TranslationStatusService) Create(projectID string, options NewTranslationStatus) (r TranslationStatusResponse, err error) {
	resp, err := c.post(c.Ctx(), fmt.Sprintf("%s/%s/%s", pathProjects, projectID, pathTranslationStatuses), &r, options)

	if err != nil {
		return
	}
	return r, apiError(resp)
}

func (c *TranslationStatusService) Retrieve(projectID string, statusID int64) (r TranslationStatusResponse, err error) {
	resp, err := c.get(c.Ctx(), fmt.Sprintf("%s/%s/%s/%d", pathProjects, projectID, pathTranslationStatuses, statusID), &r)

	if err != nil {
		return
	}
	return r, apiError(resp)
}

func (c *TranslationStatusService) Update(projectID string, statusID int64, opts UpdateTranslationStatus) (r TranslationStatusResponse, err error) {
	resp, err := c.put(c.Ctx(), fmt.Sprintf("%s/%s/%s/%d", pathProjects, projectID, pathTranslationStatuses, statusID), &r, opts)

	if err != nil {
		return
	}
	return r, apiError(resp)
}

func (c *TranslationStatusService) Delete(projectID string, statusID int64) (r DeleteTranslationStatusResponse, err error) {
	resp, err := c.delete(c.Ctx(), fmt.Sprintf("%s/%s/%s/%d", pathProjects, projectID, pathTranslationStatuses, statusID), &r)

	if err != nil {
		return
	}
	return r, apiError(resp)
}
