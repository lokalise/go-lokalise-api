package lokalise

import (
	"path"
	"strconv"
)

const (
	pathSnapshots = "snapshots"
)

type SnapshotService struct {
	BaseService
}

// ‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// Service entity objects
// _____________________________________________________________________________________________________________________

type Snapshot struct {
	WithCreationTime
	WithCreationUser

	SnapshotID int64  `json:"snapshot_id"`
	Title      string `json:"title"`
}

// ‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// Service request/response objects
// _____________________________________________________________________________________________________________________

type ListSnapshotsResponse struct {
	Paged
	WithProjectID
	Snapshots []Snapshot `json:"snapshots"`
}

type CreateSnapshotResponse struct {
	WithProjectID
	Snapshot Snapshot `json:"snapshot"`
}

type DeleteSnapshotResponse struct {
	WithProjectID
	SnapshotDeleted bool `json:"snapshot_deleted"`
}

// ‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// Service methods
// _____________________________________________________________________________________________________________________

func (c *SnapshotService) List(projectID string) (r ListSnapshotsResponse, err error) {
	path := path.Join(pathProjects, projectID, pathSnapshots)
	resp, err := c.getList(c.Ctx(), path, &r, c.PageOpts())

	if err != nil {
		return
	}
	applyPaged(resp, &r.Paged)
	return r, apiError(resp)
}

func (c *SnapshotService) Create(projectID string, title string) (r CreateSnapshotResponse, err error) {
	path := path.Join(pathProjects, projectID, pathSnapshots)
	resp, err := c.post(c.Ctx(), path, &r, map[string]interface{}{"title": title})

	if err != nil {
		return
	}
	return r, apiError(resp)
}

func (c *SnapshotService) Delete(projectID string, ID int64) (r DeleteSnapshotResponse, err error) {
	path := path.Join(pathProjects, projectID, pathSnapshots, strconv.FormatInt(ID, 10))
	resp, err := c.delete(c.Ctx(), path, &r)

	if err != nil {
		return
	}
	return r, apiError(resp)
}

func (c *SnapshotService) Restore(projectID string, ID int64) (r Project, err error) {
	path := path.Join(pathProjects, projectID, pathSnapshots, strconv.FormatInt(ID, 10))
	resp, err := c.post(c.Ctx(), path, &r, nil)

	if err != nil {
		return
	}
	return r, apiError(resp)
}
