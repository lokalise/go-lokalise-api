package lokalise

import (
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/google/go-querystring/query"
)

const (
	pathKeys = "keys"
)

// The Key service
type KeyService struct {
	BaseService

	listOpts     KeyListOptions
	retrieveOpts KeyRetrieveOptions
}

// ‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// Service entity objects
// _____________________________________________________________________________________________________________________

// The Key object
type Key struct {
	WithCreationTime

	KeyID   int64           `json:"key_id"`
	KeyName PlatformStrings `json:"key_name"`

	Filenames    PlatformStrings `json:"filenames"`
	Description  string          `json:"description"`
	Platforms    []string        `json:"platforms"`
	Tags         []string        `json:"tags"`
	Comments     []Comment       `json:"comments"`
	Screenshots  []Screenshot    `json:"screenshots"`
	Translations []Translation   `json:"translations"`

	IsPlural         bool   `json:"is_plural"`
	PluralName       string `json:"plural_name,omitempty"`
	IsHidden         bool   `json:"is_hidden"`
	IsArchived       bool   `json:"is_archived"`
	Context          string `json:"context,omitempty"`
	BaseWords        int    `json:"base_words"`
	CharLimit        int    `json:"char_limit"`
	CustomAttributes string `json:"custom_attributes,omitempty"`

	ModifiedAt   string `json:"modified_at,omitempty"`
	ModifiedAtTs int64  `json:"modified_at_timestamp,omitempty"`
}

type PlatformStrings struct {
	Ios     string `json:"ios,omitempty"`
	Android string `json:"android,omitempty"`
	Web     string `json:"web,omitempty"`
	Other   string `json:"other,omitempty"`
}

// ‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// Service request/response objects
// _____________________________________________________________________________________________________________________
type KeyRequestOptions struct {
	UseAutomations *bool `json:"use_automations,omitempty"`
}

type KeyRequestOption func(options *KeyRequestOptions)

func WithAutomations(UseAutomations bool) KeyRequestOption {
	return func(c *KeyRequestOptions) {
		c.UseAutomations = &UseAutomations
	}
}

type NewKey struct {
	// KeyName could be string or PlatformStrings
	KeyName          *interface{}
	IsPlural         *bool
	PluralName       *string
	IsHidden         *bool
	IsArchived       *bool
	Context          *string
	BaseWords        *int
	CharLimit        *int
	CustomAttributes *string
	MergeTags        *bool
	Filenames        *PlatformStrings
	Description      *string
	Platforms        *[]string
	Tags             *[]string
	Comments         *[]NewComment
	Screenshots      *[]NewScreenshot
	Translations     *[]NewTranslation
}

type CreateKeysRequest struct {
	Keys []NewKey `json:"keys"`
	KeyRequestOptions
}

// BulkUpdateKey Separate struct for bulk updating
type BulkUpdateKey struct {
	KeyID int64 `json:"key_id"`
	NewKey
}

type BulkUpdateKeysRequest struct {
	Keys []BulkUpdateKey `json:"keys"`
	KeyRequestOptions
}

// ErrorKeys is error for key create/update API
type ErrorKeys struct {
	Error
	Key struct {
		KeyName string `json:"key_name"`
	} `json:"key"`
}

type KeysResponse struct {
	Paged
	WithProjectID
	Keys   []Key       `json:"keys"`
	Errors []ErrorKeys `json:"error,omitempty"`
}

type KeyResponse struct {
	WithProjectID
	Key Key `json:"key"`
}

type DeleteKeyResponse struct {
	WithProjectID
	IsRemoved      bool  `json:"key_removed"`
	NumberOfLocked int64 `json:"keys_locked"`
}

type DeleteKeysResponse struct {
	WithProjectID
	AreRemoved     bool  `json:"keys_removed"`
	NumberOfLocked int64 `json:"keys_locked"`
}

// ‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// Service methods
// _____________________________________________________________________________________________________________________

func (c *KeyService) List(projectID string) (r KeysResponse, err error) {
	resp, err := c.getWithOptions(c.Ctx(), fmt.Sprintf("%s/%s/%s", pathProjects, projectID, pathKeys), &r, c.ListOpts())

	if err != nil {
		return
	}
	applyPaged(resp, &r.Paged)
	return r, apiError(resp)
}

func (c *KeyService) Create(projectID string, keys []NewKey, options ...KeyRequestOption) (r KeysResponse, err error) {
	request := CreateKeysRequest{
		Keys: keys,
	}

	for _, o := range options {
		o(&request.KeyRequestOptions)
	}

	resp, err := c.post(c.Ctx(), fmt.Sprintf("%s/%s/%s", pathProjects, projectID, pathKeys), &r, request)

	if err != nil {
		return
	}
	return r, apiError(resp)
}

func (c *KeyService) Retrieve(projectID string, keyID int64) (r KeyResponse, err error) {
	resp, err := c.getWithOptions(c.Ctx(), fmt.Sprintf("%s/%s/%s/%d", pathProjects, projectID, pathKeys, keyID), &r, c.RetrieveOpts())

	if err != nil {
		return
	}
	return r, apiError(resp)
}

func (c *KeyService) Update(projectID string, keyID int64, key NewKey) (r KeyResponse, err error) {
	resp, err := c.put(c.Ctx(), fmt.Sprintf("%s/%s/%s/%d", pathProjects, projectID, pathKeys, keyID), &r, key)

	if err != nil {
		return
	}
	return r, apiError(resp)
}

