package lokalise

import (
	"context"
	"fmt"
)

type TeamUserGroup struct {
	GroupID     int64        `json:"group_id"`
	TeamID      int64        `json:"team_id"`
	Name        string       `json:"Name"`
	Permissions []Permission `json:"permissions,omitempty"`
	CreatedAt   string       `json:"created_at"`
	Projects    []string     `json:"projects"`
	Members     []int64      `json:"members"`
}

type Permission struct {
	IsAdmin     bool       `json:"is_admin"`
	IsReviewer  bool       `json:"is_reviewer"`
	AdminRights []string   `json:"admin_rights"`
	Languages   []Language `json:"languages"` // todo check response, doesn't fully match to Language struct
}

type TeamUserGroupsResponse struct {
	Paged
	TeamID     int64           `json:"team_id"`
	UserGroups []TeamUserGroup `json:"user_groups"`
}

type TeamUserGroupsService struct {
	client *Client
}

func pathTeamUserGroups(teamID int64) string {
	return fmt.Sprintf("%s/%d/groups", pathTeams, teamID)
}

// Method implementations

func (c *TeamUserGroupsService) List(ctx context.Context, teamID int64, pageOptions PageOptions) (TeamUserGroupsResponse, error) {
	var res TeamUserGroupsResponse
	resp, err := c.client.getList(ctx, pathTeamUserGroups(teamID), &res, &pageOptions)
	if err != nil {
		return TeamUserGroupsResponse{}, err
	}
	applyPaged(resp, &res.Paged)
	return res, apiError(resp)
}
