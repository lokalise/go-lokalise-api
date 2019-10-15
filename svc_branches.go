package lokalise

import (
	"path"
	"strconv"
)

const (
	pathBranches = "branches"
)

type BranchService struct {
	BaseService
}

// ‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// Service entity objects
// _____________________________________________________________________________________________________________________

type Branch struct {
	WithCreationTime
	WithCreationUser

	BranchID int64  `json:"branch_id"`
	Name     string `json:"name"`
}

// ‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// Service request/response objects
// _____________________________________________________________________________________________________________________

type ListBranchesResponse struct {
	Paged
	WithProjectID
	Branches []Branch `json:"branches"`
}

type CreateBranchResponse struct {
	WithProjectID
	Branch Branch `json:"branch"`
}

type DeleteBranchResponse struct {
	WithProjectID
	BranchDeleted bool `json:"branch_deleted"`
}

// ‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// Service methods
// _____________________________________________________________________________________________________________________

func (c *BranchService) List(projectID string) (r ListBranchesResponse, err error) {
	path := path.Join(pathProjects, projectID, pathBranches)
	resp, err := c.getList(c.Ctx(), path, &r, c.PageOpts())

	if err != nil {
		return
	}
	applyPaged(resp, &r.Paged)
	return r, apiError(resp)
}

func (c *BranchService) Create(projectID string, name string) (r CreateBranchResponse, err error) {
	path := path.Join(pathProjects, projectID, pathBranches)
	resp, err := c.post(c.Ctx(), path, &r, map[string]interface{}{"name": name})

	if err != nil {
		return
	}
	return r, apiError(resp)
}

func (c *BranchService) Delete(projectID string, ID int64) (r DeleteBranchResponse, err error) {
	path := path.Join(pathProjects, projectID, pathBranches, strconv.FormatInt(ID, 10))
	resp, err := c.delete(c.Ctx(), path, &r)

	if err != nil {
		return
	}
	return r, apiError(resp)
}
