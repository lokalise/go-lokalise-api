package model

type Translation struct {
	TranslationID   int64  `json:"translation_id,omitempty"`
	KeyID           int64  `json:"key_id,omitempty"`
	LanguageISO     string `json:"language_iso,omitempty"`
	ModifiedAt      string `json:"modified_at,omitempty"`
	ModifiedBy      int64  `json:"modified_by,omitempty"`
	ModifiedByEmail string `json:"modified_by_email,omitempty"`
	Translation     string `json:"translation,omitempty"`
	IsFuzzy         bool   `json:"is_fuzzy,omitempty"`
	IsReviewed      bool   `json:"is_reviewed,omitempty"`
	Words           int64  `json:"words,omitempty"`
}

type TranslationsResponse struct {
	Paged
	Translations []Translation
}

type TranslationResponse struct {
	ProjectID   string      `json:"project_id,omitempty"`
	Translation Translation `json:"translation,omitempty"`
}
