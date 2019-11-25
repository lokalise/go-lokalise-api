package lokalise

import (
	"path"
	"strconv"
)

const (
	pathLanguages = "languages"
)

type LanguageService struct {
	BaseService
}

// ‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// Service entity objects
// _____________________________________________________________________________________________________________________

type Language struct {
	LangID      int64    `json:"lang_id,omitempty"`
	LangISO     string   `json:"lang_iso"`
	LangName    string   `json:"lang_name,omitempty"`
	IsRTL       bool     `json:"is_rtl,omitempty"`
	IsWritable  bool     `json:"is_writable"`
	PluralForms []string `json:"plural_forms,omitempty"`
}

// ‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// Service request/response objects
// _____________________________________________________________________________________________________________________

type ListLanguagesResponse struct {
	Paged
	WithProjectID
	Languages []Language `json:"languages"`
}

type NewLanguage struct {
	LangISO           string   `json:"lang_iso"`
	CustomISO         string   `json:"custom_iso,omitempty"`
	CustomName        string   `json:"custom_name,omitempty"`
	CustomPluralForms []string `json:"custom_plural_forms,omitempty"`
}

type UpdateLanguage struct {
	LangISO     string   `json:"lang_iso,omitempty"`
	LangName    string   `json:"lang_name,omitempty"`
	PluralForms []string `json:"plural_forms,omitempty"`
}

type CreateLanguageResponse struct {
	WithProjectID
	Languages []Language `json:"languages"`
}

type RetrieveLanguageResponse struct {
	WithProjectID
	Language Language `json:"language"`
}

type UpdateLanguageResponse struct {
	WithProjectID
	Language Language `json:"language"`
}

type DeleteLanguageResponse struct {
	WithProjectID
	LanguageDeleted bool `json:"language_deleted"`
}

// ‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// Service methods
// _____________________________________________________________________________________________________________________

func (c *LanguageService) ListSystem() (r ListLanguagesResponse, err error) {
	url := path.Join("system", pathLanguages)
	resp, err := c.getList(c.Ctx(), url, &r, c.PageOpts())

	if err != nil {
		return
	}
	applyPaged(resp, &r.Paged)
	return r, apiError(resp)
}

func (c *LanguageService) ListProject(projectID string) (r ListLanguagesResponse, err error) {
	url := path.Join(pathProjects, projectID, pathLanguages)
	resp, err := c.getList(c.Ctx(), url, &r, c.PageOpts())

	if err != nil {
		return
	}
	applyPaged(resp, &r.Paged)
	return r, apiError(resp)
}

func (c *LanguageService) Create(projectID string, languages []NewLanguage) (r CreateLanguageResponse, err error) {
	url := path.Join(pathProjects, projectID, pathLanguages)
	resp, err := c.post(c.Ctx(), url, &r, map[string]interface{}{"languages": languages})

	if err != nil {
		return
	}
	return r, apiError(resp)
}

func (c *LanguageService) Retrieve(projectID string, ID int64) (r RetrieveLanguageResponse, err error) {
	url := path.Join(pathProjects, projectID, pathLanguages, strconv.FormatInt(ID, 10))
	resp, err := c.get(c.Ctx(), url, &r)

	if err != nil {
		return
	}
	return r, apiError(resp)
}

func (c *LanguageService) Update(projectID string, ID int64, language UpdateLanguage) (r UpdateLanguageResponse, err error) {
	url := path.Join(pathProjects, projectID, pathLanguages, strconv.FormatInt(ID, 10))
	resp, err := c.put(c.Ctx(), url, &r, language)

	if err != nil {
		return
	}
	return r, apiError(resp)
}

func (c *LanguageService) Delete(projectID string, ID int64) (r DeleteLanguageResponse, err error) {
	url := path.Join(pathProjects, projectID, pathLanguages, strconv.FormatInt(ID, 10))
	resp, err := c.delete(c.Ctx(), url, &r)

	if err != nil {
		return
	}
	return r, apiError(resp)
}
