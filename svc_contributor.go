package lokalise

import (
	"fmt"
	"path"
	"strconv"
)

type Contributor struct {
	WithCreationTime
	WithUserID
	Email    string `json:"email"`
	Fullname string `json:"fullname"`
	Permission
}

func pathContributors(projectID string) string {
	return fmt.Sprintf("%s/%s/contributors", pathProjects, projectID)
}

type CustomContributor struct {
	Fullname string `json:"fullname,omitempty"`
	Email    string `json:"email"`
	Permission
}

type ContributorsResponse struct { // maybe rename to avoid messing
	Paged
	ProjectID    string        `json:"project_id"`
	Contributors []Contributor `json:"contributors"`
}

type ContributorResponse struct {
	ProjectID   string      `json:"project_id"`
	Contributor Contributor `json:"contributor"`
}

type DeleteContributorResponse struct {
	ProjectID string `json:"project_id"`
	IsDeleted bool   `json:"contributor_deleted"`
}

type ContributorsService struct {
	BaseService
}

// Method implementations

func (c *ContributorsService) List(projectID string) (r ContributorsResponse, err error) {
	resp, err := c.getList(c.Ctx(), pathContributors(projectID), &r, c.PageOpts())

	if err != nil {
		return
	}
	applyPaged(resp, &r.Paged)
	return r, apiError(resp)
}

func (c *ContributorsService) Create(projectID string, cs []CustomContributor) (r ContributorResponse, err error) {
	resp, err := c.post(c.Ctx(), pathContributors(projectID), &r, map[string]interface{}{"contributors": cs})

	if err != nil {
		return
	}
	return r, apiError(resp)
}

func (c *ContributorsService) Retrieve(projectID string, userID int64) (r ContributorResponse, err error) {
	url := path.Join(pathContributors(projectID), strconv.FormatInt(userID, 10))
	resp, err := c.get(c.Ctx(), url, &r)

	if err != nil {
		return
	}
	return r, apiError(resp)
}

func (c *ContributorsService) Update(projectID string, userID int64, p Permission) (r ContributorResponse, err error) {
	url := path.Join(pathContributors(projectID), strconv.FormatInt(userID, 10))
	resp, err := c.put(c.Ctx(), url, &r, p)

	if err != nil {
		return
	}
	return r, apiError(resp)
}

func (c *ContributorsService) Delete(projectID string, userID int64) (r DeleteContributorResponse, err error) {
	url := path.Join(pathContributors(projectID), strconv.FormatInt(userID, 10))
	resp, err := c.delete(c.Ctx(), url, &r)

	if err != nil {
		return
	}
	return r, apiError(resp)
}
