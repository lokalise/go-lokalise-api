package lokalise

import "fmt"

const (
	pathTemplates = "teams/%d/roles"
)

// The PermissionTemplate service
type PermissionTemplateService struct {
	BaseService
}

type PermissionTemplate struct {
	ID                             int      `json:"id"`
	Role                           string   `json:"role"`
	Permissions                    []string `json:"permissions"`
	Description                    string   `json:"description"`
	Tag                            string   `json:"tag"`
	TagColor                       string   `json:"tagColor"`
	TagInfo                        string   `json:"tagInfo"`
	DoesEnableAllReadOnlyLanguages bool     `json:"doesEnableAllReadOnlyLanguages"`
}

type PermissionRoleResponse struct {
	Roles []PermissionTemplate `json:"roles"`
}

// List all possible permission roles
func (c *PermissionTemplateService) ListPermissionRoles(teamID int64) (r PermissionRoleResponse, err error) {
	resp, err := c.getWithOptions(c.Ctx(), pathPermissionRoles(teamID), &r, c.PageOpts())

	if err != nil {
		return r, err
	}
	return r, apiError(resp)
}

func pathPermissionRoles(teamID int64) string {
	return fmt.Sprintf(pathTemplates, teamID)
}
