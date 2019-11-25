package lokalise

import (
	"errors"
	"time"
)

type Api struct {
	httpClient  *restClient
	pageOptions PageOptions

	Branches             func() *BranchService
	Comments             func() *CommentService
	Contributors         func() *ContributorService
	Files                func() *FileService
	Keys                 func() *KeyService
	Languages            func() *LanguageService
	Orders               func() *OrderService
	PaymentCards         func() *PaymentCardService
	Projects             func() *ProjectService
	Screenshots          func() *ScreenshotService
	Snapshots            func() *SnapshotService
	Tasks                func() *TaskService
	Teams                func() *TeamService
	TeamUserGroups       func() *TeamUserGroupService
	TeamUsers            func() *TeamUserService
	TranslationProviders func() *TranslationProviderService
	Translations         func() *TranslationService
	TranslationStatuses  func() *TranslationStatusService
	Webhooks             func() *WebhookService
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
	bs := BaseService{c.httpClient, c.pageOptions, nil}

	// predefined list options if any
	prjOpts := ProjectListOptions{Page: c.pageOptions.Page, Limit: c.pageOptions.Limit}
	keyOpts := KeyListOptions{Page: c.pageOptions.Page, Limit: c.pageOptions.Limit}
	taskOpts := TaskListOptions{Page: c.pageOptions.Page, Limit: c.pageOptions.Limit}
	scOpts := ScreenshotListOptions{Page: c.pageOptions.Page, Limit: c.pageOptions.Limit}
	trOpts := TranslationListOptions{Page: c.pageOptions.Page, Limit: c.pageOptions.Limit}
	fOpts := FileListOptions{Page: c.pageOptions.Page, Limit: c.pageOptions.Limit}

	c.Projects = func() *ProjectService { return &ProjectService{BaseService: bs, opts: prjOpts} }
	c.Branches = func() *BranchService { return &BranchService{bs} }
	c.Teams = func() *TeamService { return &TeamService{bs} }
	c.TeamUsers = func() *TeamUserService { return &TeamUserService{bs} }
	c.TeamUserGroups = func() *TeamUserGroupService { return &TeamUserGroupService{bs} }

	c.Contributors = func() *ContributorService { return &ContributorService{bs} }
	c.Comments = func() *CommentService { return &CommentService{bs} }
	c.Keys = func() *KeyService { return &KeyService{BaseService: bs, listOpts: keyOpts} }
	c.Tasks = func() *TaskService { return &TaskService{BaseService: bs, listOpts: taskOpts} }

	c.Screenshots = func() *ScreenshotService { return &ScreenshotService{BaseService: bs, listOpts: scOpts} }
	c.Snapshots = func() *SnapshotService { return &SnapshotService{bs} }
	c.Languages = func() *LanguageService { return &LanguageService{bs} }
	c.Translations = func() *TranslationService { return &TranslationService{BaseService: bs, opts: trOpts} }

	c.TranslationProviders = func() *TranslationProviderService { return &TranslationProviderService{bs} }
	c.TranslationStatuses = func() *TranslationStatusService { return &TranslationStatusService{bs} }
	c.Orders = func() *OrderService { return &OrderService{bs} }
	c.PaymentCards = func() *PaymentCardService { return &PaymentCardService{bs} }

	c.Webhooks = func() *WebhookService { return &WebhookService{bs} }
	c.Files = func() *FileService { return &FileService{BaseService: bs, opts: fOpts} }

	return &c, nil
}

// WithBaseURL returns a ClientOption setting the base URL of the client.
// This should only be used for testing different API versions or for using a mocked
// backend in tests.
//noinspection GoUnusedExportedFunction
func WithBaseURL(url string) ClientOption {
	return func(c *Api) error {
		c.httpClient.Client.SetHostURL(url)
		return nil
	}
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

// Sets default wait time to sleep before retrying request.
// Default is 100 milliseconds.
//noinspection GoUnusedExportedFunction
func WithRetryTimeout(t time.Duration) ClientOption {
	return func(c *Api) error {
		c.httpClient.Client.SetRetryWaitTime(t)
		return nil
	}
}

func WithConnectionTimeout(t time.Duration) ClientOption {
	return func(c *Api) error {
		c.httpClient.Client.SetTimeout(t)
		return nil
	}
}

func WithDebug(dbg bool) ClientOption {
	return func(c *Api) error {
		c.httpClient.Client.SetDebug(dbg)
		return nil
	}
}

func WithPageLimit(limit uint) ClientOption {
	return func(c *Api) error {
		c.pageOptions.Limit = limit
		return nil
	}
}

// Bool is a helper routine that allocates a new bool value
// to store v and returns a pointer to it.
func Bool(v bool) *bool {
	return &v
}
