package lokalise

import (
	"context"
	"path"
	"strconv"

	"github.com/go-resty/resty/v2"
)

const (
	pathSnapshots = "snapshots"
)

type Snapshot struct {
	SnapshotID     int64  `json:"snapshot_id,omitempty"`
	Title          string `json:"title,omitempty"`
	CreatedAt      string `json:"created_at,omitempty"`
	CreatedAtTs    int64  `json:"created_at_timestamp,omitempty"`
	CreatedBy      int64  `json:"created_by,omitempty"`
	CreatedByEmail string `json:"created_by_email,omitempty"`
}

type ListSnapshotsResponse struct {
	Paged
	ProjectID string     `json:"project_id,omitempty"`
	Snapshots []Snapshot `json:"snapshots,omitempty"`
}

type CreateSnapshotResponse struct {
	ProjectID string   `json:"project_id,omitempty"`
	Snapshot  Snapshot `json:"snapshot,omitempty"`
}

type DeleteSnapshotResponse struct {
	ProjectID       string `json:"project_id,omitempty"`
	SnapshotDeleted bool   `json:"snapshot_deleted,omitempty"`
}

type SnapshotsService struct {
	client *Client
}

type SnapshotsOptions struct {
	PageOptions
}

func (options SnapshotsOptions) Apply(req *resty.Request) {
	options.PageOptions.Apply(req)
}

func (c *SnapshotsService) List(ctx context.Context, projectID string, pageOptions SnapshotsOptions) (result ListSnapshotsResponse, err error) {
	path := path.Join(pathProjects, projectID, pathSnapshots)
	resp, err := c.client.getList(ctx, path, &result, pageOptions)
	if err != nil {
		return result, err
	}

	applyPaged(resp, &result.Paged)
	return result, apiError(resp)
}

func (c *SnapshotsService) Create(ctx context.Context, projectID string, title string) (result CreateSnapshotResponse, err error) {
	path := path.Join(pathProjects, projectID, pathSnapshots)
	resp, err := c.client.post(ctx, path, &result, map[string]interface{}{"title": title})
	if err != nil {
		return result, err
	}
	return result, apiError(resp)
}

func (c *SnapshotsService) Delete(ctx context.Context, projectID string, ID int64) (result DeleteSnapshotResponse, err error) {
	path := path.Join(pathProjects, projectID, pathSnapshots, strconv.FormatInt(ID, 10))
	resp, err := c.client.delete(ctx, path, &result)
	if err != nil {
		return result, err
	}
	return result, apiError(resp)
}

func (c *SnapshotsService) Restore(ctx context.Context, projectID string, ID int64) (result Project, err error) {
	path := path.Join(pathProjects, projectID, pathSnapshots, strconv.FormatInt(ID, 10))
	resp, err := c.client.post(ctx, path, &result, nil)
	if err != nil {
		return result, err
	}
	return result, apiError(resp)
}
