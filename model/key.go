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
	BaseWords        int64           `json:"base_words,omitempty"`
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

type CreateKeyRequest struct {
	KeyName          interface{}            `json:"key_name,omitempty"` // KeyName could be string or PlatformStrings
	Description      string                 `json:"description,omitempty"`
	Platforms        []string               `json:"platforms,omitempty"`
	Filenames        PlatformStrings        `json:"filenames,omitempty"`
	Tags             []string               `json:"tags,omitempty"`
	Comments         []CreateKeyComment     `json:"comments,omitempty"`
	Screenshots      []CreateKeyScreenshot  `json:"screenshots,omitempty"`
	Translations     []CreateKeyTranslation `json:"translations,omitempty"`
	IsPlural         bool                   `json:"is_plural,omitempty"`
	PluralName       string                 `json:"plural_name,omitempty"`
	IsHidden         bool                   `json:"is_hidden,omitempty"`
	IsArchived       bool                   `json:"is_archived,omitempty"`
	Context          string                 `json:"context,omitempty"`
	CustomAttributes []interface{}          `json:"custom_attributes,omitempty"`
}

type CreateKeyComment struct {
	Comment string `json:"comment,omitempty"`
}

type CreateKeyScreenshot struct {
	Title          string   `json:"title,omitempty"`
	Description    string   `json:"description,omitempty"`
	ScreenshotTags []string `json:"screenshot_tags,omitempty"`
	Data           []byte   `json:"data,omitempty"`
}

type CreateKeyTranslation struct {
	LanguageISO string `json:"language_iso,omitempty"`
	Translation string `json:"translation,omitempty"`
	IsFuzzy     bool   `json:"is_fuzzy,omitempty"`
	IsReviewed  bool   `json:"is_reviewed,omitempty"`
}
