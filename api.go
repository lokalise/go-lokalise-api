// Package lokalise provides functions to access the Lokalise web API.
package lokalise

import (
	"errors"
)

type Api struct {
	httpClient *restClient

	Projects             func() *ProjectService
	Teams                func() *TeamService
	TeamUsers            func() *TeamUserService
	TeamUserGroups       func() *TeamUserGroupService
	Contributors         func() *ContributorService
	Comments             func() *CommentService
	Keys                 func() *KeyService
	Tasks                func() *TaskService
	Screenshots          func() *ScreenshotService
	Snapshots            func() *SnapshotService
	Languages            func() *LanguageService
	Translations         func() *TranslationService
	TranslationProviders func() *TranslationProviderService
	TranslationStatuses  func() *TranslationStatusService
	Orders               func() *OrderService
	PaymentCards         func() *PaymentCardService
	Webhooks             func() *WebhookService
	Files                func() *FileService
}

type ClientOption func(*Api) error

func New(apiToken string, options ...ClientOption) (*Api, error) {
	c := Api{}
	c.httpClient = newClient(apiToken)

	for _, o := range options {
		err := o(&c)
		if err != nil {
			return nil, err
		}
	}
	bs := BaseService{c.httpClient, PageOptions{}, nil}

	c.Projects = func() *ProjectService { return &ProjectService{BaseService: bs} }
	c.Teams = func() *TeamService { return &TeamService{bs} }
	c.TeamUsers = func() *TeamUserService { return &TeamUserService{bs} }
	c.TeamUserGroups = func() *TeamUserGroupService { return &TeamUserGroupService{bs} }

	c.Contributors = func() *ContributorService { return &ContributorService{bs} }
	c.Comments = func() *CommentService { return &CommentService{bs} }
	c.Keys = func() *KeyService { return &KeyService{BaseService: bs} }
	c.Tasks = func() *TaskService { return &TaskService{BaseService: bs} }

	c.Screenshots = func() *ScreenshotService { return &ScreenshotService{BaseService: bs} }
	c.Snapshots = func() *SnapshotService { return &SnapshotService{bs} }
	c.Languages = func() *LanguageService { return &LanguageService{bs} }
	c.Translations = func() *TranslationService { return &TranslationService{BaseService: bs} }

	c.TranslationProviders = func() *TranslationProviderService { return &TranslationProviderService{bs} }
	c.TranslationStatuses = func() *TranslationStatusService { return &TranslationStatusService{bs} }
	c.Orders = func() *OrderService { return &OrderService{bs} }
	c.PaymentCards = func() *PaymentCardService { return &PaymentCardService{bs} }

	c.Webhooks = func() *WebhookService { return &WebhookService{bs} }
	c.Files = func() *FileService { return &FileService{BaseService: bs} }

	return &c, nil
}

// WithRetryCount returns a client ClientOption setting the retry count of outgoing requests.
// if count is zero retries are disabled.
func WithRetryCount(count int) ClientOption {
	return func(c *Api) error {
		if count < 0 {
			return errors.New("lokalise: retry count must be positive")
		}
		c.httpClient.Client.SetRetryCount(count)
		return nil
	}
}

// WithBaseURL returns a ClientOption setting the base URL of the client.
//
// This should only be used for testing different API versions or for using a mocked
// backend in tests.
//noinspection GoUnusedExportedFunction
func WithBaseURL(url string) ClientOption {
	return func(c *Api) error {
		c.httpClient.Client.SetHostURL(url)
		return nil
	}
}
