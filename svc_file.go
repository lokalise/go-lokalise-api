package lokalise

import (
	"fmt"
	"github.com/google/go-querystring/query"

	"github.com/go-resty/resty/v2"
)

const (
	pathFiles = "files"
)

type FileService struct {
	BaseService

	opts FileListOptions
}

// ‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// Service entity objects
// _____________________________________________________________________________________________________________________

type File struct {
	Filename string `json:"filename"`
	KeyCount int64  `json:"key_count"`
}

type FileUpload struct {
	Data     string   `json:"data"`
	Filename string   `json:"filename"`
	LangISO  string   `json:"lang_iso"`
	Tags     []string `json:"tags,omitempty"`

	ConvertPlaceholders                 *bool   `json:"convert_placeholders,omitempty"`
	DetectICUPlurals                    bool    `json:"detect_icu_plurals,omitempty"`
	TagInsertedKeys                     *bool   `json:"tag_inserted_keys,omitempty"`
	TagUpdatedKeys                      *bool   `json:"tag_updated_keys,omitempty"`
	TagSkippedKeys                      bool    `json:"tag_skipped_keys,omitempty"`
	ReplaceModified                     bool    `json:"replace_modified,omitempty"`
	SlashNToLinebreak                   *bool   `json:"slashn_to_linebreak,omitempty"`
	KeysToValues                        bool    `json:"keys_to_values,omitempty"`
	DistinguishByFile                   bool    `json:"distinguish_by_file,omitempty"`
	ApplyTM                             bool    `json:"apply_tm,omitempty"`
	HiddenFromContributors              bool    `json:"hidden_from_contributors,omitempty"`
	CleanupMode                         bool    `json:"cleanup_mode,omitempty"`
	CustomTranslationStatusIds          []int64 `json:"custom_translation_status_ids,omitempty"`
	CustomTranslationStatusInsertedKeys *bool   `json:"custom_translation_status_inserted_keys,omitempty"`
	CustomTranslationStatusUpdatedKeys  *bool   `json:"custom_translation_status_updated_keys,omitempty"`
	CustomTranslationStatusSkippedKeys  *bool   `json:"custom_translation_status_skipped_keys,omitempty"`
	Queue                               bool    `json:"queue"`
	SkipDetectLangIso                   bool    `json:"skip_detect_lang_iso,omitempty"`
	UseAutomations                      *bool   `json:"use_automations,omitempty"`
}

type FileDownload struct {
	Format                     string            `json:"format"`
	OriginalFilenames          *bool             `json:"original_filenames,omitempty"`
	BundleStructure            string            `json:"bundle_structure,omitempty"`
	DirectoryPrefix            *string           `json:"directory_prefix,omitempty"`
	AllPlatforms               bool              `json:"all_platforms,omitempty"`
	FilterLangs                []string          `json:"filter_langs,omitempty"`
	FilterData                 []string          `json:"filter_data,omitempty"`
	FilterFilenames            []string          `json:"filter_filenames,omitempty"`
	AddNewlineEOF              bool              `json:"add_newline_eof,omitempty"`
	CustomTranslationStatusIDs []string          `json:"custom_translation_status_ids,omitempty"`
	IncludeTags                []string          `json:"include_tags,omitempty"`
	ExcludeTags                []string          `json:"exclude_tags,omitempty"`
	ExportSort                 string            `json:"export_sort,omitempty"`
	ExportEmptyAs              string            `json:"export_empty_as,omitempty"`
	IncludeComments            bool              `json:"include_comments,omitempty"`
	IncludeDescription         *bool             `json:"include_description,omitempty"`
	IncludeProjectIDs          []string          `json:"include_pids,omitempty"`
	Triggers                   []string          `json:"triggers,omitempty"`
	FilterRepositories         []string          `json:"filter_repositories,omitempty"`
	ReplaceBreaks              *bool             `json:"replace_breaks,omitempty"`
	DisableReferences          bool              `json:"disable_references,omitempty"`
	PluralFormat               string            `json:"plural_format,omitempty"`
	PlaceholderFormat          string            `json:"placeholder_format,omitempty"`
	WebhookURL                 string            `json:"webhook_url,omitempty"`
	LanguageMapping            []LanguageMapping `json:"language_mapping,omitempty"`
	ICUNumeric                 bool              `json:"icu_numeric,omitempty"`
	EscapePercent              bool              `json:"escape_percent,omitempty"`
	Indentation                string            `json:"indentation,omitempty"`
	YAMLIncludeRoot            bool              `json:"yaml_include_root,omitempty"`
	JSONUnescapedSlashes       bool              `json:"json_unescaped_slashes,omitempty"`
	JavaPropertiesEncoding     string            `json:"java_properties_encoding,omitempty"`
	JavaPropertiesSeparator    string            `json:"java_properties_separator,omitempty"`
	BundleDescription          string            `json:"bundle_description,omitempty"`
}

