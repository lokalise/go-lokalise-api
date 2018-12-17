package model

type PlatformStrings struct {
	Ios     string `json:"ios,omitempty"`
	Android string `json:"android,omitempty"`
	Web     string `json:"web,omitempty"`
	Other   string `json:"other,omitempty"`
}

type Key struct {
	KeyID        int64           `json:"key_id,omitempty"`
	CreatedAt    string          `json:"created_at,omitempty"`
	KeyName      PlatformStrings `json:"key_name,omitempty"`
	Filenames    PlatformStrings `json:"filenames,omitempty"`
	Description  string          `json:"description,omitempty"`
	Platforms    []string        `json:"platforms,omitempty"`
	Tags         []string        `json:"tags,omitempty"`
	Comments     []Comment       `json:"comments,omitempty"`
	Screenshots  []Screenshot    `json:"screenshots,omitempty"`
	Translations []Translation   `json:"translations,omitempty"`
}

type KeysResponse struct {
	Paged
	ProjectID string `json:"project_id,omitempty"`
	Keys      []Key  `json:"keys,omitempty"`
}
