package lokalise

import (
	"context"
	"fmt"

	"github.com/go-resty/resty/v2"
)

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

type TranslationsService struct {
	client *Client
}

const (
	pathTranslations = "translations"
)

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

func (c *TranslationsService) List(ctx context.Context, projectID string, pageOptions TranslationsOptions) (TranslationsResponse, error) {
	var res TranslationsResponse
	resp, err := c.client.getList(ctx, fmt.Sprintf("%s/%s/%s", pathProjects, projectID, pathTranslations), &res, pageOptions)
	if err != nil {
		return TranslationsResponse{}, err
	}
	applyPaged(resp, &res.Paged)
	return res, apiError(resp)
}

func (c *TranslationsService) Retrieve(ctx context.Context, projectID string, translationID int64) (TranslationResponse, error) {
	var res TranslationResponse
	resp, err := c.client.get(ctx, fmt.Sprintf("%s/%s/%s/%d", pathProjects, projectID, pathTranslations, translationID), &res)
	if err != nil {
		return TranslationResponse{}, err
	}
	return res, apiError(resp)
}

func (c *TranslationsService) Update(ctx context.Context, projectID string, translationID int64, translation string, isFuzzy, isReviewed bool) (TranslationResponse, error) {
	var res TranslationResponse
	resp, err := c.client.put(ctx, fmt.Sprintf("%s/%s/%s/%d", pathProjects, projectID, pathTranslations, translationID), &res, map[string]interface{}{
		"translation": translation,
		"is_fuzzy":    isFuzzy,
		"is_reviewed": isReviewed,
	})
	if err != nil {
		return TranslationResponse{}, err
	}
	return res, apiError(resp)
}
