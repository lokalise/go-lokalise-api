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

type Order struct {
	WithProjectID
	WithCreationTime
	OrderID             int64            `json:"order_id"`
	CardID              int64            `json:"card_id"`
	Status              string           `json:"status"`
	CreatedBy           int64            `json:"created_by,omitempty"`
	CreatedByEmail      string           `json:"created_by_email,omitempty"`
	SourceLangISO       string           `json:"source_language_iso"`
	TargetLangISOs      []string         `json:"target_language_isos"`
	Keys                []string         `json:"keys,omitempty"`
	SourceWords         map[string]int64 `json:"source_words"`
	ProviderSlug        string           `json:"provider_slug"`
	TranslationStyle    string           `json:"translation_style,omitempty"`
	TranslationTierID   int64            `json:"translation_tier"` // this should be TranslationTier.TierID
	TranslationTierName string           `json:"translation_tier_name"`
	Briefing            string           `json:"briefing,omitempty"`
	Total               float64          `json:"total"`
}

type CreateOrder struct {
	WithProjectID
	CardID            int64    `json:"card_id"`
	Briefing          string   `json:"briefing,omitempty"`
	SourceLangISO     string   `json:"source_language_iso"`
	TargetLangISOs    []string `json:"target_language_isos"`
	Keys              []string `json:"keys,omitempty"`
	ProviderSlug      string   `json:"provider_slug"`
	TranslationTierID int64    `json:"translation_tier"`
	DryRun            bool     `json:"dry_run"`
	TranslationStyle  string   `json:"translation_style,omitempty"`
}

type OrdersResponse struct {
	Paged
	Orders []Order `json:"orders"`
}

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

func (c *OrderService) Retrieve(teamID, orderID int64) (r Order, err error) {
	resp, err := c.get(c.Ctx(), fmt.Sprintf("%s/%d/%s/%d", pathTeams, teamID, pathOrders, orderID), &r)

	if err != nil {
		return
	}
	return r, apiError(resp)
}
