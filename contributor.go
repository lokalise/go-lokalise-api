package lokalise

import (
	"context"
	"fmt"
	"path"
	"strconv"
)

type Contributor struct {
	UserId      int64  `json:"user_id"` // todo usage pointer with omitempty option should allow to leave only one struct (inst of CustomContributor)
	Email       string `json:"email"`
	Fullname    string `json:"fullname"`
	CreatedAt   string `json:"created_at"` // string or time.Time? todo move to CreatedAt struct
	CreatedAtTs int64  `json:"created_at_timestamp"`
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
	client *Client
}

// Method implementations

func (c *ContributorsService) List(ctx context.Context, projectID string, pageOptions PageOptions) (ContributorsResponse, error) {
	var res ContributorsResponse
	resp, err := c.client.getList(ctx, pathContributors(projectID), &res, &pageOptions)

	if err != nil {
		return ContributorsResponse{}, err // todo dont create everywhere a new object, res already exists
	}
	applyPaged(resp, &res.Paged)
	return res, apiError(resp)
}

func (c *ContributorsService) Create(
	ctx context.Context,
	projectID string,
	contributors []CustomContributor,
) (result ContributorResponse, err error) {
	resp, err := c.client.post(ctx, pathContributors(projectID), &result, map[string]interface{}{
		"contributors": contributors,
	})

	if err != nil {
		return result, err
	}
	return result, apiError(resp)
}

func (c *ContributorsService) Retrieve(ctx context.Context, projectID string, userID int64) (result ContributorResponse, err error) {
	url := path.Join(pathContributors(projectID), strconv.FormatInt(userID, 10))
	resp, err := c.client.get(ctx, url, &result)

	if err != nil {
		return result, err
	}
	return result, apiError(resp)
}

func (c *ContributorsService) Update(ctx context.Context, projectID string, userID int64, p Permission) (result ContributorResponse, err error) {
	url := path.Join(pathContributors(projectID), strconv.FormatInt(userID, 10))
	resp, err := c.client.put(ctx, url, &result, p)

	if err != nil {
		return result, err
	}
	return result, apiError(resp)
}

func (c *ContributorsService) Delete(ctx context.Context, projectID string, userID int64) (result DeleteContributorResponse, err error) {
	url := path.Join(pathContributors(projectID), strconv.FormatInt(userID, 10))
	resp, err := c.client.delete(ctx, url, &result)

	if err != nil {
		return result, err
	}
	return result, apiError(resp)
}
