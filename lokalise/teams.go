package lokalise

import (
	"context"

	"github.com/17media/go-lokalise-api/model"
)

type TeamsService struct {
	client *Client
}

const (
	pathTeams = "teams"
)

func (c *TeamsService) List(ctx context.Context, pageOptions PageOptions) (model.TeamsResponse, error) {
	var res model.TeamsResponse
	resp, err := c.client.getList(ctx, pathTeams, &res, &pageOptions)
	if err != nil {
		return model.TeamsResponse{}, err
	}
	applyPaged(resp, &res.Paged)
	return res, apiError(resp)
}
