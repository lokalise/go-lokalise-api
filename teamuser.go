package lokalise

import (
	"context"
	"fmt"
)

type TeamUsersService struct {
	client *Client
}

func pathTeamUsers(teamID int64) string {
	return fmt.Sprintf("%s/%d/users", pathTeams, teamID)
}

func (c *TeamUsersService) List(ctx context.Context, teamID int64, pageOptions PageOptions) (TeamUsersResponse, error) {
	var res TeamUsersResponse
	resp, err := c.client.getList(ctx, pathTeamUsers(teamID), &res, &pageOptions)
	if err != nil {
		return TeamUsersResponse{}, err
	}
	applyPaged(resp, &res.Paged)
	return res, apiError(resp)
}

func (c *TeamUsersService) Retrieve(ctx context.Context, teamID, userID int64) (TeamUserResponse, error) {
	var res TeamUserResponse
	resp, err := c.client.get(ctx, fmt.Sprintf("%s/%d", pathTeamUsers(teamID), userID), &res)
	if err != nil {
		return TeamUserResponse{}, err
	}
	return res, apiError(resp)
}

func (c *TeamUsersService) UpdateRole(ctx context.Context, teamID, userID int64, role TeamUserRole) (TeamUserResponse, error) {
	var res TeamUserResponse
	resp, err := c.client.put(ctx, fmt.Sprintf("%s/%d", pathTeamUsers(teamID), userID), &res, map[string]interface{}{
		"role": role,
	})
	if err != nil {
		return TeamUserResponse{}, err
	}
	return res, apiError(resp)
}

func (c *TeamUsersService) Delete(ctx context.Context, teamID, userID int64) (TeamUserDeleteResponse, error) {
	var res TeamUserDeleteResponse
	resp, err := c.client.delete(ctx, fmt.Sprintf("%s/%d", pathTeamUsers(teamID), userID), &res)
	if err != nil {
		return TeamUserDeleteResponse{}, err
	}
	return res, apiError(resp)
}
