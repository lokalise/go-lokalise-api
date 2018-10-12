package model

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

const (
	TeamUserRoleOwner  TeamUserRole = "owner"
	TeamUserRoleAdmin  TeamUserRole = "admin"
	TeamUserRoleMember TeamUserRole = "member"
)

type TeamUserDeleteResponse struct {
	TeamID  int64 `json:"team_id,omitempty"`
	Deleted bool  `json:"team_user_deleted"`
}
