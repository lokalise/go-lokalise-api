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

type QueuedProcess struct {
	ID      string      `json:"process_id"`
	Type    string      `json:"type"`
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Details interface{} `json:"details"`
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
