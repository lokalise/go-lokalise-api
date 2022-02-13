package lokalise

import (
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

	ModifiedAtTs             int64 `json:"modified_at_timestamp"`
	TranslationsModifiedAtTs int64 `json:"translations_modified_at_timestamp"`
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
	KeyName      interface{}      `json:"key_name,omitempty"` // could be empty in case of updating
	Description  string           `json:"description,omitempty"`
	Platforms    []string         `json:"platforms,omitempty"` // could be empty in case of updating
	Filenames    *PlatformStrings `json:"filenames,omitempty"`
	Tags         []string         `json:"tags,omitempty"`
	MergeTags    bool             `json:"merge_tags,omitempty"`
	Comments     []NewComment     `json:"comments,omitempty"`
	Screenshots  []NewScreenshot  `json:"screenshots,omitempty"`
	Translations []NewTranslation `json:"translations,omitempty"`

	IsPlural         bool   `json:"is_plural,omitempty"`
	PluralName       string `json:"plural_name,omitempty"`
	IsHidden         bool   `json:"is_hidden,omitempty"`
	IsArchived       bool   `json:"is_archived,omitempty"`
	Context          string `json:"context,omitempty"`
	BaseWords        int    `json:"base_words,omitempty"`
	CharLimit        int    `json:"char_limit,omitempty"`
	CustomAttributes string `json:"custom_attributes,omitempty"`
}

type CreateKeysRequest struct {
	Keys []NewKey `json:"keys"`
	KeyRequestOptions
}

// Separate struct for bulk updating
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
	resp, err := c.getList(c.Ctx(), fmt.Sprintf("%s/%s/%s", pathProjects, projectID, pathKeys), &r, c.ListOpts())

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
	resp, err := c.getList(c.Ctx(), fmt.Sprintf("%s/%s/%s/%d", pathProjects, projectID, pathKeys, keyID), &r, c.RetrieveOpts())

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

// ‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// Additional methods
// _____________________________________________________________________________________________________________________

// List options
type KeyListOptions struct {
	// page options
	Page  uint `url:"page,omitempty"`
	Limit uint `url:"limit,omitempty"`

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