type LanguageMapping struct {
	OriginalLangISO string `json:"original_language_iso"`
	CustomLangISO   string `json:"custom_language_iso"`
}

// ‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// Service request/response objects
// _____________________________________________________________________________________________________________________

type FilesResponse struct {
	Paged
	WithProjectID
	Files []File `json:"files"`
}

type FileUploadResponse struct {
	WithProjectID
	Process QueuedProcess `json:"process"`
}

type FileDownloadResponse struct {
	WithProjectID
	BundleURL string `json:"bundle_url"`
}

// ‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// Service methods
// _____________________________________________________________________________________________________________________

func (c *FileService) List(projectID string) (r FilesResponse, err error) {
	resp, err := c.getWithOptions(c.Ctx(), fmt.Sprintf("%s/%s/%s", pathProjects, projectID, pathFiles), &r, c.ListOpts())

	if err != nil {
		return
	}
	applyPaged(resp, &r.Paged)
	return r, apiError(resp)
}

func (c *FileService) Upload(projectID string, file FileUpload) (r FileUploadResponse, err error) {
	if file.CustomTranslationStatusSkippedKeys == nil {
		file.CustomTranslationStatusSkippedKeys = Bool(false)
	}
	if file.CustomTranslationStatusUpdatedKeys == nil {
		file.CustomTranslationStatusUpdatedKeys = Bool(true)
	}
	if file.CustomTranslationStatusInsertedKeys == nil {
		file.CustomTranslationStatusInsertedKeys = Bool(true)
	}

	file.Queue = true

	resp, err := c.post(c.Ctx(), fmt.Sprintf("%s/%s/%s/%s", pathProjects, projectID, pathFiles, "upload"), &r, file)

	if err != nil {
		return
	}
	return r, apiError(resp)
}

func (c *FileService) Download(projectID string, downloadOptions FileDownload) (r FileDownloadResponse, err error) {
	url := fmt.Sprintf("%s/%s/%s/%s", pathProjects, projectID, pathFiles, "download")
	resp, err := c.post(c.Ctx(), url, &r, downloadOptions)

	if err != nil {
		return
	}
	return r, apiError(resp)
}

// ‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// Additional service structs & methods
// _____________________________________________________________________________________________________________________

type FileListOptions struct {
	Limit    uint   `url:"limit,omitempty"`
	Page     uint   `url:"page,omitempty"`
	Filename string `url:"filter_filename,omitempty"`
}

func (options FileListOptions) Apply(req *resty.Request) {
	v, _ := query.Values(options)
	req.SetQueryString(v.Encode())
}

func (c *FileService) ListOpts() FileListOptions        { return c.opts }
func (c *FileService) SetListOptions(o FileListOptions) { c.opts = o }
func (c *FileService) WithListOptions(o FileListOptions) *FileService {
	c.opts = o
	return c
}
