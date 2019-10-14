package lokalise

import (
	"fmt"
)

type ContributorService struct {
	BaseService
}

// ‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// Service entity objects
// _____________________________________________________________________________________________________________________

type Contributor struct {
	WithCreationTime
	WithUserID

	Email    string `json:"email"`
	Fullname string `json:"fullname"`
	Permission
}

// ‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// Service request/response objects
// _____________________________________________________________________________________________________________________

type NewContributor struct {
	Email    string `json:"email"`
	Fullname string `json:"fullname,omitempty"`
	Permission
}

type ContributorsResponse struct {
	Paged
	WithProjectID
	Contributors []Contributor `json:"contributors"`
}

type ContributorResponse struct {
	WithProjectID
	Contributor Contributor `json:"contributor"`
}

type DeleteContributorResponse struct {
	WithProjectID
	IsDeleted bool `json:"contributor_deleted"`
}

// ‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// Service methods
// _____________________________________________________________________________________________________________________

func (c *ContributorService) List(projectID string) (r ContributorsResponse, err error) {
	resp, err := c.getList(c.Ctx(), pathContributors(projectID), &r, c.PageOpts())

	if err != nil {
		return
	}
	applyPaged(resp, &r.Paged)
	return r, apiError(resp)
}

func (c *ContributorService) Create(projectID string, cs []NewContributor) (r ContributorsResponse, err error) {
	resp, err := c.post(c.Ctx(), pathContributors(projectID), &r, map[string]interface{}{"contributors": cs})

	if err != nil {
		return
	}
	return r, apiError(resp)
}

func (c *ContributorService) Retrieve(projectID string, userID int64) (r ContributorResponse, err error) {
	resp, err := c.get(c.Ctx(), pathContributorByID(projectID, userID), &r)

	if err != nil {
		return
	}
	return r, apiError(resp)
}

func (c *ContributorService) Update(projectID string, userID int64, p Permission) (r ContributorResponse, err error) {
	resp, err := c.put(c.Ctx(), pathContributorByID(projectID, userID), &r, p)

	if err != nil {
		return
	}
	return r, apiError(resp)
}

func (c *ContributorService) Delete(projectID string, userID int64) (r DeleteContributorResponse, err error) {
	resp, err := c.delete(c.Ctx(), pathContributorByID(projectID, userID), &r)

	if err != nil {
		return
	}
	return r, apiError(resp)
}

func pathContributors(projectID string) string {
	return fmt.Sprintf("%s/%s/contributors", pathProjects, projectID)
}

func pathContributorByID(projectID string, userID int64) string {
	return fmt.Sprintf("%s/%s/contributors/%d", pathProjects, projectID, userID)
}
