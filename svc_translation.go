package lokalise

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/google/go-querystring/query"
)

const (
	pathTranslations = "translations"
)

// TranslationService supports List, Retrieve and Update commands
type TranslationService struct {
	BaseService

	opts         TranslationListOptions
	retrieveOpts TranslationRetrieveOptions
}

// ‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// Service entity objects
// _____________________________________________________________________________________________________________________

type Translation struct {
	TranslationID   int64  `json:"translation_id"`
	Translation     string `json:"translation"` // could be string or json in case it includes plural forms and is_plural is true.
	KeyID           int64  `json:"key_id"`
	LanguageISO     string `json:"language_iso"`
	ModifiedAt      string `json:"modified_at"`
	ModifiedAtTs    int64  `json:"modified_at_timestamp"`
	ModifiedBy      int64  `json:"modified_by"`
	ModifiedByEmail string `json:"modified_by_email"`
	IsUnverified    bool   `json:"is_unverified"`
	IsReviewed      bool   `json:"is_reviewed"`
	ReviewedBy      int64  `json:"reviewed_by"`
	Words           int64  `json:"words"`
	TaskID          int64  `json:"task_id"`
	SegmentNumber   int64  `json:"segment_number"`

	CustomTranslationStatuses []TranslationStatus `json:"custom_translation_statuses"`
}

// ‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// Service request/response objects
// _____________________________________________________________________________________________________________________

// NewTranslation Used for NewKey
type NewTranslation struct {
	LanguageISO                    string  `json:"language_iso"`
	Translation                    string  `json:"translation"`
	IsFuzzy                        *bool   `json:"is_fuzzy,omitempty"`
	IsReviewed                     bool    `json:"is_reviewed,omitempty"`
	CustomTranslationStatusIds     []int64 `json:"custom_translation_status_ids,omitempty"`
	MergeCustomTranslationStatuses bool    `json:"merge_custom_translation_statuses,omitempty"`
}

func (t NewTranslation) MarshalJSON() ([]byte, error) {
	type Alias NewTranslation

	var translation interface{} = t.Translation

	if json.Valid([]byte(t.Translation)) {
		_ = json.Unmarshal([]byte(t.Translation), &translation)
	}

	return json.Marshal(&struct {
		LanguageISO string      `json:"language_iso"`
		Translation interface{} `json:"translation"`
		Alias
	}{
		LanguageISO: t.LanguageISO,
		Translation: translation,
		Alias:       (Alias)(t),
	})
}

