package lokalise

import (
	"fmt"
)

type TeamUserService struct {
	BaseService
}

// ‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// Service entity objects
// _____________________________________________________________________________________________________________________

//noinspection GoUnusedConst
const (
	TeamUserRoleOwner  TeamUserRole = "owner"
	TeamUserRoleAdmin  TeamUserRole = "admin"
	TeamUserRoleMember TeamUserRole = "member"
)

type TeamUserRole string

type TeamUser struct {
	WithCreationTime
	WithUserID

	Email    string       `json:"email"`
	Fullname string       `json:"fullname"`
	Role     TeamUserRole `json:"role"`
}

// ‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// Service request/response objects
// _____________________________________________________________________________________________________________________

type TeamUserResponse struct {
	WithTeamID
	TeamUser TeamUser `json:"team_user"`
}

type DeleteTeamUserResponse struct {
	WithTeamID
	Deleted bool `json:"team_user_deleted"`
}

type TeamUsersResponse struct {
	Paged
	WithTeamID
	TeamUsers []TeamUser `json:"team_users"`
}

// ‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// Service methods
// _____________________________________________________________________________________________________________________

func (c *TeamUserService) List(teamID int64) (r TeamUsersResponse, err error) {
	resp, err := c.getList(c.Ctx(), pathTeamUsers(teamID), &r, c.PageOpts())

	if err != nil {
		return
	}
	applyPaged(resp, &r.Paged)
	return r, apiError(resp)
}

func (c *TeamUserService) Retrieve(teamID, userID int64) (res TeamUserResponse, err error) {
	resp, err := c.get(c.Ctx(), fmt.Sprintf("%s/%d", pathTeamUsers(teamID), userID), &res)

	if err != nil {
		return
	}
	return res, apiError(resp)
}

func (c *TeamUserService) UpdateRole(teamID, userID int64, role TeamUserRole) (r TeamUserResponse, err error) {
	resp, err := c.put(c.Ctx(), fmt.Sprintf("%s/%d", pathTeamUsers(teamID), userID), &r, map[string]interface{}{
		"role": role,
	})
	if err != nil {
		return
	}
	return r, apiError(resp)
}

func (c *TeamUserService) Delete(teamID, userID int64) (r DeleteTeamUserResponse, err error) {
	resp, err := c.delete(c.Ctx(), fmt.Sprintf("%s/%d", pathTeamUsers(teamID), userID), &r)

	if err != nil {
		return
	}
	return r, apiError(resp)
}

func pathTeamUsers(teamID int64) string {
	return fmt.Sprintf("%s/%d/users", pathTeams, teamID)
}
