package lokalise

import (
	"context"
	"fmt"
)

type TeamUser struct {
	UserID    int64        `json:"user_id,omitempty"`
	Email     string       `json:"email,omitempty"`
	Fullname  string       `json:"fullname,omitempty"`
	CreatedAt string       `json:"created_at,omitempty"`
	Role      TeamUserRole `json:"role,omitempty"`
}

type TeamUserResponse struct {
	TeamID   int64    `json:"team_id,omitempty"`
	TeamUser TeamUser `json:"team_user,omitempty"`
}

type TeamUserRole string

//noinspection GoUnusedConst
const (
	TeamUserRoleOwner  TeamUserRole = "owner"
	TeamUserRoleAdmin  TeamUserRole = "admin"
	TeamUserRoleMember TeamUserRole = "member"
)

type TeamUserDeleteResponse struct {
	TeamID  int64 `json:"team_id,omitempty"`
	Deleted bool  `json:"team_user_deleted"`
}

type TeamUsersResponse struct {
	Paged
	TeamID    int64      `json:"team_id,omitempty"`
	TeamUsers []TeamUser `json:"team_users,omitempty"`
}

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
