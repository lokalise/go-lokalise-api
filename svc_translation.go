package lokalise

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

const (
	pathTranslations = "translations"
)

type TranslationsService struct {
	BaseService
}

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
	Translations []Translation `json:"translations"`
}

type TranslationResponse struct {
	WithProjectID
	Translation Translation `json:"translation,omitempty"`
}

type TranslationsOptions struct {
	PageOptions
	DisableReferences bool
}

func (options TranslationsOptions) Apply(req *resty.Request) {
	options.PageOptions.Apply(req)
	if options.DisableReferences {
		req.SetQueryParam("disable_references", "1")
	}
}

func (c *TranslationsService) List(projectID string, pageOptions TranslationsOptions) (r TranslationsResponse, err error) {
	resp, err := c.getList(c.Ctx(), fmt.Sprintf("%s/%s/%s", pathProjects, projectID, pathTranslations), &r, pageOptions)

	if err != nil {
		return
	}
	applyPaged(resp, &r.Paged)
	return r, apiError(resp)
}

func (c *TranslationsService) Retrieve(projectID string, translationID int64) (r TranslationResponse, err error) {
	resp, err := c.get(c.Ctx(), fmt.Sprintf("%s/%s/%s/%d", pathProjects, projectID, pathTranslations, translationID), &r)

	if err != nil {
		return
	}
	return r, apiError(resp)
}

func (c *TranslationsService) Update(projectID string, translationID int64, translation string, isFuzzy, isReviewed bool) (r TranslationResponse, err error) {
	resp, err := c.put(
		c.Ctx(),
		fmt.Sprintf("%s/%s/%s/%d", pathProjects, projectID, pathTranslations, translationID),
		&r,
		map[string]interface{}{
			"translation": translation,
			"is_fuzzy":    isFuzzy,
			"is_reviewed": isReviewed,
		},
	)

	if err != nil {
		return
	}
	return r, apiError(resp)
}
