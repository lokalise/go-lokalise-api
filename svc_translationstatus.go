package lokalise

import (
	"fmt"
)

const (
	pathTranslationStatuses = "custom_translation_statuses"
)

type TranslationStatusesService struct {
	BaseService
}

type TranslationStatus struct {
	StatusID int64  `json:"status_id"`
	Title    string `json:"title"`
	Color    string `json:"color"`
}

type TranslationStatusesResponse struct {
	Paged
	WithProjectID
	TranslationStatuses []TranslationStatus `json:"custom_translation_statuses"`
}

type TranslationStatusColorsResponse struct {
	Colors []string `json:"colors"`
}

type TranslationStatusResponse struct {
	WithProjectID
	TranslationStatus TranslationStatus `json:"custom_translation_status"`
}

type TranslationStatusDeleteResponse struct {
	WithProjectID
	Deleted bool `json:"deleted"`
}

type CreateTranslationStatus struct {
	Title string `json:"title"`
	Color string `json:"color"`
}

type UpdateTranslationStatus struct {
	Title string `json:"title"`
	Color string `json:"color"`
}

func (c *TranslationStatusesService) List(projectID string) (r TranslationStatusesResponse, err error) {
	resp, err := c.getList(c.Ctx(), fmt.Sprintf("%s/%s/%s", pathProjects, projectID, pathTranslationStatuses), &r, c.PageOpts())

	if err != nil {
		return
	}
	applyPaged(resp, &r.Paged)
	return r, apiError(resp)
}

func (c *TranslationStatusesService) ListColors(projectID string) (r TranslationStatusColorsResponse, err error) {
	resp, err := c.get(c.Ctx(), fmt.Sprintf("%s/%s/%s/%s", pathProjects, projectID, pathTranslationStatuses, "colors"), &r)

	if err != nil {
		return
	}
	return r, apiError(resp)
}

func (c *TranslationStatusesService) Create(projectID string, options CreateTranslationStatus) (r TranslationStatusResponse, err error) {
	resp, err := c.post(c.Ctx(), fmt.Sprintf("%s/%s/%s", pathProjects, projectID, pathTranslationStatuses), &r, options)

	if err != nil {
		return
	}
	return r, apiError(resp)
}

func (c *TranslationStatusesService) Retrieve(projectID string, statusID int64) (r TranslationStatusResponse, err error) {
	resp, err := c.get(c.Ctx(), fmt.Sprintf("%s/%s/%s/%d", pathProjects, projectID, pathTranslationStatuses, statusID), &r)

	if err != nil {
		return
	}
	return r, apiError(resp)
}

func (c *TranslationStatusesService) Update(projectID string, statusID int64, opts UpdateTranslationStatus) (r TranslationStatusResponse, err error) {
	resp, err := c.put(c.Ctx(), fmt.Sprintf("%s/%s/%s/%d", pathProjects, projectID, pathTranslationStatuses, statusID), &r, opts)

	if err != nil {
		return
	}
	return r, apiError(resp)
}

func (c *TranslationStatusesService) Delete(projectID string, statusID int64) (r TranslationStatusDeleteResponse, err error) {
	resp, err := c.delete(c.Ctx(), fmt.Sprintf("%s/%s/%s/%d", pathProjects, projectID, pathTranslationStatuses, statusID), &r)

	if err != nil {
		return
	}
	return r, apiError(resp)
}
