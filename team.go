package lokalise

import (
	"context"
)

type TeamsResponse struct {
	Paged
	Teams []Team `json:"teams,omitempty"`
}

type Team struct {
	TeamID       int64  `json:"team_id,omitempty"`
	Name         string `json:"name,omitempty"`
	CreatedAt    string `json:"created_at,omitempty"`
	Plan         string `json:"plan,omitempty"`
	QuotaUsage   Quota  `json:"quota_usage,omitempty"`
	QuotaAllowed Quota  `json:"quota_allowed,omitempty"`
}

type Quota struct {
	Users    int64 `json:"users,omitempty"`
	Keys     int64 `json:"keys,omitempty"`
	Projects int64 `json:"projects,omitempty"`
	MAU      int64 `json:"mau,omitempty"`
}

type TeamsService struct {
	client *Client
}

const (
	pathTeams = "teams"
)

func (c *TeamsService) List(ctx context.Context, pageOptions PageOptions) (TeamsResponse, error) {
	var res TeamsResponse
	resp, err := c.client.getList(ctx, pathTeams, &res, &pageOptions)
	if err != nil {
		return TeamsResponse{}, err
	}
	applyPaged(resp, &res.Paged)
	return res, apiError(resp)
}
