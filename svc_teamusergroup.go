package lokalise

import (
	"fmt"
	"path"
	"strconv"
)

type TeamUserGroup struct {
	WithCreationTime
	WithTeamID
	GroupID     int64        `json:"group_id"`
	Name        string       `json:"name"`
	Permissions []Permission `json:"permissions,omitempty"`
	Projects    []string     `json:"projects"`
	Members     []int64      `json:"members"`
}

type Permission struct {
	IsAdmin     bool       `json:"is_admin,omitempty"`
	IsReviewer  bool       `json:"is_reviewer,omitempty"`
	AdminRights []string   `json:"admin_rights,omitempty"` // todo make admin rights as constants available in lib
	Languages   []Language `json:"languages,omitempty"`    // todo check response for the different entities
}

type CustomGroup struct {
	Name string `json:"name"`
	Permission
}

type TeamUserGroupsResponse struct {
	Paged
	WithTeamID
	UserGroups []TeamUserGroup `json:"user_groups"`
}

type CreateGroupResponse struct {
	WithTeamID
	Group TeamUserGroup `json:"group"`
}

type DeleteGroupResponse struct {
	WithTeamID
	IsDeleted bool `json:"group_deleted"`
}

type TeamUserGroupsService struct {
	BaseService
}

func pathTeamUserGroups(teamID int64) string {
	return fmt.Sprintf("%s/%d/groups", pathTeams, teamID)
}

// Method implementations

func (c *TeamUserGroupsService) List(teamID int64) (r TeamUserGroupsResponse, err error) {
	resp, err := c.getList(c.Ctx(), pathTeamUserGroups(teamID), &r, c.PageOpts())

	if err != nil {
		return
	}
	applyPaged(resp, &r.Paged)
	return r, apiError(resp)
}

func (c *TeamUserGroupsService) Create(teamID int64, group CustomGroup) (r CreateGroupResponse, err error) {
	resp, err := c.post(c.Ctx(), pathTeamUserGroups(teamID), &r, group)

	if err != nil {
		return r, err
	}
	return r, apiError(resp)
}

func (c *TeamUserGroupsService) Retrieve(teamID, groupID int64) (r TeamUserGroup, err error) {
	url := path.Join(pathTeamUserGroups(teamID), strconv.FormatInt(groupID, 10))
	resp, err := c.get(c.Ctx(), url, &r)

	if err != nil {
		return r, err
	}
	return r, apiError(resp)
}

func (c *TeamUserGroupsService) Update(teamID, groupID int64, group CustomGroup) (r CreateGroupResponse, err error) {
	url := path.Join(pathTeamUserGroups(teamID), strconv.FormatInt(groupID, 10))
	resp, err := c.put(c.Ctx(), url, &r, group)

	if err != nil {
		return r, err
	}
	return r, apiError(resp)
}

func (c *TeamUserGroupsService) AddProjects(teamID, groupID int64, projects []string) (r CreateGroupResponse, err error) {
	url := path.Join(pathTeamUserGroups(teamID), strconv.FormatInt(groupID, 10), "projects", "add")
	resp, err := c.put(c.Ctx(), url, &r, map[string]interface{}{
		"projects": projects,
	})

	if err != nil {
		return r, err
	}
	return r, apiError(resp)
}

func (c *TeamUserGroupsService) RemoveProjects(teamID, groupID int64, projects []string) (r CreateGroupResponse, err error) {
	url := path.Join(pathTeamUserGroups(teamID), strconv.FormatInt(groupID, 10), "projects", "remove")
	resp, err := c.put(c.Ctx(), url, &r, map[string]interface{}{
		"projects": projects,
	})

	if err != nil {
		return r, err
	}
	return r, apiError(resp)
}

func (c *TeamUserGroupsService) AddMembers(teamID, groupID int64, users []int64) (r CreateGroupResponse, err error) {
	url := path.Join(pathTeamUserGroups(teamID), strconv.FormatInt(groupID, 10), "members", "add")
	resp, err := c.put(c.Ctx(), url, &r, map[string]interface{}{
		"users": users,
	})

	if err != nil {
		return r, err
	}
	return r, apiError(resp)
}

func (c *TeamUserGroupsService) RemoveMembers(teamID, groupID int64, users []int64) (r CreateGroupResponse, err error) {
	url := path.Join(pathTeamUserGroups(teamID), strconv.FormatInt(groupID, 10), "members", "remove")
	resp, err := c.put(c.Ctx(), url, &r, map[string]interface{}{
		"users": users,
	})

	if err != nil {
		return r, err
	}
	return r, apiError(resp)
}

func (c *TeamUserGroupsService) Delete(teamID, groupID int64) (r DeleteGroupResponse, err error) {
	url := path.Join(pathTeamUserGroups(teamID), strconv.FormatInt(groupID, 10))
	resp, err := c.delete(c.Ctx(), url, &r)

	if err != nil {
		return r, err
	}
	return r, apiError(resp)
}
