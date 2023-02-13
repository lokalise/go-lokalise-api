package lokalise

import (
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/google/go-querystring/query"
)

const (
	pathProjects = "projects"
)

// The Project service
type ProjectService struct {
	BaseService

	opts ProjectListOptions
}

// ‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// Service entity objects
// _____________________________________________________________________________________________________________________

type Project struct {
	WithProjectID
	WithTeamID
	WithCreationTime
	WithCreationUser

	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"project_type,omitempty"`

	BaseLangID  int64  `json:"base_language_id"`
	BaseLangISO string `json:"base_language_iso"`

	Settings   *ProjectSettings   `json:"settings"`
	Statistics *ProjectStatistics `json:"statistics"`
}

type ProjectSettings struct {
	PerPlatformKeyNames       bool `json:"per_platform_key_names"`
	Reviewing                 bool `json:"reviewing"`
	Upvoting                  bool `json:"upvoting"`
	AutoToggleUnverified      bool `json:"auto_toggle_unverified"`
	OfflineTranslation        bool `json:"offline_translation"`
	KeyEditing                bool `json:"key_editing"`
	InlineMachineTranslations bool `json:"inline_machine_translations"`
}

type QAIssues struct {
	NotReviewed                   int64 `json:"not_reviewed"`
	Unverified                    int64 `json:"unverified"`
	SpellingGrammar               int64 `json:"spelling_grammar"`
	InconsistentPlaceholders      int64 `json:"inconsistent_placeholders"`
	InconsistentHtml              int64 `json:"inconsistent_html"`
	DifferentNumberOfUrls         int64 `json:"different_number_of_urls"`
	DifferentUrls                 int64 `json:"different_urls"`
	LeadingWhitespace             int64 `json:"leading_whitespace"`
	TrailingWhitespace            int64 `json:"trailing_whitespace"`
	DifferentNumberOfEmailAddress int64 `json:"different_number_of_email_address"`
	DifferentEmailAddress         int64 `json:"different_email_address"`
	DifferentBrackets             int64 `json:"different_brackets"`
	DifferentNumbers              int64 `json:"different_numbers"`
	DoubleSpace                   int64 `json:"double_space"`
	SpecialPlaceholder            int64 `json:"special_placeholder"`
}

type ProjectStatistics struct {
	ProgressTotal int64                `json:"progress_total"`
	KeysTotal     int64                `json:"keys_total"`
	Team          int64                `json:"team"`
	BaseWords     int64                `json:"base_words"`
	QAIssuesTotal int64                `json:"qa_issues_total"`
	QAIssues      QAIssues             `json:"qa_issues"`
	Languages     []LanguageStatistics `json:"languages"`
}

type LanguageStatistics struct {
	LanguageID  int64  `json:"language_id"`
	LanguageISO string `json:"language_iso"`
	Progress    int64  `json:"progress"`
	WordsToDo   int64  `json:"words_to_do"`
}

// ‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// Service request/response objects
// _____________________________________________________________________________________________________________________

type NewProject struct {
	WithTeamID // optional

	Name        string `json:"name"`
	Description string `json:"description,omitempty"`

	Languages   []NewLanguage `json:"languages,omitempty"`
	BaseLangISO string        `json:"base_lang_iso,omitempty"` // Default: en // Differs from Project struct!

	// Allowed values are localization_files, paged_documents
	// Default: localization_files
	ProjectType string `json:"project_type,omitempty"`
}

type UpdateProject struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

type ProjectsResponse struct {
	Paged
	Projects []Project `json:"projects"`
}

type TruncateProjectResponse struct {
	WithProjectID
	KeysDeleted bool `json:"keys_deleted"`
}

type DeleteProjectResponse struct {
	WithProjectID
	Deleted bool `json:"project_deleted"`
}

// ‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// Service methods
// _____________________________________________________________________________________________________________________

// Retrieves a list of projects available to the user
func (c *ProjectService) List() (r ProjectsResponse, err error) {
	resp, err := c.getWithOptions(c.Ctx(), pathProjects, &r, c.ListOpts())

	if err != nil {
		return
	}
	applyPaged(resp, &r.Paged)
	return r, apiError(resp)
}

// Creates a new project in the specified team. Requires Admin role in the team.
func (c *ProjectService) Create(project NewProject) (r Project, err error) {
	resp, err := c.post(c.Ctx(), pathProjects, &r, project)

	if err != nil {
		return
	}
	return r, apiError(resp)
}

func (c *ProjectService) Retrieve(projectID string) (r Project, err error) {
	resp, err := c.get(c.Ctx(), fmt.Sprintf("%s/%s", pathProjects, projectID), &r)

	if err != nil {
		return
	}
	return r, apiError(resp)
}

func (c *ProjectService) Update(projectID string, project UpdateProject) (r Project, err error) {
	resp, err := c.put(c.Ctx(), fmt.Sprintf("%s/%s", pathProjects, projectID), &r, project)

	if err != nil {
		return
	}
	return r, apiError(resp)
}

func (c *ProjectService) Truncate(projectID string) (r TruncateProjectResponse, err error) {
	resp, err := c.put(c.Ctx(), fmt.Sprintf("%s/%s/empty", pathProjects, projectID), &r, nil)

	if err != nil {
		return
	}
	return r, apiError(resp)
}

func (c *ProjectService) Delete(projectID string) (r DeleteProjectResponse, err error) {
	resp, err := c.delete(c.Ctx(), fmt.Sprintf("%s/%s", pathProjects, projectID), &r)

	if err != nil {
		return
	}
	return r, apiError(resp)
}

// ‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// Additional service methods
// _____________________________________________________________________________________________________________________

// List options
type ProjectListOptions struct {
	// page options
	Page  uint `url:"page,omitempty"`
	Limit uint `url:"limit,omitempty"`

	FilterTeamID int64  `url:"filter_team_id,omitempty"`
	FilterNames  string `url:"filter_names,omitempty"`

	// Possible values are 1 and 0, default 1
	IncludeStat     string `url:"include_statistics,omitempty"`
	IncludeSettings string `url:"include_settings,omitempty"`
}

func (options ProjectListOptions) Apply(req *resty.Request) {
	v, _ := query.Values(options)
	req.SetQueryString(v.Encode())
}

func (c *ProjectService) ListOpts() ProjectListOptions {
	return c.opts
}

func (c *ProjectService) WithListOptions(o ProjectListOptions) *ProjectService {
	c.opts = o
	return c
}

func (c *ProjectService) SetListOptions(o ProjectListOptions) {
	c.opts = o
}
