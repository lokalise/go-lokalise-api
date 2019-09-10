package lokalise

import (
	"context"
	"fmt"
	"path"
	"strconv"
)

type TeamUserGroup struct {
	GroupID     int64        `json:"group_id"`
	TeamID      int64        `json:"team_id"`
	Name        string       `json:"name"`
	Permissions []Permission `json:"permissions,omitempty"`
	CreatedAt   string       `json:"created_at"`
	Projects    []string     `json:"projects"`
	Members     []int64      `json:"members"`
}

type Permission struct { // split near languages for more convenient implementation in group creation (customGroup)
	IsAdmin     bool       `json:"is_admin"`
	IsReviewer  bool       `json:"is_reviewer"`
	AdminRights []string   `json:"admin_rights,omitempty"` // todo make admin rights as constants available in lib
	Languages   []Language `json:"languages,omitempty"`    // todo check response, doesn't fully match to Language struct
}

type CustomGroup struct {
	Name       string `json:"name"`
	Permission        // for creation lang are also not a slice of the language structs
}

type TeamUserGroupsResponse struct {
	Paged
	TeamID     int64           `json:"team_id"`
	UserGroups []TeamUserGroup `json:"user_groups"`
}

type CreateGroupResponse struct {
	TeamID int64         `json:"team_id"`
	Group  TeamUserGroup `json:"group"`
}

type DeleteGroupResponse struct {
	TeamID    int64 `json:"team_id"`
	IsDeleted bool  `json:"group_deleted"`
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

func (c *TeamUserGroupsService) Create(ctx context.Context, teamID int64, group CustomGroup) (result CreateGroupResponse, err error) {
	resp, err := c.client.post(ctx, pathTeamUserGroups(teamID), &result, group)

	if err != nil {
		return result, err
	}
	return result, apiError(resp)
}

func (c *TeamUserGroupsService) Retrieve(ctx context.Context, teamID, groupID int64) (result TeamUserGroup, err error) {
	url := path.Join(pathTeamUserGroups(teamID), strconv.FormatInt(groupID, 10))
	resp, err := c.client.get(ctx, url, &result)

	if err != nil {
		return result, err
	}
	return result, apiError(resp)
}

func (c *TeamUserGroupsService) Update(ctx context.Context, teamID, groupID int64, group CustomGroup) (result CreateGroupResponse, err error) {
	url := path.Join(pathTeamUserGroups(teamID), strconv.FormatInt(groupID, 10))
	resp, err := c.client.put(ctx, url, &result, group)

	if err != nil {
		return result, err
	}
	return result, apiError(resp)
}

// todo make internal (unexported) function for all of the next fours
// wrap em all!
func (c *TeamUserGroupsService) AddProjects(ctx context.Context, teamID, groupID int64, projects []string) (result CreateGroupResponse, err error) {
	url := path.Join(pathTeamUserGroups(teamID), strconv.FormatInt(groupID, 10), "projects", "add")
	resp, err := c.client.put(ctx, url, &result, map[string]interface{}{
		"projects": projects,
	})

	if err != nil {
		return result, err
	}
	return result, apiError(resp)
}

func (c *TeamUserGroupsService) RemoveProjects(ctx context.Context, teamID, groupID int64, projects []string) (result CreateGroupResponse, err error) {
	url := path.Join(pathTeamUserGroups(teamID), strconv.FormatInt(groupID, 10), "projects", "remove")
	resp, err := c.client.put(ctx, url, &result, map[string]interface{}{
		"projects": projects,
	})

	if err != nil {
		return result, err
	}
	return result, apiError(resp)
}

func (c *TeamUserGroupsService) AddMembers(ctx context.Context, teamID, groupID int64, users []int64) (result CreateGroupResponse, err error) {
	url := path.Join(pathTeamUserGroups(teamID), strconv.FormatInt(groupID, 10), "members", "add")
	resp, err := c.client.put(ctx, url, &result, map[string]interface{}{
		"users": users,
	})

	if err != nil {
		return result, err
	}
	return result, apiError(resp)
}

func (c *TeamUserGroupsService) RemoveMembers(ctx context.Context, teamID, groupID int64, users []int64) (result CreateGroupResponse, err error) {
	url := path.Join(pathTeamUserGroups(teamID), strconv.FormatInt(groupID, 10), "members", "remove")
	resp, err := c.client.put(ctx, url, &result, map[string]interface{}{
		"users": users,
	})

	if err != nil {
		return result, err
	}
	return result, apiError(resp)
}

func (c *TeamUserGroupsService) Delete(ctx context.Context, teamID, groupID int64) (result DeleteGroupResponse, err error) {
	url := path.Join(pathTeamUserGroups(teamID), strconv.FormatInt(groupID, 10))
	resp, err := c.client.delete(ctx, url, &result)

	if err != nil {
		return result, err
	}
	return result, apiError(resp)
}
