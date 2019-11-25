package lokalise

import (
	"fmt"
)

const (
	pathWebhooks = "webhooks"
)

type WebhookService struct {
	BaseService
}

type Webhook struct {
	WebhookID    string      `json:"webhook_id"`
	URL          string      `json:"url"`
	Secret       string      `json:"secret"`
	Events       []string    `json:"events"`
	EventLangMap []EventLang `json:"event_lang_map,omitempty"`
}

type CreateWebhook struct {
	URL          string      `json:"url"`
	Events       []string    `json:"events"`
	EventLangMap []EventLang `json:"event_lang_map,omitempty"`
}

type UpdateWebhook struct { // all optional
	URL          string      `json:"url,omitempty"`
	Events       []string    `json:"events,omitempty"`
	EventLangMap []EventLang `json:"event_lang_map,omitempty"`
}

type EventLang struct {
	Event        string   `json:"event"`
	LangISOCodes []string `json:"lang_iso_codes"`
}

type WebhooksResponse struct {
	Paged
	WithProjectID
	Webhooks []Webhook `json:"webhooks"`
}

type WebhookResponse struct {
	WithProjectID
	Webhook Webhook `json:"webhook"`
}

type DeleteWebhookResponse struct {
	WithProjectID
	Deleted bool `json:"webhook_deleted"`
}

func (c *WebhookService) List(projectID string) (r WebhooksResponse, err error) {
	resp, err := c.getList(c.Ctx(), fmt.Sprintf("%s/%s/%s", pathProjects, projectID, pathWebhooks), &r, c.PageOpts())

	if err != nil {
		return
	}
	applyPaged(resp, &r.Paged)
	return r, apiError(resp)
}

func (c *WebhookService) Create(projectID string, wh CreateWebhook) (r WebhookResponse, err error) {
	resp, err := c.post(c.Ctx(), fmt.Sprintf("%s/%s/%s", pathProjects, projectID, pathWebhooks), &r, wh)

	if err != nil {
		return
	}
	return r, apiError(resp)
}

func (c *WebhookService) Update(projectID string, webhookID string, opts UpdateWebhook) (r WebhookResponse, err error) {
	resp, err := c.put(c.Ctx(), fmt.Sprintf("%s/%s/%s/%s", pathProjects, projectID, pathWebhooks, webhookID), &r, opts)

	if err != nil {
		return
	}
	return r, apiError(resp)
}

func (c *WebhookService) Retrieve(projectID string, webhookID string) (r WebhookResponse, err error) {
	resp, err := c.get(c.Ctx(), fmt.Sprintf("%s/%s/%s/%s", pathProjects, projectID, pathWebhooks, webhookID), &r)

	if err != nil {
		return
	}
	return r, apiError(resp)
}

func (c *WebhookService) Delete(projectID string, webhookID string) (r DeleteWebhookResponse, err error) {
	resp, err := c.delete(c.Ctx(), fmt.Sprintf("%s/%s/%s/%s", pathProjects, projectID, pathWebhooks, webhookID), &r)

	if err != nil {
		return
	}
	return r, apiError(resp)
}
