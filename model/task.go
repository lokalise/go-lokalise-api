package model

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
