package lokalise

import (
	"context"
)

type BaseService struct {
	*restClient

	PageOptions
	ctx context.Context
}

func (s *BaseService) Ctx() context.Context {
	if s.ctx != nil {
		return s.ctx
	}
	return context.Background()
}

func (s *BaseService) SetContext(c context.Context) {
	s.ctx = c
}

func (s *BaseService) PageOpts() PageOptions {
	return s.PageOptions
}

func (s *BaseService) SetPageOptions(opts PageOptions) {
	if opts.Limit != 0 {
		s.Limit = opts.Limit
	}
	if opts.Page != 0 {
		s.Page = opts.Page
	}
}

// ‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// Additional subtypes
// _____________________________________________________________________________________________________________________

type WithCreationTime struct {
	CreatedAt   string `json:"created_at"`
	CreatedAtTs int64  `json:"created_at_timestamp"`
}

type WithCreationUser struct {
	CreatedBy      int64  `json:"created_by"`
	CreatedByEmail string `json:"created_by_email"`
}

// could be optional for some creation structs, i.e. NewProject
type WithTeamID struct {
	TeamID int64 `json:"team_id,omitempty"`
}

type WithProjectID struct {
	ProjectID string `json:"project_id,omitempty"`
}

type WithUserID struct {
	UserID int64 `json:"user_id,omitempty"`
}
