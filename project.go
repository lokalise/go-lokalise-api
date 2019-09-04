package lokalise

import (
	"context"
	"fmt"

	"github.com/go-resty/resty/v2"
)

type Project struct {
	ProjectID      string `json:"project_id,omitempty"`
	Name           string `json:"name,omitempty"`
	Description    string `json:"description,omitempty"`
	CreatedAt      string `json:"created_at,omitempty"`
	CreatedBy      int64  `json:"created_by,omitempty"`
	CreatedByEmail string `json:"created_by_email,omitempty"`
	TeamID         int64  `json:"team_id,omitempty"`
}

type ProjectsResponse struct {
	Paged
	Projects []Project `json:"projects,omitempty"`
}

type ProjectEmptyResponse struct {
	ProjectID   string `json:"project_id,omitempty"`
	KeysDeleted bool   `json:"keys_deleted,omitempty"`
}

type ProjectDeleteResponse struct {
	ProjectID string `json:"project_id,omitempty"`
	Deleted   bool   `json:"project_deleted,omitempty"`
}

type ProjectsService struct {
	client *Client
}

const (
	pathProjects = "projects"
)

type ProjectsOptions struct {
	PageOptions
	TeamID int64 `json:"teamID,omitempty"`
}

func (options ProjectsOptions) Apply(req *resty.Request) {
	options.PageOptions.Apply(req)
	if options.TeamID != 0 {
		req.SetQueryParam("filter_team_id", fmt.Sprintf("%d", options.TeamID))
	}
}

func (c *ProjectsService) List(ctx context.Context, pageOptions ProjectsOptions) (ProjectsResponse, error) {
	var res ProjectsResponse
	resp, err := c.client.getList(ctx, pathProjects, &res, pageOptions)
	if err != nil {
		return ProjectsResponse{}, err
	}
	applyPaged(resp, &res.Paged)
	return res, apiError(resp)
}

func (c *ProjectsService) Create(ctx context.Context, name, description string) (Project, error) {
	var res Project
	resp, err := c.client.post(ctx, pathProjects, &res, map[string]interface{}{
		"name":        name,
		"description": description,
	})
	if err != nil {
		return Project{}, err
	}
	return res, apiError(resp)
}

func (c *ProjectsService) CreateForTeam(ctx context.Context, name, description string, teamID int64) (Project, error) {
	var res Project
	resp, err := c.client.post(ctx, pathProjects, &res, map[string]interface{}{
		"name":        name,
		"description": description,
		"team_id":     teamID,
	})
	if err != nil {
		return Project{}, err
	}
	return res, apiError(resp)
}

func (c *ProjectsService) Retrieve(ctx context.Context, projectID string) (Project, error) {
	var res Project
	resp, err := c.client.get(ctx, fmt.Sprintf("%s/%s", pathProjects, projectID), &res)
	if err != nil {
		return Project{}, err
	}
	return res, apiError(resp)
}

func (c *ProjectsService) Update(ctx context.Context, projectID, name, description string) (Project, error) {
	var res Project
	resp, err := c.client.put(ctx, fmt.Sprintf("%s/%s", pathProjects, projectID), &res, map[string]interface{}{
		"name":        name,
		"description": description,
	})
	if err != nil {
		return Project{}, err
	}
	return res, apiError(resp)
}

func (c *ProjectsService) Empty(ctx context.Context, projectID string) (ProjectEmptyResponse, error) {
	var res ProjectEmptyResponse
	resp, err := c.client.put(ctx, fmt.Sprintf("%s/%s/empty", pathProjects, projectID), &res, nil)
	if err != nil {
		return ProjectEmptyResponse{}, err
	}
	return res, apiError(resp)
}

func (c *ProjectsService) Delete(ctx context.Context, projectID string) (ProjectDeleteResponse, error) {
	var res ProjectDeleteResponse
	resp, err := c.client.delete(ctx, fmt.Sprintf("%s/%s", pathProjects, projectID), &res)
	if err != nil {
		return ProjectDeleteResponse{}, err
	}
	return res, apiError(resp)
}
