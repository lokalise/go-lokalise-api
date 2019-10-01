package lokalise

import (
	"path"
	"strconv"
)

const (
	pathSnapshots = "snapshots"
)

type Snapshot struct {
	WithCreationTime
	SnapshotID     int64  `json:"snapshot_id,omitempty"`
	Title          string `json:"title,omitempty"`
	CreatedBy      int64  `json:"created_by,omitempty"`
	CreatedByEmail string `json:"created_by_email,omitempty"`
}

type ListSnapshotsResponse struct {
	Paged
	WithProjectID
	Snapshots []Snapshot `json:"snapshots,omitempty"`
}

type CreateSnapshotResponse struct {
	WithProjectID
	Snapshot Snapshot `json:"snapshot,omitempty"`
}

type DeleteSnapshotResponse struct {
	WithProjectID
	SnapshotDeleted bool `json:"snapshot_deleted,omitempty"`
}

type SnapshotsService struct {
	BaseService
}

func (c *SnapshotsService) List(projectID string) (r ListSnapshotsResponse, err error) {
	path := path.Join(pathProjects, projectID, pathSnapshots)
	resp, err := c.getList(c.Ctx(), path, &r, c.PageOpts())

	if err != nil {
		return
	}
	applyPaged(resp, &r.Paged)
	return r, apiError(resp)
}

func (c *SnapshotsService) Create(projectID string, title string) (r CreateSnapshotResponse, err error) {
	path := path.Join(pathProjects, projectID, pathSnapshots)
	resp, err := c.post(c.Ctx(), path, &r, map[string]interface{}{"title": title})

	if err != nil {
		return
	}
	return r, apiError(resp)
}

func (c *SnapshotsService) Delete(projectID string, ID int64) (r DeleteSnapshotResponse, err error) {
	path := path.Join(pathProjects, projectID, pathSnapshots, strconv.FormatInt(ID, 10))
	resp, err := c.delete(c.Ctx(), path, &r)

	if err != nil {
		return
	}
	return r, apiError(resp)
}

func (c *SnapshotsService) Restore(projectID string, ID int64) (r Project, err error) {
	path := path.Join(pathProjects, projectID, pathSnapshots, strconv.FormatInt(ID, 10))
	resp, err := c.post(c.Ctx(), path, &r, nil)

	if err != nil {
		return
	}
	return r, apiError(resp)
}
