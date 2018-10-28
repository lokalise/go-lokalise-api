package model

type Project struct {
	ProjectID      string `json:"project_id,omitempty"`
	Name           string `json:"name,omitempty"`
	Description    string `json:"description,omitempty"`
	CreatedAt      string `json:"created_at,omitempty"`
	CreatedBy      int64  `json:"created_by,omitempty"`
	CreatedByEmail string `json:"created_by_email,omitempty"`
	TeamID         int64  `json:"team_id,omitempty"`
}

type ProjectsResponse struct {
	Paged
	Projects []Project `json:"projects,omitempty"`
}

type ProjectEmptyResponse struct {
	ProjectID   string `json:"project_id,omitempty"`
	KeysDeleted bool   `json:"keys_deleted,omitempty"`
}
