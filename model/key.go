package model

type PlatformStrings struct {
	Ios     string `json:"ios,omitempty"`
	Android string `json:"android,omitempty"`
	Web     string `json:"web,omitempty"`
	Other   string `json:"other,omitempty"`
}

type Key struct {
	KeyID            int64           `json:"key_id,omitempty"`
	CreatedAt        string          `json:"created_at,omitempty"`
	KeyName          interface{}     `json:"key_name,omitempty"` // KeyName could be string or PlatformStrings
	Filenames        PlatformStrings `json:"filenames,omitempty"`
	Description      string          `json:"description,omitempty"`
	Platforms        []string        `json:"platforms,omitempty"`
	Tags             []string        `json:"tags,omitempty"`
	Comments         []Comment       `json:"comments,omitempty"`
	Screenshots      []Screenshot    `json:"screenshots,omitempty"`
	Translations     []Translation   `json:"translations,omitempty"`
	IsPlural         bool            `json:"is_plural,omitempty"`
	PluralName       string          `json:"plural_name,omitempty"`
	IsHidden         bool            `json:"is_hidden,omitempty"`
	IsArchived       bool            `json:"is_archived,omitempty"`
	Context          string          `json:"context,omitempty"`
	CharLimit        int             `json:"char_limit,omitempty"`
	CustomAttributes []interface{}   `json:"custom_attributes,omitempty"`
}

// ErrorKey is key info from error for key create/update API
type ErrorKey struct {
	KeyName string `json:"key_name,omitempty"`
}

// ErrorKeys is error for key create/update API
type ErrorKeys struct {
	Error
	Key ErrorKey `json:"key,omitempty"`
}

type KeysResponse struct {
	Paged
	ProjectID string      `json:"project_id,omitempty"`
	Keys      []Key       `json:"keys,omitempty"`
	Errors    []ErrorKeys `json:"error,omitempty"`
}
