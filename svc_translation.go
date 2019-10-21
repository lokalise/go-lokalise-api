package lokalise

import (
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/google/go-querystring/query"
)

const (
	pathTranslations = "translations"
)

type TranslationService struct {
	BaseService

	opts TranslationListOptions
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
	IsFuzzy         bool   `json:"is_fuzzy"`
	IsReviewed      bool   `json:"is_reviewed"`
	ReviewedBy      int64  `json:"reviewed_by"`
	Words           int64  `json:"words"`

	CustomTranslationStatuses []TranslationStatus `json:"custom_translation_statuses"`
}

// ‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// Service request/response objects
// _____________________________________________________________________________________________________________________

// Used for NewKey
type NewTranslation struct {
	LanguageISO string `json:"language_iso"`
	Translation string `json:"translation"`
	IsFuzzy     *bool  `json:"is_fuzzy,omitempty"`
	IsReviewed  bool   `json:"is_reviewed,omitempty"`
}

type UpdateTranslation struct {
	Translation                string   `json:"translation"`
	IsFuzzy                    *bool    `json:"is_fuzzy,omitempty"`
	IsReviewed                 bool     `json:"is_reviewed,omitempty"`
	CustomTranslationStatusIDs []string `json:"custom_translation_status_ids,omitempty"`
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
	resp, err := c.getList(c.Ctx(), fmt.Sprintf("%s/%s/%s", pathProjects, projectID, pathTranslations), &r, c.ListOpts())

	if err != nil {
		return
	}
	applyPaged(resp, &r.Paged)
	return r, apiError(resp)
}

func (c *TranslationService) Retrieve(projectID string, translationID int64) (r TranslationResponse, err error) {
	resp, err := c.get(c.Ctx(), fmt.Sprintf("%s/%s/%s/%d", pathProjects, projectID, pathTranslations, translationID), &r)

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

	FilterLangID     string `url:"filter_lang_id,omitempty"`
	FilterIsReviewed uint8  `url:"filter_is_reviewed,omitempty"`
	FilterFuzzy      uint8  `url:"filter_fuzzy,omitempty"`
	FilterQAIssues   string `url:"filter_qa_issues,omitempty"`
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
