package lokalise

import (
	"context"
	"fmt"

	"github.com/go-resty/resty"
	"github.com/lokalise/lokalise-go-sdk/model"
)

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

func (c *ProjectsService) List(ctx context.Context, pageOptions ProjectsOptions) (model.ProjectsResponse, error) {
	var res model.ProjectsResponse
	resp, err := c.client.getList(ctx, pathProjects, &res, pageOptions)
	if err != nil {
		return model.ProjectsResponse{}, err
	}
	applyPaged(resp, &res.Paged)
	return res, apiError(resp)
}
