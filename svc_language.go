package lokalise

import (
	"path"
	"strconv"
)

const (
	pathLanguages = "languages"
)

type Language struct {
	LangID      int64    `json:"lang_id,omitempty"`
	LangISO     string   `json:"lang_iso,omitempty"`
	LangName    string   `json:"lang_name,omitempty"`
	IsRTL       bool     `json:"is_rtl,omitempty"`
	IsWritable  bool     `json:"is_writable,omitempty"`
	PluralForms []string `json:"plural_forms,omitempty"`
}

type ListLanguagesResponse struct {
	Paged
	WithProjectID
	Languages []Language `json:"languages,omitempty"`
}

type CustomLanguage struct {
	LangISO           string   `json:"lang_iso"`
	CustomISO         string   `json:"custom_iso,omitempty"`
	CustomName        string   `json:"custom_name,omitempty"`
	CustomPluralForms []string `json:"custom_plural_forms,omitempty"`
}

type CreateLanguageResponse struct {
	WithProjectID
	Languages []Language `json:"languages,omitempty"`
}

type RetrieveLanguageResponse struct {
	WithProjectID
	Language Language `json:"language,omitempty"`
}

type UpdateLanguageResponse struct {
	WithProjectID
	Language Language `json:"language,omitempty"`
}

type DeleteLanguageResponse struct {
	WithProjectID
	LanguageDeleted bool `json:"language_deleted,omitempty"`
}

type LanguagesService struct {
	BaseService
}

func (c *LanguagesService) ListSystem() (r ListLanguagesResponse, err error) {
	url := path.Join("system", pathLanguages)
	resp, err := c.getList(c.Ctx(), url, &r, c.PageOpts())

	if err != nil {
		return
	}
	applyPaged(resp, &r.Paged)
	return r, apiError(resp)
}

func (c *LanguagesService) ListProject(projectID string) (r ListLanguagesResponse, err error) {
	url := path.Join(pathProjects, projectID, pathLanguages)
	resp, err := c.getList(c.Ctx(), url, &r, c.PageOpts())

	if err != nil {
		return
	}
	applyPaged(resp, &r.Paged)
	return r, apiError(resp)
}

func (c *LanguagesService) Create(projectID string, languages []CustomLanguage) (r CreateLanguageResponse, err error) {
	url := path.Join(pathProjects, projectID, pathLanguages)
	resp, err := c.post(c.Ctx(), url, &r, map[string]interface{}{"languages": languages})

	if err != nil {
		return
	}
	return r, apiError(resp)
}

func (c *LanguagesService) Retrieve(projectID string, ID int64) (r RetrieveLanguageResponse, err error) {
	url := path.Join(pathProjects, projectID, pathLanguages, strconv.FormatInt(ID, 10))
	resp, err := c.get(c.Ctx(), url, &r)

	if err != nil {
		return
	}
	return r, apiError(resp)
}

func (c *LanguagesService) Update(projectID string, ID int64, language Language) (r UpdateLanguageResponse, err error) {
	url := path.Join(pathProjects, projectID, pathLanguages, strconv.FormatInt(ID, 10))
	resp, err := c.put(c.Ctx(), url, &r, language)

	if err != nil {
		return
	}
	return r, apiError(resp)
}

func (c *LanguagesService) Delete(projectID string, ID int64) (r DeleteLanguageResponse, err error) {
	url := path.Join(pathProjects, projectID, pathLanguages, strconv.FormatInt(ID, 10))
	resp, err := c.delete(c.Ctx(), url, &r)

	if err != nil {
		return
	}
	return r, apiError(resp)
}
