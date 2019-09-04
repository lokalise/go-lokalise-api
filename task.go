package lokalise

import (
	"context"
	"fmt"

	"github.com/go-resty/resty/v2"
)

type TaskStatus string

const (
	StatusCompleted  TaskStatus = "completed"
	StatusInProgress TaskStatus = "in progress"
	StatusCreated    TaskStatus = "created"
	StatusQueued     TaskStatus = "queued"
)

type TaskType string

const (
	TaskTypeTranslation TaskType = "translation"
	TaskTypeReview      TaskType = "review"
)

type Task struct {
	TaskID             int64      `json:"task_id,omitempty"`
	Title              string     `json:"title,omitempty"`
	Description        string     `json:"description,omitempty"`
	Status             TaskStatus `json:"status,omitempty"`
	Progress           int        `json:"progress,omitempty"`
	DueDate            string     `json:"due_date,omitempty"`
	KeysCount          int64      `json:"keys_count,omitempty"`
	WordsCount         int64      `json:"words_count,omitempty"`
	CreatedAt          string     `json:"created_at,omitempty"`
	CreatedBy          int64      `json:"created_by,omitempty"`
	CreatedByEmail     string     `json:"created_by_email,omitempty"`
	CanBeParent        bool       `json:"can_be_parent"`
	TaskType           TaskType   `json:"task_type,omitempty"`
	ParentTaskID       int64      `json:"parent_task_id,omitempty"`
	ClosingTags        []string   `json:"closing_tags,omitempty"`
	LockTranslations   bool       `json:"do_lock_translations"`
	Languages          []Language `json:"languages,omitempty"`
	AutoCloseLanguages bool       `json:"auto_close_languages,omitempty"`
	AutoCloseTask      bool       `json:"auto_close_task,omitempty"`
	CompletedAt        string     `json:"completed_at,omitempty"`
	CompletedBy        int64      `json:"completed_by,omitempty"`
	CompletedByEmail   string     `json:"completed_by_email,omitempty"`
}

type LanguageStatus string

const (
	LanguageStatusCompleted  LanguageStatus = "completed"
	LanguageStatusInProgress LanguageStatus = "in progress"
	LanguageStatusCreated    LanguageStatus = "created"
)

type Language struct {
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
	UserID   int64  `json:"user_id,omitempty"`
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
	ProjectID string `json:"project_id,omitempty"`
	Tasks     []Task `json:"tasks,omitempty"`
}

type TaskResponse struct {
	ProjectID string `json:"project_id,omitempty"`
	Task      Task   `json:"task,omitempty"`
}

type TaskDeleteResponse struct {
	ProjectID string `json:"project_id,omitempty"`
	Deleted   bool   `json:"task_deleted,omitempty"`
}

type TasksService struct {
	client *Client
}

const (
	pathTasks = "tasks"
)

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

func (c *TasksService) List(ctx context.Context, projectID string, pageOptions TasksOptions) (TasksResponse, error) {
	var res TasksResponse
	resp, err := c.client.getList(ctx, fmt.Sprintf("%s/%s/%s", pathProjects, projectID, pathTasks), &res, pageOptions)
	if err != nil {
		return TasksResponse{}, err
	}
	applyPaged(resp, &res.Paged)
	return res, apiError(resp)
}

func (c *TasksService) Create(ctx context.Context, projectID string, task CreateTaskRequest) (TaskResponse, error) {
	var res TaskResponse
	resp, err := c.client.post(ctx, fmt.Sprintf("%s/%s/%s", pathProjects, projectID, pathTasks), &res, task)
	if err != nil {
		return TaskResponse{}, err
	}
	return res, apiError(resp)
}

func (c *TasksService) Retrieve(ctx context.Context, projectID string, taskID int64) (TaskResponse, error) {
	var res TaskResponse
	resp, err := c.client.get(ctx, fmt.Sprintf("%s/%s/%s/%d", pathProjects, projectID, pathTasks, taskID), &res)
	if err != nil {
		return TaskResponse{}, err
	}
	return res, apiError(resp)
}

func (c *TasksService) Update(ctx context.Context, projectID string, taskID int64, task UpdateTaskRequest) (TaskResponse, error) {
	var res TaskResponse
	resp, err := c.client.put(ctx, fmt.Sprintf("%s/%s/%s/%d", pathProjects, projectID, pathTasks, taskID), &res, task)
	if err != nil {
		return TaskResponse{}, err
	}
	return res, apiError(resp)
}

func (c *TasksService) Delete(ctx context.Context, projectID string, taskID int64) (TaskDeleteResponse, error) {
	var res TaskDeleteResponse
	resp, err := c.client.delete(ctx, fmt.Sprintf("%s/%s/%s/%d", pathProjects, projectID, pathTasks, taskID), &res)
	if err != nil {
		return TaskDeleteResponse{}, err
	}
	return res, apiError(resp)
}
