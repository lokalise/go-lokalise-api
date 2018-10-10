package model

type TeamUser struct {
	UserID    int64  `json:"user_id,omitempty"`
	Email     string `json:"email,omitempty"`
	Fullname  string `json:"fullname,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	Role      string `json:"role,omitempty"`
}

type TeamUserResponse struct {
	TeamID   int64    `json:"team_id,omitempty"`
	TeamUser TeamUser `json:"team_user,omitempty"`
}
