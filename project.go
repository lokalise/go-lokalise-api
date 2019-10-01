package lokalise

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

type Project struct {
	WithProjectID
	WithTeamID
	WithCreationTime
	Name           string `json:"name,omitempty"`
	Description    string `json:"description,omitempty"`
	CreatedBy      int64  `json:"created_by,omitempty"`
	CreatedByEmail string `json:"created_by_email,omitempty"`
}

type ProjectsResponse struct {
	Paged
	Projects []Project `json:"projects,omitempty"`
}

type ProjectEmptyResponse struct {
	WithProjectID
	KeysDeleted bool `json:"keys_deleted,omitempty"`
}

type ProjectDeleteResponse struct {
	WithProjectID
	Deleted bool `json:"project_deleted,omitempty"`
}

type ProjectsService struct {
	BaseService
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

// Methods implementations

func (c *ProjectsService) List(pageOptions ProjectsOptions) (r ProjectsResponse, err error) {
	resp, err := c.getList(c.Ctx(), pathProjects, &r, pageOptions)

	if err != nil {
		return
	}
	applyPaged(resp, &r.Paged)
	return r, apiError(resp)
}

func (c *ProjectsService) Create(name, description string) (r Project, err error) {
	resp, err := c.post(c.Ctx(), pathProjects, &r, map[string]interface{}{
		"name":        name,
		"description": description,
	})

	if err != nil {
		return
	}
	return r, apiError(resp)
}

func (c *ProjectsService) CreateForTeam(name, description string, teamID int64) (r Project, err error) {
	resp, err := c.post(c.Ctx(), pathProjects, &r, map[string]interface{}{
		"name":        name,
		"description": description,
		"team_id":     teamID,
	})

	if err != nil {
		return
	}
	return r, apiError(resp)
}

func (c *ProjectsService) Retrieve(projectID string) (r Project, err error) {
	resp, err := c.get(c.Ctx(), fmt.Sprintf("%s/%s", pathProjects, projectID), &r)

	if err != nil {
		return
	}
	return r, apiError(resp)
}

func (c *ProjectsService) Update(projectID, name, description string) (r Project, err error) {
	resp, err := c.put(c.Ctx(), fmt.Sprintf("%s/%s", pathProjects, projectID), &r, map[string]interface{}{
		"name":        name,
		"description": description,
	})

	if err != nil {
		return
	}
	return r, apiError(resp)
}

func (c *ProjectsService) Empty(projectID string) (r ProjectEmptyResponse, err error) {
	resp, err := c.put(c.Ctx(), fmt.Sprintf("%s/%s/empty", pathProjects, projectID), &r, nil)

	if err != nil {
		return
	}
	return r, apiError(resp)
}

func (c *ProjectsService) Delete(projectID string) (r ProjectDeleteResponse, err error) {
	resp, err := c.delete(c.Ctx(), fmt.Sprintf("%s/%s", pathProjects, projectID), &r)

	if err != nil {
		return
	}
	return r, apiError(resp)
}
