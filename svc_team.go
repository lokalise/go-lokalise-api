package lokalise

import "fmt"

const (
	pathTeams = "teams"
)

// The Team service
type TeamService struct {
	BaseService
}

// ‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// Service entity objects
// _____________________________________________________________________________________________________________________

type Team struct {
	WithCreationTime
	WithTeamID

	Name         string `json:"name"`
	Plan         string `json:"plan"` // e.g. "Essential", "Trial"
	QuotaUsage   Quota  `json:"quota_usage"`
	QuotaAllowed Quota  `json:"quota_allowed"`
}

type Quota struct {
	Users    int64 `json:"users"`
	Keys     int64 `json:"keys"`
	Projects int64 `json:"projects"`
	MAU      int64 `json:"mau"`
}

type PermissionRole struct {
	ID                             int      `json:"id"`
	Role                           string   `json:"role"`
	Permissions                    []string `json:"permissions"`
	Description                    string   `json:"description"`
	Tag                            string   `json:"tag"`
	TagColor                       string   `json:"tagColor"`
	DoesEnableAllReadOnlyLanguages bool     `json:"doesEnableAllReadOnlyLanguages"`
}

type PermissionRoleResponse struct {
	Roles []PermissionRole `json:"roles"`
}

// ‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// Service request/response objects
// _____________________________________________________________________________________________________________________

type TeamsResponse struct {
	Paged
	Teams []Team `json:"teams"`
}

// ‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// Service methods
// _____________________________________________________________________________________________________________________

// Lists all teams available to the user
func (c *TeamService) List() (r TeamsResponse, err error) {
	resp, err := c.getWithOptions(c.Ctx(), pathTeams, &r, c.PageOpts())

	if err != nil {
		return
	}
	applyPaged(resp, &r.Paged)
	return r, apiError(resp)
}

// List all possible permission roles
func (c *TeamService) ListPermissionRoles(teamID int64) (r PermissionRoleResponse, err error) {
	resp, err := c.getWithOptions(c.Ctx(), pathPermissionRoles(teamID), &r, c.PageOpts())

	if err != nil {
		return r, err
	}
	return r, apiError(resp)
}

func pathPermissionRoles(teamID int64) string {
	return fmt.Sprintf("teams/%d/roles", teamID)
}
