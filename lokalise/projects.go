package lokalise

import (
	"context"

	"github.com/lokalise/lokalise-go-sdk/model"
)

type ProjectsService struct {
	client *Client
}

const (
	pathProjects = "projects"
)

func (c *ProjectsService) List(ctx context.Context, pageOptions PageOptions) (model.ProjectsResponse, error) {
	var res model.ProjectsResponse
	resp, err := c.client.getList(ctx, pathProjects, &res, pageOptions)
	if err != nil {
		return model.ProjectsResponse{}, err
	}
	applyPaged(resp, &res.Paged)
	return res, apiError(resp)
}
