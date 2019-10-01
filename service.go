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

func (s *BaseService) WithContext(c context.Context) *BaseService {
	s.ctx = c
	return s
}

func (s *BaseService) PageOpts() PageOptions {
	return s.PageOptions
}

func (s *BaseService) WithPageOptions(Limit, Offset int64) *BaseService {
	s.Limit = Limit
	s.Page = Offset
	return s
}

// Additional subtypes
// -------------------

type WithCreationTime struct {
	CreatedAt   string `json:"created_at,omitempty"`
	CreatedAtTs int64  `json:"created_at_timestamp,omitempty"`
}

type WithTeamID struct {
	TeamID int64 `json:"team_id"`
}

type WithProjectID struct {
	ProjectID string `json:"project_id,omitempty"`
}

type WithUserID struct {
	UserID int64 `json:"user_id,omitempty"`
}
