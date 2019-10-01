// Package lokalise provides functions to access the Lokalise web API.
package lokalise

import (
	"errors"
)

type Api struct {
	httpClient *restClient

	Projects             func() *ProjectsService
	Teams                func() *TeamsService
	TeamUsers            func() *TeamUsersService
	TeamUserGroups       func() *TeamUserGroupsService
	Contributors         func() *ContributorsService
	Comments             func() *CommentService
	Keys                 func() *KeysService
	Tasks                func() *TasksService
	Screenshots          func() *ScreenshotsService
	Snapshots            func() *SnapshotsService
	Languages            func() *LanguagesService
	Translations         func() *TranslationsService
	TranslationProviders func() *TranslationProviderService
	TranslationStatuses  func() *TranslationStatusesService
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

	c.Projects = func() *ProjectsService { return &ProjectsService{bs} }
	c.Teams = func() *TeamsService { return &TeamsService{bs} }
	c.TeamUsers = func() *TeamUsersService { return &TeamUsersService{bs} }
	c.TeamUserGroups = func() *TeamUserGroupsService { return &TeamUserGroupsService{bs} }

	c.Contributors = func() *ContributorsService { return &ContributorsService{bs} }
	c.Comments = func() *CommentService { return &CommentService{bs} }
	c.Keys = func() *KeysService { return &KeysService{bs} }
	c.Tasks = func() *TasksService { return &TasksService{bs} }

	c.Screenshots = func() *ScreenshotsService { return &ScreenshotsService{bs} }
	c.Snapshots = func() *SnapshotsService { return &SnapshotsService{bs} }
	c.Languages = func() *LanguagesService { return &LanguagesService{bs} }
	c.Translations = func() *TranslationsService { return &TranslationsService{bs} }

	c.TranslationProviders = func() *TranslationProviderService { return &TranslationProviderService{bs} }
	c.TranslationStatuses = func() *TranslationStatusesService { return &TranslationStatusesService{bs} }
	c.Orders = func() *OrderService { return &OrderService{bs} }
	c.PaymentCards = func() *PaymentCardService { return &PaymentCardService{bs} }

	c.Webhooks = func() *WebhookService { return &WebhookService{bs} }
	c.Files = func() *FileService { return &FileService{bs} }

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
