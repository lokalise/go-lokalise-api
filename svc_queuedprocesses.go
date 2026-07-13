package lokalise

import (
	"fmt"
)

const (
	pathQueuedProcesses = "processes"
)

type QueuedProcessService struct {
	BaseService
}

// ‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// Service entity objects
// _____________________________________________________________________________________________________________________

type ProcessDetails struct {
	ItemsToProcess *int                 `json:"items_to_process,omitempty"`
	ItemsProcessed *int                 `json:"items_processed,omitempty"`
	DownloadUrl    string               `json:"download_url,omitempty"`
	Progress       string               `json:"progress,omitempty"`
	Stage          string               `json:"stage,omitempty"`
	Files          []ProcessDetailsFile `json:"files,omitempty"`
}

type ProcessDetailsFile struct {
	Status           string `json:"status"`
	Message          string `json:"message"`
	NameOriginal     string `json:"name_original"`
	NameCustom       string `json:"name_custom"`
	WordCountTotal   int    `json:"word_count_total"`
	KeyCountTotal    int    `json:"key_count_total"`
	KeyCountInserted int    `json:"key_count_inserted"`
	KeyCountUpdated  int    `json:"key_count_updated"`
	KeyCountSkipped  int    `json:"key_count_skipped"`
}

type QueuedProcess struct {
	ID      string         `json:"process_id"`
	Type    string         `json:"type"`
	Status  string         `json:"status"`
	Message string         `json:"message"`
	Details ProcessDetails `json:"details"`
	WithCreationUser
	WithCreationTime
}

// ‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// Service request/response objects
// _____________________________________________________________________________________________________________________

type QueuedProcessesResponse struct {
	WithProjectID
	Processes []QueuedProcess `json:"processes"`
}

type QueuedProcessResponse struct {
	WithProjectID
	Process QueuedProcess `json:"process"`
}

// ‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// Service methods
// _____________________________________________________________________________________________________________________

func (c *QueuedProcessService) List(projectID string) (r QueuedProcessesResponse, err error) {
	resp, err := c.get(c.Ctx(), fmt.Sprintf("%s/%s/%s", pathProjects, projectID, pathQueuedProcesses), &r)

	if err != nil {
		return
	}
	return r, apiError(resp)
}

func (c *QueuedProcessService) Retrieve(projectID string, processID string) (r QueuedProcessResponse, err error) {
	resp, err := c.get(c.Ctx(), pathQueuedProcessById(projectID, processID), &r)

	if err != nil {
		return
	}
	return r, apiError(resp)
}

func pathQueuedProcessById(projectID string, processID string) string {
	return fmt.Sprintf("%s/%s/%s/%s", pathProjects, projectID, pathQueuedProcesses, processID)
}
