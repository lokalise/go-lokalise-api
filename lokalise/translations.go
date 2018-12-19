package lokalise

import (
	"context"
	"fmt"

	"github.com/17media/go-lokalise-api/model"
	"github.com/go-resty/resty"
)

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

func (c *TranslationsService) List(ctx context.Context, projectID string, pageOptions TranslationsOptions) (model.TranslationsResponse, error) {
	var res model.TranslationsResponse
	resp, err := c.client.getList(ctx, fmt.Sprintf("%s/%s/%s", pathProjects, projectID, pathTranslations), &res, pageOptions)
	if err != nil {
		return model.TranslationsResponse{}, err
	}
	applyPaged(resp, &res.Paged)
	return res, apiError(resp)
}

func (c *TranslationsService) Retrieve(ctx context.Context, projectID string, translationID int64) (model.TranslationResponse, error) {
	var res model.TranslationResponse
	resp, err := c.client.get(ctx, fmt.Sprintf("%s/%s/%s/%d", pathProjects, projectID, pathTranslations, translationID), &res)
	if err != nil {
		return model.TranslationResponse{}, err
	}
	return res, apiError(resp)
}

func (c *TranslationsService) Update(ctx context.Context, projectID string, translationID int64, translation string, isFuzzy, isReviewed bool) (model.TranslationResponse, error) {
	var res model.TranslationResponse
	resp, err := c.client.put(ctx, fmt.Sprintf("%s/%s/%s/%d", pathProjects, projectID, pathTranslations, translationID), &res, map[string]interface{}{
		"translation": translation,
		"is_fuzzy":    isFuzzy,
		"is_reviewed": isReviewed,
	})
	if err != nil {
		return model.TranslationResponse{}, err
	}
	return res, apiError(resp)
}
