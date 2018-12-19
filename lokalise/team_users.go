package lokalise

import (
	"context"
	"fmt"

	"github.com/17media/go-lokalise-api/model"
)

type TeamUsersService struct {
	client *Client
}

func pathTeamUsers(teamID int64) string {
	return fmt.Sprintf("%s/%d/users", pathTeams, teamID)
}

func (c *TeamUsersService) List(ctx context.Context, teamID int64, pageOptions PageOptions) (model.TeamUsersResponse, error) {
	var res model.TeamUsersResponse
	resp, err := c.client.getList(ctx, pathTeamUsers(teamID), &res, &pageOptions)
	if err != nil {
		return model.TeamUsersResponse{}, err
	}
	applyPaged(resp, &res.Paged)
	return res, apiError(resp)
}

func (c *TeamUsersService) Retrieve(ctx context.Context, teamID, userID int64) (model.TeamUserResponse, error) {
	var res model.TeamUserResponse
	resp, err := c.client.get(ctx, fmt.Sprintf("%s/%d", pathTeamUsers(teamID), userID), &res)
	if err != nil {
		return model.TeamUserResponse{}, err
	}
	return res, apiError(resp)
}

func (c *TeamUsersService) UpdateRole(ctx context.Context, teamID, userID int64, role model.TeamUserRole) (model.TeamUserResponse, error) {
	var res model.TeamUserResponse
	resp, err := c.client.put(ctx, fmt.Sprintf("%s/%d", pathTeamUsers(teamID), userID), &res, map[string]interface{}{
		"role": role,
	})
	if err != nil {
		return model.TeamUserResponse{}, err
	}
	return res, apiError(resp)
}

func (c *TeamUsersService) Delete(ctx context.Context, teamID, userID int64) (model.TeamUserDeleteResponse, error) {
	var res model.TeamUserDeleteResponse
	resp, err := c.client.delete(ctx, fmt.Sprintf("%s/%d", pathTeamUsers(teamID), userID), &res)
	if err != nil {
		return model.TeamUserDeleteResponse{}, err
	}
	return res, apiError(resp)
}
