package lokalise

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

//noinspection GoUnusedConst
const (
	pathTasks = "tasks"

	StatusCompleted  TaskStatus = "completed"
	StatusInProgress TaskStatus = "in progress"
	StatusCreated    TaskStatus = "created"
	StatusQueued     TaskStatus = "queued"

	TaskTypeTranslation TaskType = "translation"
	TaskTypeReview      TaskType = "review"

	LanguageStatusCompleted  LanguageStatus = "completed"
	LanguageStatusInProgress LanguageStatus = "in progress"
	LanguageStatusCreated    LanguageStatus = "created"
)

type TaskStatus string
type TaskType string
type LanguageStatus string

type Task struct {
	WithCreationTime
	TaskID             int64          `json:"task_id,omitempty"`
	Title              string         `json:"title,omitempty"`
	Description        string         `json:"description,omitempty"`
	Status             TaskStatus     `json:"status,omitempty"`
	Progress           int            `json:"progress,omitempty"`
	DueDate            string         `json:"due_date,omitempty"`
	KeysCount          int64          `json:"keys_count,omitempty"`
	WordsCount         int64          `json:"words_count,omitempty"`
	CreatedBy          int64          `json:"created_by,omitempty"`
	CreatedByEmail     string         `json:"created_by_email,omitempty"`
	CanBeParent        bool           `json:"can_be_parent"`
	TaskType           TaskType       `json:"task_type,omitempty"`
	ParentTaskID       int64          `json:"parent_task_id,omitempty"`
	ClosingTags        []string       `json:"closing_tags,omitempty"`
	LockTranslations   bool           `json:"do_lock_translations"`
	Languages          []TaskLanguage `json:"languages,omitempty"`
	AutoCloseLanguages bool           `json:"auto_close_languages,omitempty"`
	AutoCloseTask      bool           `json:"auto_close_task,omitempty"`
	CompletedAt        string         `json:"completed_at,omitempty"`
	CompletedBy        int64          `json:"completed_by,omitempty"`
	CompletedByEmail   string         `json:"completed_by_email,omitempty"`
}

type TaskLanguage struct { // todo embed Lang
	LanguageISO                      string           `json:"language_iso,omitempty"`
	Users                            []TaskUser       `json:"users,omitempty"`
	Groups                           []TaskGroup      `json:"groups,omitempty"`
	Keys                             []int64          `json:"keys,omitempty"`
	Status                           LanguageStatus   `json:"status,omitempty"`
	Progress                         int              `json:"progress,omitempty"`
	InitialTranslationMemoryLeverage map[string]int64 `json:"initial_tm_leverage,omitempty"`
	KeysCount                        int64            `json:"keys_count,omitempty"`
	WordsCount                       int64            `json:"words_count,omitempty"`
	CompletedAt                      string           `json:"completed_at,omitempty"`
	CompletedBy                      int64            `json:"completed_by,omitempty"`
	CompletedByEmail                 string           `json:"completed_by_email,omitempty"`
}

type TaskUser struct {
	WithUserID
	Email    string `json:"email,omitempty"`
	Fullname string `json:"fullname,omitempty"`
}

type TaskGroup struct {
	ID   int64  `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type CreateTaskRequest struct {
	Title                            string                      `json:"title,omitempty"`
	Description                      string                      `json:"description,omitempty"`
	DueDate                          string                      `json:"due_date,omitempty"`
	Keys                             []int64                     `json:"keys,omitempty"`
	Languages                        []CreateTaskLanguageRequest `json:"languages,omitempty"`
	AutoCloseLanguages               bool                        `json:"auto_close_languages,omitempty"`
	AutoCloseTask                    bool                        `json:"auto_close_task,omitempty"`
	InitialTranslationMemoryLeverage bool                        `json:"initial_tm_leverage,omitempty"`
	TaskType                         TaskType                    `json:"task_type,omitempty"`
	ParentTaskID                     int64                       `json:"parent_task_id,omitempty"`
	ClosingTags                      []string                    `json:"closing_tags,omitempty"`
	LockTranslations                 bool                        `json:"do_lock_translations,omitempty"`
}

type CreateTaskLanguageRequest struct {
	LanguageISO string  `json:"language_iso,omitempty"`
	Users       []int64 `json:"users,omitempty"`
	Groups      []int64 `json:"groups,omitempty"`
}

type UpdateTaskRequest struct {
	Title              string                      `json:"title,omitempty"`
	Description        string                      `json:"description,omitempty"`
	DueDate            string                      `json:"due_date,omitempty"`
	Languages          []UpdateTaskLanguageRequest `json:"languages,omitempty"`
	AutoCloseLanguages bool                        `json:"auto_close_languages,omitempty"`
	AutoCloseTask      bool                        `json:"auto_close_task,omitempty"`
	CloseTask          bool                        `json:"close_task,omitempty"`
	ClosingTags        []string                    `json:"closing_tags,omitempty"`
	LockTranslations   bool                        `json:"do_lock_translations,omitempty"`
}

type UpdateTaskLanguageRequest struct {
	LanguageISO   string      `json:"language_iso,omitempty"`
	Users         []int64     `json:"users,omitempty"`
	Groups        []TaskGroup `json:"groups,omitempty"`
	CloseLanguage bool        `json:"close_language,omitempty"`
}

type TasksResponse struct {
	Paged
	WithProjectID
	Tasks []Task `json:"tasks,omitempty"`
}

type TaskResponse struct {
	WithProjectID
	Task Task `json:"task,omitempty"`
}

type TaskDeleteResponse struct {
	WithProjectID
	Deleted bool `json:"task_deleted,omitempty"`
}

type TasksService struct {
	BaseService
}

type TasksOptions struct {
	PageOptions
	Title string
}

func (options TasksOptions) Apply(req *resty.Request) {
	options.PageOptions.Apply(req)
	if options.Title != "" {
		req.SetQueryParam("filter_title", options.Title)
	}
}

func (c *TasksService) List(projectID string, pageOptions TasksOptions) (r TasksResponse, err error) {
	resp, err := c.getList(c.Ctx(), fmt.Sprintf("%s/%s/%s", pathProjects, projectID, pathTasks), &r, pageOptions)

	if err != nil {
		return
	}
	applyPaged(resp, &r.Paged)
	return r, apiError(resp)
}

func (c *TasksService) Create(projectID string, task CreateTaskRequest) (r TaskResponse, err error) {
	resp, err := c.post(c.Ctx(), fmt.Sprintf("%s/%s/%s", pathProjects, projectID, pathTasks), &r, task)

	if err != nil {
		return
	}
	return r, apiError(resp)
}

func (c *TasksService) Retrieve(projectID string, taskID int64) (r TaskResponse, err error) {
	resp, err := c.get(c.Ctx(), fmt.Sprintf("%s/%s/%s/%d", pathProjects, projectID, pathTasks, taskID), &r)

	if err != nil {
		return
	}
	return r, apiError(resp)
}

func (c *TasksService) Update(projectID string, taskID int64, task UpdateTaskRequest) (r TaskResponse, err error) {
	resp, err := c.put(c.Ctx(), fmt.Sprintf("%s/%s/%s/%d", pathProjects, projectID, pathTasks, taskID), &r, task)

	if err != nil {
		return
	}
	return r, apiError(resp)
}

func (c *TasksService) Delete(projectID string, taskID int64) (r TaskDeleteResponse, err error) {
	resp, err := c.delete(c.Ctx(), fmt.Sprintf("%s/%s/%s/%d", pathProjects, projectID, pathTasks, taskID), &r)

	if err != nil {
		return
	}
	return r, apiError(resp)
}
