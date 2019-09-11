package lokalise

import (
	"context"
	"path"
	"strconv"

	"github.com/go-resty/resty/v2"
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
	ProjectID string     `json:"project_id,omitempty"`
	Languages []Language `json:"languages,omitempty"`
}

type CustomLanguage struct {
	LangISO           string   `json:"lang_iso"`
	CustomISO         string   `json:"custom_iso,omitempty"`
	CustomName        string   `json:"custom_name,omitempty"`
	CustomPluralForms []string `json:"custom_plural_forms,omitempty"`
}

type CreateLanguageResponse struct {
	ProjectID string     `json:"project_id,omitempty"`
	Languages []Language `json:"languages,omitempty"`
}

type RetrieveLanguageResponse struct {
	ProjectID string   `json:"project_id,omitempty"`
	Language  Language `json:"language,omitempty"`
}

type UpdateLanguageResponse struct {
	ProjectID string   `json:"project_id,omitempty"`
	Language  Language `json:"language,omitempty"`
}

type DeleteLanguageResponse struct {
	ProjectID       string `json:"project_id,omitempty"`
	LanguageDeleted bool   `json:"language_deleted,omitempty"`
}

type LanguagesService struct {
	client *Client
}

type LanguagesOptions struct {
	PageOptions
}

func (options LanguagesOptions) Apply(req *resty.Request) {
	options.PageOptions.Apply(req)
}

func (c *LanguagesService) ListSystem(ctx context.Context, pageOptions LanguagesOptions) (result ListLanguagesResponse, err error) {
	url := path.Join("system", pathLanguages)
	resp, err := c.client.getList(ctx, url, &result, pageOptions)
	if err != nil {
		return result, err
	}

	applyPaged(resp, &result.Paged)
	return result, apiError(resp)
}

func (c *LanguagesService) ListProject(ctx context.Context, projectID string, pageOptions LanguagesOptions) (result ListLanguagesResponse, err error) {
	url := path.Join(pathProjects, projectID, pathLanguages)
	resp, err := c.client.getList(ctx, url, &result, pageOptions)
	if err != nil {
		return result, err
	}

	applyPaged(resp, &result.Paged)
	return result, apiError(resp)
}

func (c *LanguagesService) Create(ctx context.Context, projectID string, languages []CustomLanguage) (result CreateLanguageResponse, err error) {
	url := path.Join(pathProjects, projectID, pathLanguages)
	resp, err := c.client.post(ctx, url, &result, map[string]interface{}{"languages": languages})
	if err != nil {
		return result, err
	}
	return result, apiError(resp)
}

func (c *LanguagesService) Retrieve(ctx context.Context, projectID string, ID int64) (result RetrieveLanguageResponse, err error) {
	url := path.Join(pathProjects, projectID, pathLanguages, strconv.FormatInt(ID, 10))
	resp, err := c.client.get(ctx, url, &result)
	if err != nil {
		return result, err
	}
	return result, apiError(resp)
}

func (c *LanguagesService) Update(ctx context.Context, projectID string, ID int64, language Language) (result UpdateLanguageResponse, err error) {
	url := path.Join(pathProjects, projectID, pathLanguages, strconv.FormatInt(ID, 10))
	resp, err := c.client.put(ctx, url, &result, language)
	if err != nil {
		return result, err
	}
	return result, apiError(resp)
}

func (c *LanguagesService) Delete(ctx context.Context, projectID string, ID int64) (result DeleteLanguageResponse, err error) {
	url := path.Join(pathProjects, projectID, pathLanguages, strconv.FormatInt(ID, 10))
	resp, err := c.client.delete(ctx, url, &result)
	if err != nil {
		return result, err
	}
	return result, apiError(resp)
}
