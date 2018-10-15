package lokalise

import (
	"context"
	"fmt"

	"github.com/go-resty/resty"
	"github.com/lokalise/lokalise-go-sdk/model"
)

type TeamUsersService struct {
	httpClient *resty.Client
}

const (
	pathTeams = "teams"
)

func pathTeamUsers(teamID int64) string {
	return fmt.Sprintf("%s/%d/users", pathTeams, teamID)
}

func (c *TeamUsersService) List(ctx context.Context, teamID int64, pageOptions PageOptions) (model.TeamUsersResponse, error) {
	var res model.TeamUsersResponse
	req := c.httpClient.R().
		SetResult(&res).
		SetContext(ctx)
	applyPageOptions(req, pageOptions)
	resp, err := req.Get(pathTeamUsers(teamID))
	if err != nil {
		return model.TeamUsersResponse{}, err
	}
	applyPaged(resp, &res.Paged)
	return res, apiError(resp)
}

func (c *TeamUsersService) Retrieve(ctx context.Context, teamID, userID int64) (model.TeamUserResponse, error) {
	var res model.TeamUserResponse
	req := c.httpClient.R().
		SetResult(&res).
		SetContext(ctx)
	resp, err := req.Get(fmt.Sprintf("%s/%d", pathTeamUsers(teamID), userID))
	if err != nil {
		return model.TeamUserResponse{}, err
	}
	return res, apiError(resp)
}

func (c *TeamUsersService) UpdateRole(ctx context.Context, teamID, userID int64, role model.TeamUserRole) (model.TeamUserResponse, error) {
	var res model.TeamUserResponse
	req := c.httpClient.R().
		SetResult(&res).
		SetContext(ctx).
		SetBody(map[string]interface{}{
			"role": role,
		})
	resp, err := req.Put(fmt.Sprintf("%s/%d", pathTeamUsers(teamID), userID))
	if err != nil {
		return model.TeamUserResponse{}, err
	}
	return res, apiError(resp)
}

func (c *TeamUsersService) Delete(ctx context.Context, teamID, userID int64) (model.TeamUserDeleteResponse, error) {
	var res model.TeamUserDeleteResponse
	req := c.httpClient.R().
		SetResult(&res).
		SetContext(ctx)
	resp, err := req.Delete(fmt.Sprintf("%s/%d", pathTeamUsers(teamID), userID))
	if err != nil {
		return model.TeamUserDeleteResponse{}, err
	}
	return res, apiError(resp)
}
