package lokalise

import (
	"fmt"
)

const (
	pathOrders = "orders"
)

type OrderService struct {
	BaseService
}

// ‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// Service entity objects
// _____________________________________________________________________________________________________________________

type Order struct {
	WithProjectID
	WithCreationTime
	WithCreationUser

	OrderID             string           `json:"order_id"`
	CardID              int64            `json:"card_id"`
	Status              string           `json:"status"`
	SourceLangISO       string           `json:"source_language_iso"`
	TargetLangISOs      []string         `json:"target_language_isos"`
	Keys                []int64          `json:"keys"`
	SourceWords         map[string]int64 `json:"source_words"`
	ProviderSlug        string           `json:"provider_slug"`
	TranslationStyle    string           `json:"translation_style,omitempty"`
	TranslationTierID   int64            `json:"translation_tier"`
	TranslationTierName string           `json:"translation_tier_name"`
	Briefing            string           `json:"briefing,omitempty"`
	Total               float64          `json:"total"`
}

// ‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// Service request/response objects
// _____________________________________________________________________________________________________________________

type CreateOrder struct {
	ProjectID         string   `json:"project_id"`
	CardID            int64    `json:"card_id"`
	Briefing          string   `json:"briefing"`
	SourceLangISO     string   `json:"source_language_iso"`
	TargetLangISOs    []string `json:"target_language_isos"`
	Keys              []int    `json:"keys"`
	ProviderSlug      string   `json:"provider_slug"`
	TranslationTierID int64    `json:"translation_tier"`
	DryRun            bool     `json:"dry_run,omitempty"`
	TranslationStyle  string   `json:"translation_style,omitempty"`
}

type OrdersResponse struct {
	Paged
	Orders []Order `json:"orders"`
}

// ‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// Service methods
// _____________________________________________________________________________________________________________________

func (c *OrderService) List(teamID int64) (r OrdersResponse, err error) {
	resp, err := c.getList(c.Ctx(), fmt.Sprintf("%s/%d/%s", pathTeams, teamID, pathOrders), &r, c.PageOpts())

	if err != nil {
		return
	}
	applyPaged(resp, &r.Paged)
	return r, apiError(resp)
}

func (c *OrderService) Create(teamID int64, order CreateOrder) (r Order, err error) {
	resp, err := c.post(c.Ctx(), fmt.Sprintf("%s/%d/%s", pathTeams, teamID, pathOrders), &r, order)

	if err != nil {
		return
	}
	return r, apiError(resp)
}

func (c *OrderService) Retrieve(teamID int64, orderID string) (r Order, err error) {
	resp, err := c.get(c.Ctx(), fmt.Sprintf("%s/%d/%s/%s", pathTeams, teamID, pathOrders, orderID), &r)

	if err != nil {
		return
	}
	return r, apiError(resp)
}
