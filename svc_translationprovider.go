package lokalise

import "fmt"

const (
	pathTranslationProviders = "translation_providers"
)

type TranslationProviderService struct {
	BaseService
}

type TranslationProvider struct {
	ProviderID   int64             `json:"provider_id"`
	Name         string            `json:"name"`
	Slug         string            `json:"slug"`
	PricePairMin string            `json:"price_pair_min"`
	WebsiteURL   string            `json:"website_url,omitempty"`
	Description  string            `json:"description,omitempty"`
	Tiers        []TranslationTier `json:"tiers,omitempty"`
	Pairs        []TranslationPair `json:"pairs"`
}

type TranslationTier struct {
	TierID int64  `json:"tier_id"`
	Title  string `json:"title"`
}

type TranslationPair struct {
	TierID       int64  `json:"tier_id"`
	FromLangISO  string `json:"from_lang_iso"`
	FromLangName string `json:"from_lang_name"`
	ToLangISO    string `json:"to_lang_iso"`
	ToLangName   string `json:"to_lang_name"`
	PricePerWord string `json:"price_per_word"`
}

// ‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// Service request/response objects
// _____________________________________________________________________________________________________________________

type TranslationProvidersResponse struct {
	Paged
	TranslationProviders []TranslationProvider `json:"translation_providers"`
}

// ‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// Service methods
// _____________________________________________________________________________________________________________________

func (c *TranslationProviderService) List(teamID int64) (r TranslationProvidersResponse, err error) {
	resp, err := c.getList(c.Ctx(), fmt.Sprintf("%s/%d/%s", pathTeams, teamID, pathTranslationProviders), &r, c.PageOpts())

	if err != nil {
		return
	}
	applyPaged(resp, &r.Paged)
	return r, apiError(resp)
}

func (c *TranslationProviderService) Retrieve(teamID, providerID int64) (r TranslationProvider, err error) {
	resp, err := c.get(c.Ctx(), fmt.Sprintf("%s/%d/%s/%d", pathTeams, teamID, "translation_providers", providerID), &r)

	if err != nil {
		return
	}
	return r, apiError(resp)
}