func (c *KeyService) BulkUpdate(projectID string, keys []BulkUpdateKey, options ...KeyRequestOption) (r KeysResponse, err error) {
	request := BulkUpdateKeysRequest{
		Keys: keys,
	}

	for _, o := range options {
		o(&request.KeyRequestOptions)
	}

	resp, err := c.put(c.Ctx(), fmt.Sprintf("%s/%s/%s", pathProjects, projectID, pathKeys), &r, request)

	if err != nil {
		return
	}
	return r, apiError(resp)
}

func (c *KeyService) Delete(projectID string, keyID int64) (r DeleteKeyResponse, err error) {
	resp, err := c.delete(c.Ctx(), fmt.Sprintf("%s/%s/%s/%d", pathProjects, projectID, pathKeys, keyID), &r)

	if err != nil {
		return
	}
	return r, apiError(resp)
}

func (c *KeyService) BulkDelete(projectID string, keyIDs []int64) (r DeleteKeysResponse, err error) {
	resp, err := c.deleteWithBody(c.Ctx(), fmt.Sprintf("%s/%s/%s", pathProjects, projectID, pathKeys), &r,
		map[string]interface{}{
			"keys": keyIDs,
		},
	)

	if err != nil {
		return
	}
	return r, apiError(resp)
}

// MarshalJSON Preserve fields for BulkUpdateKey when custom marshaling of anonymous fields are used
func (k BulkUpdateKey) MarshalJSON() ([]byte, error) {
	jsonNewKey, err := json.Marshal(k.NewKey)
	if err != nil {
		return nil, err
	}

	jsonNewKey[0] = ','
	jsonKeyId := []byte(fmt.Sprintf(`{"key_id":%d`, k.KeyID))

	return append(jsonKeyId, jsonNewKey...), nil
}

// MarshalJSON Remove null tags array, preserve empty array in json
func (k NewKey) MarshalJSON() ([]byte, error) {
	data := map[string]interface{}

	if k.KeyName != nil {
		data["key_name"] = k.KeyName
	}

	if k.IsPlural != nil {
		data["is_plural"] = k.IsPlural
	}

	if k.PluralName != nil {
		data["plural_name"] = k.PluralName
	}

	if k.IsHidden != nil {
		data["is_hidden"] = k.IsHidden
	}

	if k.IsArchived != nil {
		data["is_archived"] = k.IsArchived
	}

	if k.Context != nil {
		data["context"] = k.Context
	}

	if k.BaseWords != nil {
		data["base_words"] = k.BaseWords
	}

	if k.CharLimit != nil {
		data["char_limit"] = k.CharLimit
	}

	if k.CustomAttributes != nil {
		data["custom_attributes"] = k.CustomAttributes
	}

	if k.MergeTags != nil {
		data["merge_tags"] = k.MergeTags
	}

	if k.Filenames != nil {
		data["filenames"] = k.Filenames
	}

	if k.Description != nil {
		data["description"] = k.Description
	}

	if k.Platforms != nil {
		data["platforms"] = k.Platforms
	}

	if k.Tags != nil {
		data["tags"] = k.Tags
	}

	if k.Comments != nil {
		data["comments"] = k.Comments
	}

	if k.Screenshots != nil {
		data["screenshots"] = k.Screenshots
	}

	if k.Translations != nil {
		data["translations"] = k.Translations
	}

	return json.Marshal(data)
}

// ‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// Additional methods
// _____________________________________________________________________________________________________________________

// List options
type KeyListOptions struct {
	// page options
	Pagination string `url:"pagination,omitempty"`
	Page       uint   `url:"page,omitempty"`
	Limit      uint   `url:"limit,omitempty"`
	Cursor     string `url:"cursor,omitempty"`

	// Possible values are 1 and 0.
	DisableReferences   uint8 `url:"disable_references,omitempty"`
	IncludeComments     uint8 `url:"include_comments,omitempty"`
	IncludeScreenshots  uint8 `url:"include_screenshots,omitempty"`
	IncludeTranslations uint8 `url:"include_translations,omitempty"`

	FilterTranslationLangIDs string `url:"filter_translation_lang_ids,omitempty"`
	FilterTags               string `url:"filter_tags,omitempty"`
	FilterFilenames          string `url:"filter_filenames,omitempty"`
	FilterKeys               string `url:"filter_keys,omitempty"`
	FilterKeyIDs             string `url:"filter_key_ids,omitempty"`
	FilterPlatforms          string `url:"filter_platforms,omitempty"`
	FilterUntranslated       string `url:"filter_untranslated,omitempty"`
	FilterQAIssues           string `url:"filter_qa_issues,omitempty"`
}

func (options KeyListOptions) Apply(req *resty.Request) {
	v, _ := query.Values(options)
	req.SetQueryString(v.Encode())
}

// Retrieve options
type KeyRetrieveOptions struct {
	DisableReferences uint8 `url:"disable_references,omitempty"`
}

func (options KeyRetrieveOptions) Apply(req *resty.Request) {
	v, _ := query.Values(options)
	req.SetQueryString(v.Encode())
}

func (c *KeyService) ListOpts() KeyListOptions        { return c.listOpts }
func (c *KeyService) SetListOptions(o KeyListOptions) { c.listOpts = o }
func (c *KeyService) WithListOptions(o KeyListOptions) *KeyService {
	c.listOpts = o
	return c
}

func (c *KeyService) RetrieveOpts() KeyRetrieveOptions        { return c.retrieveOpts }
func (c *KeyService) SetRetrieveOptions(o KeyRetrieveOptions) { c.retrieveOpts = o }
func (c *KeyService) WithRetrieveOptions(o KeyRetrieveOptions) *KeyService {
	c.retrieveOpts = o
	return c
}
