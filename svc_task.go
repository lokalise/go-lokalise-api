package lokalise

import (
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/google/go-querystring/query"
)

type TaskService struct {
	BaseService

	listOpts TaskListOptions
}

// ‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// Service entity objects
// _____________________________________________________________________________________________________________________

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

type (
	TaskStatus     string
	TaskType       string
	LanguageStatus string
)

type Task struct {
	WithCreationTime
	WithCreationUser

	TaskID                     int64          `json:"task_id"`
	Title                      string         `json:"title"`
	Description                string         `json:"description"`
	Status                     TaskStatus     `json:"status"`
	Progress                   int            `json:"progress"`
	DueDate                    string         `json:"due_date"`
	KeysCount                  int64          `json:"keys_count"`
	WordsCount                 int64          `json:"words_count"`
	CanBeParent                bool           `json:"can_be_parent"`
	TaskType                   TaskType       `json:"task_type"`
	ParentTaskID               int64          `json:"parent_task_id"`
	ClosingTags                []string       `json:"closing_tags"`
	LockTranslations           bool           `json:"do_lock_translations"`
	Languages                  []TaskLanguage `json:"languages"`
	AutoCloseLanguages         bool           `json:"auto_close_languages"`
	AutoCloseTask              bool           `json:"auto_close_task"`
	CompletedAt                string         `json:"completed_at"`
	CompletedAtTs              int64          `json:"completed_at_timestamp"`
	CompletedBy                int64          `json:"completed_by"`
	CompletedByEmail           string         `json:"completed_by_email"`
	CustomTranslationStatusIDs []int64        `json:"custom_translation_status_ids,omitempty"`
}

type TaskLanguage struct {
	LanguageISO       string           `json:"language_iso"`
	Users             []TaskUser       `json:"users"`
	Groups            []TaskGroup      `json:"groups"`
	Keys              []int64          `json:"keys"`
	Status            LanguageStatus   `json:"status"`
	Progress          int              `json:"progress"`
	InitialTMLeverage map[string]int64 `json:"initial_tm_leverage"`
	KeysCount         int64            `json:"keys_count"`
	WordsCount        int64            `json:"words_count"`
	CompletedAt       string           `json:"completed_at"`
	CompletedAtTs     int64            `json:"completed_at_timestamp"`
	CompletedBy       int64            `json:"completed_by"`
	CompletedByEmail  string           `json:"completed_by_email"`
}

type TaskUser struct {
	WithUserID
	Email    string `json:"email"`
	Fullname string `json:"fullname"`
}

type TaskGroup struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// ‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// Service request/response objects
// _____________________________________________________________________________________________________________________

type CreateTask struct {
	Title     string           `json:"title"`
	Languages []CreateTaskLang `json:"languages"`

	Description        string   `json:"description,omitempty"`
	DueDate            string   `json:"due_date,omitempty"`
	Keys               []int64  `json:"keys,omitempty"`
	AutoCloseLanguages *bool    `json:"auto_close_languages,omitempty"`
	AutoCloseTask      *bool    `json:"auto_close_task,omitempty"`
	InitialTMLeverage  bool     `json:"initial_tm_leverage,omitempty"`
	TaskType           TaskType `json:"task_type,omitempty"`
	ParentTaskID       int64    `json:"parent_task_id,omitempty"`
	ClosingTags        []string `json:"closing_tags,omitempty"`
	LockTranslations   bool     `json:"do_lock_translations,omitempty"`

	CustomTranslationStatusIDs []int64 `json:"custom_translation_status_ids,omitempty"`
}

type CreateTaskLang struct {
	LanguageISO   string  `json:"language_iso,omitempty"`
	Users         []int64 `json:"users,omitempty"`
	Groups        []int64 `json:"groups,omitempty"`
	CloseLanguage bool    `json:"close_language,omitempty"` // only for updating
}

type UpdateTask struct {
	Title              string           `json:"title,omitempty"`
	Description        string           `json:"description,omitempty"`
	DueDate            string           `json:"due_date,omitempty"`
	Languages          []CreateTaskLang `json:"languages,omitempty"`
	AutoCloseLanguages *bool            `json:"auto_close_languages,omitempty"`
	AutoCloseTask      *bool            `json:"auto_close_task,omitempty"`
	CloseTask          bool             `json:"close_task,omitempty"`
	ClosingTags        []string         `json:"closing_tags,omitempty"`
	LockTranslations   bool             `json:"do_lock_translations,omitempty"`
}

type TasksResponse struct {
	Paged
	WithProjectID
	Tasks []Task `json:"tasks"`
}

type TaskResponse struct {
	WithProjectID
	Task Task `json:"task"`
}

type DeleteTaskResponse struct {
	WithProjectID
	Deleted bool `json:"task_deleted"`
}

// ‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// Service methods
// _____________________________________________________________________________________________________________________

func (c *TaskService) List(projectID string) (r TasksResponse, err error) {
	resp, err := c.getList(c.Ctx(), fmt.Sprintf("%s/%s/%s", pathProjects, projectID, pathTasks), &r, c.ListOpts())

	if err != nil {
		return
	}
	applyPaged(resp, &r.Paged)
	return r, apiError(resp)
}

func (c *TaskService) Create(projectID string, task CreateTask) (r TaskResponse, err error) {
	resp, err := c.post(c.Ctx(), fmt.Sprintf("%s/%s/%s", pathProjects, projectID, pathTasks), &r, task)

	if err != nil {
		return
	}
	return r, apiError(resp)
}

func (c *TaskService) Retrieve(projectID string, taskID int64) (r TaskResponse, err error) {
	resp, err := c.get(c.Ctx(), fmt.Sprintf("%s/%s/%s/%d", pathProjects, projectID, pathTasks, taskID), &r)

	if err != nil {
		return
	}
	return r, apiError(resp)
}

func (c *TaskService) Update(projectID string, taskID int64, task UpdateTask) (r TaskResponse, err error) {
	resp, err := c.put(c.Ctx(), fmt.Sprintf("%s/%s/%s/%d", pathProjects, projectID, pathTasks, taskID), &r, task)

	if err != nil {
		return
	}
	return r, apiError(resp)
}

func (c *TaskService) Delete(projectID string, taskID int64) (r DeleteTaskResponse, err error) {
	resp, err := c.delete(c.Ctx(), fmt.Sprintf("%s/%s/%s/%d", pathProjects, projectID, pathTasks, taskID), &r)

	if err != nil {
		return
	}
	return r, apiError(resp)
}

// ‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// Additional methods
// _____________________________________________________________________________________________________________________

type TaskListOptions struct {
	Limit          uint   `url:"limit,omitempty"`
	Page           uint   `url:"page,omitempty"`
	FilterStatuses string `url:"filter_statuses,omitempty"`
	Title          string `url:"filter_title,omitempty"`
}

func (options TaskListOptions) Apply(req *resty.Request) {
	v, _ := query.Values(options)
	req.SetQueryString(v.Encode())
}

func (c *TaskService) ListOpts() TaskListOptions        { return c.listOpts }
func (c *TaskService) SetListOptions(o TaskListOptions) { c.listOpts = o }
func (c *TaskService) WithListOptions(o TaskListOptions) *TaskService {
	c.listOpts = o
	return c
}