func (t *NewTranslation) UnmarshalJSON(data []byte) error {
	type Alias NewTranslation
	aux := &struct {
		LanguageISO string      `json:"language_iso"`
		Translation interface{} `json:"translation"`
		*Alias
	}{
		Alias: (*Alias)(t),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	switch aux.Translation.(type) {
	case map[string]interface{}:
		marshal, err := json.Marshal(aux.Translation)
		if err != nil {
			return err
		}
		t.Translation = string(marshal)
	case string:
		t.Translation = fmt.Sprintf("%+v", aux.Translation)
	default:
		return fmt.Errorf("NewTranslation tranlation type is of unknown type")
	}

	t.LanguageISO = aux.LanguageISO
	t.IsFuzzy = aux.IsFuzzy
	t.IsReviewed = aux.IsReviewed
	t.CustomTranslationStatusIds = aux.CustomTranslationStatusIds
	t.MergeCustomTranslationStatuses = aux.MergeCustomTranslationStatuses

	return nil
}

type UpdateTranslation struct {
	Translation                string   `json:"translation"`
	IsUnverified               *bool    `json:"is_unverified,omitempty"`
	IsReviewed                 bool     `json:"is_reviewed,omitempty"`
	CustomTranslationStatusIDs []string `json:"custom_translation_status_ids,omitempty"`
}

func (t UpdateTranslation) MarshalJSON() ([]byte, error) {
	type Alias UpdateTranslation

	var translation interface{} = t.Translation

	if json.Valid([]byte(t.Translation)) {
		_ = json.Unmarshal([]byte(t.Translation), &translation)
	}

	return json.Marshal(&struct {
		Translation interface{} `json:"translation"`
		Alias
	}{
		Translation: translation,
		Alias:       (Alias)(t),
	})
}

func (t *UpdateTranslation) UnmarshalJSON(data []byte) error {
	type Alias UpdateTranslation
	aux := &struct {
		LanguageISO string      `json:"language_iso"`
		Translation interface{} `json:"translation"`
		*Alias
	}{
		Alias: (*Alias)(t),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	switch aux.Translation.(type) {
	case map[string]interface{}:
		marshal, err := json.Marshal(aux.Translation)
		if err != nil {
			return err
		}
		t.Translation = string(marshal)
	case string:
		t.Translation = fmt.Sprintf("%+v", aux.Translation)
	default:
		return fmt.Errorf("NewTranslation tranlation type is of unknown type")
	}

	t.IsUnverified = aux.IsUnverified
	t.IsReviewed = aux.IsReviewed
	t.CustomTranslationStatusIDs = aux.CustomTranslationStatusIDs

	return nil
}

type TranslationsResponse struct {
	Paged
	Translations []Translation `json:"translations"`
}

type TranslationResponse struct {
	WithProjectID
	Translation Translation `json:"translation"`
}

// ‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// Service methods
// _____________________________________________________________________________________________________________________

func (c *TranslationService) List(projectID string) (r TranslationsResponse, err error) {
	resp, err := c.getWithOptions(c.Ctx(), fmt.Sprintf("%s/%s/%s", pathProjects, projectID, pathTranslations), &r, c.ListOpts())

	if err != nil {
		return
	}
	applyPaged(resp, &r.Paged)
	return r, apiError(resp)
}

func (c *TranslationService) Retrieve(projectID string, translationID int64) (r TranslationResponse, err error) {
	resp, err := c.getWithOptions(c.Ctx(), fmt.Sprintf("%s/%s/%s/%d", pathProjects, projectID, pathTranslations, translationID), &r, c.RetrieveOpts())

	if err != nil {
		return
	}
	return r, apiError(resp)
}

func (c *TranslationService) Update(projectID string, translationID int64, opts UpdateTranslation) (r TranslationResponse, err error) {
	resp, err := c.put(c.Ctx(), fmt.Sprintf("%s/%s/%s/%d", pathProjects, projectID, pathTranslations, translationID), &r, opts)

	if err != nil {
		return
	}
	return r, apiError(resp)
}

// ‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// Additional methods
// _____________________________________________________________________________________________________________________

type TranslationListOptions struct {
	// page options
	Page  uint `url:"page,omitempty"`
	Limit uint `url:"limit,omitempty"`

	// Possible values are 1 and 0.
	DisableReferences uint8 `url:"disable_references,omitempty"`

	FilterLangID       int64  `url:"filter_lang_id,omitempty"`
	FilterIsReviewed   uint8  `url:"filter_is_reviewed,omitempty"`
	FilterUnverified   uint8  `url:"filter_unverified,omitempty"`
	FilterUntranslated uint8  `url:"filter_untranslated,omitempty"`
	FilterQAIssues     string `url:"filter_qa_issues,omitempty"`
	FilterActiveTaskID int64  `url:"filter_active_task_id,omitempty"`
}

func (options TranslationListOptions) Apply(req *resty.Request) {
	v, _ := query.Values(options)
	req.SetQueryString(v.Encode())
}

func (c *TranslationService) ListOpts() TranslationListOptions        { return c.opts }
func (c *TranslationService) SetListOptions(o TranslationListOptions) { c.opts = o }
func (c *TranslationService) WithListOptions(o TranslationListOptions) *TranslationService {
	c.opts = o
	return c
}

type TranslationRetrieveOptions struct {
	DisableReferences uint8 `url:"disable_references,omitempty"`
}

func (options TranslationRetrieveOptions) Apply(req *resty.Request) {
	v, _ := query.Values(options)
	req.SetQueryString(v.Encode())
}

func (c *TranslationService) RetrieveOpts() TranslationRetrieveOptions        { return c.retrieveOpts }
func (c *TranslationService) SetRetrieveOptions(o TranslationRetrieveOptions) { c.retrieveOpts = o }
func (c *TranslationService) WithRetrieveOptions(o TranslationRetrieveOptions) *TranslationService {
	c.retrieveOpts = o
	return c
}
