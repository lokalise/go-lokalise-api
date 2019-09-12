package lokalise

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"
)

const (
	pathKeys = "keys"
)

type PlatformStrings struct {
	Ios     string `json:"ios,omitempty"`
	Android string `json:"android,omitempty"`
	Web     string `json:"web,omitempty"`
	Other   string `json:"other,omitempty"`
}

type CustomAttributes struct {
	Attributes map[string]interface{}
}

type Key struct { // todo pointers for Update method
	KeyID            int64             `json:"key_id,omitempty"`
	CreatedAt        string            `json:"created_at,omitempty"`
	KeyName          interface{}       `json:"key_name,omitempty"` // KeyName could be string or PlatformStrings
	Filenames        PlatformStrings   `json:"filenames,omitempty"`
	Description      string            `json:"description,omitempty"`
	Platforms        []string          `json:"platforms,omitempty"`
	Tags             []string          `json:"tags,omitempty"`
	Comments         []Comment         `json:"comments,omitempty"`
	Screenshots      []Screenshot      `json:"screenshots,omitempty"`
	Translations     []Translation     `json:"translations,omitempty"`
	IsPlural         bool              `json:"is_plural,omitempty"`
	PluralName       string            `json:"plural_name,omitempty"`
	IsHidden         bool              `json:"is_hidden,omitempty"`
	IsArchived       bool              `json:"is_archived,omitempty"`
	Context          string            `json:"context,omitempty"`
	CharLimit        int               `json:"char_limit,omitempty"`
	CustomAttributes *CustomAttributes `json:"custom_attributes,string,omitempty"`
}

func (ca *CustomAttributes) UnmarshalJSON(data []byte) error {

	caJsonString := ""
	var customAttributes CustomAttributes

	// First unmarshal the data to a string.
	if err := json.Unmarshal(data, &caJsonString); err != nil {
		return err
	}

	// Escape if the string is empty
	if caJsonString == "" {
		return nil
	}

	// Unmarshal the string further into a map[string]interface{} structure
	if err := json.Unmarshal([]byte(caJsonString), &customAttributes.Attributes); err != nil {
		return err
	}
	ca.Attributes = customAttributes.Attributes

	return nil
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

type KeyResponse struct {
	ProjectID string `json:"project_id,omitempty"`
	Key       Key    `json:"key,omitempty"`
}

type DeleteKeyResponse struct {
	ProjectID      string `json:"project_id,omitempty"`
	IsRemoved      bool   `json:"key_removed"`
	NumberOfLocked int64  `json:"keys_locked"`
}

type DeleteKeysResponse struct {
	ProjectID      string `json:"project_id,omitempty"`
	AreRemoved     bool   `json:"keys_removed"`
	NumberOfLocked int64  `json:"keys_locked"`
}

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

type RetrieveKeyOptions struct {
	DisableReferences bool `json:"disable_references"`
}

func (options RetrieveKeyOptions) Apply(req *resty.Request) {
	if options.DisableReferences {
		req.SetQueryParam("disable_references", "1")
	}
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

func (c *KeysService) List(ctx context.Context, projectID string, options ListKeysOptions) (KeysResponse, error) {
	var res KeysResponse
	resp, err := c.client.getList(ctx, fmt.Sprintf("%s/%s/%s", pathProjects, projectID, pathKeys), &res, options)
	if err != nil {
		return KeysResponse{}, err
	}
	applyPaged(resp, &res.Paged)
	return res, apiError(resp)
}

func (c *KeysService) Create(ctx context.Context, projectID string, keys []Key) (KeysResponse, error) {
	var res KeysResponse
	resp, err := c.client.post(ctx, fmt.Sprintf("%s/%s/%s", pathProjects, projectID, pathKeys), &res, map[string]interface{}{
		"keys": keys,
	})
	if err != nil {
		return KeysResponse{}, err
	}
	return res, apiError(resp)
}

func (c *KeysService) Retrieve(ctx context.Context, projectID string, keyID int64, options RetrieveKeyOptions) (KeyResponse, error) {
	var res KeyResponse
	resp, err := c.client.getList(ctx, fmt.Sprintf("%s/%s/%s/%d", pathProjects, projectID, pathKeys, keyID), &res, options)
	if err != nil {
		return KeyResponse{}, err
	}
	return res, apiError(resp)
}

func (c *KeysService) Update(ctx context.Context, projectID string, keyID int64, key Key) (KeyResponse, error) {
	var res KeyResponse
	resp, err := c.client.put(ctx, fmt.Sprintf("%s/%s/%s/%d", pathProjects, projectID, pathKeys, keyID), &res, key)
	if err != nil {
		return KeyResponse{}, err
	}
	return res, apiError(resp)
}

func (c *KeysService) BulkUpdate(ctx context.Context, projectID string, keys []Key) (KeysResponse, error) {
	var res KeysResponse
	resp, err := c.client.put(ctx, fmt.Sprintf("%s/%s/%s", pathProjects, projectID, pathKeys), &res, map[string]interface{}{
		"keys": keys,
	})
	if err != nil {
		return KeysResponse{}, err
	}
	return res, apiError(resp)
}

func (c *KeysService) Delete(ctx context.Context, projectID string, keyID int64) (DeleteKeyResponse, error) {
	var res DeleteKeyResponse
	resp, err := c.client.delete(ctx, fmt.Sprintf("%s/%s/%s/%d", pathProjects, projectID, pathKeys, keyID), &res)
	if err != nil {
		return DeleteKeyResponse{}, err
	}
	return res, apiError(resp)
}

func (c *KeysService) BulkDelete(ctx context.Context, projectID string, keyIDs []int64) (DeleteKeysResponse, error) {
	var res DeleteKeysResponse
	resp, err := c.client.deleteWithBody(ctx, fmt.Sprintf("%s/%s/%s", pathProjects, projectID, pathKeys), &res, map[string]interface{}{
		"keys": keyIDs,
	})
	if err != nil {
		return DeleteKeysResponse{}, err
	}
	return res, apiError(resp)
}
