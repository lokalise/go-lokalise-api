package lokalise

import (
	"encoding/json"
	"errors"
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

type QueuedProcessType string

const (
	FileImport   QueuedProcessType = "file-import"
	SketchImport QueuedProcessType = "sketch-import"
)

func (qpt *QueuedProcessType) UnmarshalJSON(b []byte) error {
	// Define a secondary type to avoid ending up with a recursive call to json.Unmarshal
	type QPT QueuedProcessType
	var r = (*QPT)(qpt)
	err := json.Unmarshal(b, &r)
	if err != nil {
		panic(err)
	}
	switch *qpt {
	case FileImport, SketchImport:
		return nil
	}
	return errors.New("invalid QueuedProcess type")
}

type QueuedProcessStatus string

const (
	Queued    QueuedProcessStatus = "queued"
	Running   QueuedProcessStatus = "running"
	Cancelled QueuedProcessStatus = "cancelled"
	Finished  QueuedProcessStatus = "finished"
	Failed    QueuedProcessStatus = "failed"
)

func (qps *QueuedProcessStatus) UnmarshalJSON(b []byte) error {
	// Define a secondary type to avoid ending up with a recursive call to json.Unmarshal
	type QPS QueuedProcessStatus
	var r = (*QPS)(qps)
	err := json.Unmarshal(b, &r)
	if err != nil {
		panic(err)
	}
	switch *qps {
	case Queued, Running, Cancelled, Finished, Failed:
		return nil
	}
	return errors.New("invalid QueuedProcess status")
}

type QueuedProcess struct {
	ID      string              `json:"process_id"`
	Type    QueuedProcessType   `json:"type"`
	Status  QueuedProcessStatus `json:"status"`
	Message string              `json:"message"`
	Url     string              `json:"url"`
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

func (c *QueuedProcessService) RetrieveDetailed(projectID string, processID string, processType QueuedProcessType) (r interface{}, err error) {
	resp, err := c.get(c.Ctx(), pathQueuedProcessByIdAndType(projectID, processType, processID), &r)
	if err != nil {
		return
	}

	return r, apiError(resp)
}

func pathQueuedProcessByIdAndType(projectID string, processType QueuedProcessType, processID string) string {
	return fmt.Sprintf("%s/%s/%s/%s/%s", pathProjects, projectID, pathQueuedProcesses, processType, processID)
}

func pathQueuedProcessById(projectID string, processID string) string {
	return fmt.Sprintf("%s/%s/%s/%s", pathProjects, projectID, pathQueuedProcesses, processID)
}
