package lokalise

const (
	pathTeams = "teams"
)

type TeamsResponse struct {
	Paged
	Teams []Team `json:"teams,omitempty"`
}

type Team struct {
	WithCreationTime
	WithTeamID
	Name         string `json:"name,omitempty"`
	Plan         string `json:"plan,omitempty"`
	QuotaUsage   Quota  `json:"quota_usage,omitempty"`
	QuotaAllowed Quota  `json:"quota_allowed,omitempty"`
}

type Quota struct {
	Users    int64 `json:"users,omitempty"`
	Keys     int64 `json:"keys,omitempty"`
	Projects int64 `json:"projects,omitempty"`
	MAU      int64 `json:"mau,omitempty"`
}

type TeamsService struct {
	BaseService
}

func (c *TeamsService) List() (r TeamsResponse, err error) {
	resp, err := c.getList(c.Ctx(), pathTeams, &r, c.PageOpts())

	if err != nil {
		return
	}
	applyPaged(resp, &r.Paged)
	return r, apiError(resp)
}
