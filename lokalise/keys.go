package lokalise

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-resty/resty"

	"github.com/17media/go-lokalise-api/model"
)

const (
	pathKeys = "keys"
)

type KeysService struct {
	client *Client
}

type ListKeysOptions struct {
	PageOptions
	IncludeTranslations       bool
	DisableReferences         bool
	IncludeComments           bool
	IncludeScreenshots        bool
	filterTags                []string
	filterKeys                []string
	filterKeyIDs              []string
	filterPlatforms           []string
	filterPlaceholderMismatch bool
}

func (options ListKeysOptions) Apply(req *resty.Request) {
	options.PageOptions.Apply(req)
	if options.IncludeTranslations {
		req.SetQueryParam("include_translations", "1")
	}
	if options.DisableReferences {
		req.SetQueryParam("disable_references", "1")
	}
	if options.IncludeComments {
		req.SetQueryParam("include_comments", "1")
	}
	if options.IncludeScreenshots {
		req.SetQueryParam("include_screenshots", "1")
	}
	if len(options.filterTags) > 0 {
		req.SetQueryParam("filter_tags", strings.Join(options.filterTags, ","))
	}
	if len(options.filterKeys) > 0 {
		req.SetQueryParam("filter_keys", strings.Join(options.filterKeys, ","))
	}
	if len(options.filterKeyIDs) > 0 {
		req.SetQueryParam("filter_key_ids", strings.Join(options.filterKeyIDs, ","))
	}
	if len(options.filterPlatforms) > 0 {
		req.SetQueryParam("filter_platforms", strings.Join(options.filterPlatforms, ","))
	}
	if options.filterPlaceholderMismatch {
		req.SetQueryParam("filter_placeholder_mismatch", "1")
	}
}

func (c *KeysService) List(ctx context.Context, projectID string, options ListKeysOptions) (model.KeysResponse, error) {
	var res model.KeysResponse
	resp, err := c.client.getList(ctx, fmt.Sprintf("%s/%s/%s", pathProjects, projectID, pathKeys), &res, options)
	if err != nil {
		return model.KeysResponse{}, err
	}
	applyPaged(resp, &res.Paged)
	return res, apiError(resp)
}

func (c *KeysService) Create(ctx context.Context, projectID string, keys []model.Key) (model.KeysResponse, error) {
	var res model.KeysResponse
	resp, err := c.client.post(ctx, fmt.Sprintf("%s/%s/%s", pathProjects, projectID, pathKeys), &res, map[string]interface{}{
		"keys": keys,
	})
	if err != nil {
		return model.KeysResponse{}, err
	}
	return res, apiError(resp)
}
