package lokalise

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

const (
	pathFiles = "files"
)

type FileService struct {
	BaseService
}

type File struct {
	Filename string `json:"filename"`
	KeyCount int64  `json:"key_count"`
}

type FileUpload struct {
	FileUploadOptions

	Data        string   `json:"data"`
	Filename    string   `json:"filename"`
	LangISO     string   `json:"lang_iso"`
	Tags        []string `json:"tags"`
	CleanupMode bool     `json:"cleanup_mode"`
}

type FileUploadOptions struct {
	ConvertPlaceholders    bool `json:"convert_placeholders"`
	DetectICUPlurals       bool `json:"detect_icu_plurals"`
	TagInsertedKeys        bool `json:"tag_inserted_keys"`
	TagUpdatedKeys         bool `json:"tag_updated_keys"`
	TagSkippedKeys         bool `json:"tag_skipped_keys"`
	ReplaceModified        bool `json:"replace_modified"`
	SlashNToLinebreak      bool `json:"slashn_to_linebreak"`
	KeysToValues           bool `json:"keys_to_values"`
	DistinguishByFile      bool `json:"distinguish_by_file"`
	ApplyTM                bool `json:"apply_tm"`
	HiddenFromContributors bool `json:"hidden_from_contributors"`
}

type FileDownloadOptions struct {
	Format string `json:"format"`

	OriginalFilenames          bool              `json:"original_filenames,omitempty"`
	BundleStructure            bool              `json:"bundle_structure,omitempty"`
	DirectoryPrefix            string            `json:"directory_prefix,omitempty"`
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
	IncludeDescription         bool              `json:"include_description,omitempty"`
	IncludePids                bool              `json:"include_pids,omitempty"`
	Triggers                   []string          `json:"triggers,omitempty"`
	FilterRepositories         []string          `json:"filter_repositories,omitempty"`
	ReplaceBreaks              bool              `json:"replace_breaks,omitempty"`
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

type FileOptions struct {
	PageOptions
	Filename string `url:"filter_filename"`
}

func (options FileOptions) Apply(req *resty.Request) {
	options.PageOptions.Apply(req)
	if options.Filename != "" {
		req.SetQueryParam("filter_filename", options.Filename)
	}
}

type FilesResponse struct {
	Paged
	WithProjectID
	Files []File `json:"files"`
}

type FileUploadResponse struct {
	WithProjectID
	Filename string `json:"file"`
	Result   struct {
		Skipped  int64 `json:"skipped,omitempty"`
		Inserted int64 `json:"inserted,omitempty"`
		Updated  int64 `json:"updated,omitempty"`
	} `json:"result"`
}

type FileDownloadResponse struct {
	WithProjectID
	BundleURL string `json:"bundle_url"`
}

func (c *FileService) List(projectID string, opts TasksOptions) (r FilesResponse, err error) {
	resp, err := c.getList(c.Ctx(), fmt.Sprintf("%s/%s/%s", pathProjects, projectID, pathFiles), &r, opts)

	if err != nil {
		return
	}
	applyPaged(resp, &r.Paged)
	return r, apiError(resp)
}

func (c *FileService) Upload(projectID string, file FileUpload) (r FileUploadResponse, err error) {
	resp, err := c.post(c.Ctx(), fmt.Sprintf("%s/%s/%s/%s", pathProjects, projectID, pathFiles, "upload"), &r, file)

	if err != nil {
		return
	}
	return r, apiError(resp)
}

func (c *FileService) Download(projectID string, opts FileDownloadOptions) (r FileDownloadResponse, err error) {
	resp, err := c.post(c.Ctx(), fmt.Sprintf("%s/%s/%s/%s", pathProjects, projectID, pathFiles, "download"), &r, opts)

	if err != nil {
		return
	}
	return r, apiError(resp)
}
