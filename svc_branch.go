package lokalise

import (
	"path"
	"strconv"
)

const (
	pathBranches = "branches"
)

// BranchService handles communication with the branch related
// methods of the Lokalise API.
//
// Lokalise API docs: https://lokalise.com/api2docs/curl/#resource-branches
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

type CreateBranchRequest struct {
	Name string `json:"name"`
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

// List retrieves a list of branches
//
// Lokalise API docs: https://lokalise.com/api2docs/curl/#transition-list-all-branches-get
func (c *BranchService) List(projectID string) (r ListBranchesResponse, err error) {
	endpoint := path.Join(pathProjects, projectID, pathBranches)
	resp, err := c.getList(c.Ctx(), endpoint, &r, c.PageOpts())

	if err != nil {
		return
	}
	applyPaged(resp, &r.Paged)
	return r, apiError(resp)
}

// Create creates a branch in the project. Requires admin right.
//
// Lokalise API docs: https://lokalise.com/api2docs/curl/#transition-create-a-branch-post
func (c *BranchService) Create(projectID string, name string) (r CreateBranchResponse, err error) {
	endpoint := path.Join(pathProjects, projectID, pathBranches)
	resp, err := c.post(c.Ctx(), endpoint, &r, CreateBranchRequest{Name: name})

	if err != nil {
		return
	}
	return r, apiError(resp)
}

// Delete deletes a configured branch in the project. Requires admin right.
//
// Lokalise API docs: https://lokalise.com/api2docs/curl/#transition-delete-a-branch-delete
func (c *BranchService) Delete(projectID string, ID int64) (r DeleteBranchResponse, err error) {
	endpoint := path.Join(pathProjects, projectID, pathBranches, strconv.FormatInt(ID, 10))
	resp, err := c.delete(c.Ctx(), endpoint, &r)

	if err != nil {
		return
	}
	return r, apiError(resp)
}
