package lokalise

import (
	"context"
	"fmt"

	"github.com/go-resty/resty"
	"github.com/lokalise/go-lokalise-api/model"
)

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

func (c *TasksService) List(ctx context.Context, projectID string, pageOptions TasksOptions) (model.TasksResponse, error) {
	var res model.TasksResponse
	resp, err := c.client.getList(ctx, fmt.Sprintf("%s/%s/%s", pathProjects, projectID, pathTasks), &res, pageOptions)
	if err != nil {
		return model.TasksResponse{}, err
	}
	applyPaged(resp, &res.Paged)
	return res, apiError(resp)
}

func (c *TasksService) Create(ctx context.Context, projectID string, task model.CreateTaskRequest) (model.TaskResponse, error) {
	var res model.TaskResponse
	resp, err := c.client.post(ctx, fmt.Sprintf("%s/%s/%s", pathProjects, projectID, pathTasks), &res, task)
	if err != nil {
		return model.TaskResponse{}, err
	}
	return res, apiError(resp)
}

func (c *TasksService) Retrieve(ctx context.Context, projectID string, taskID int64) (model.TaskResponse, error) {
	var res model.TaskResponse
	resp, err := c.client.get(ctx, fmt.Sprintf("%s/%s/%s/%d", pathProjects, projectID, pathTasks, taskID), &res)
	if err != nil {
		return model.TaskResponse{}, err
	}
	return res, apiError(resp)
}

func (c *TasksService) Update(ctx context.Context, projectID string, taskID int64, task model.UpdateTaskRequest) (model.TaskResponse, error) {
	var res model.TaskResponse
	resp, err := c.client.put(ctx, fmt.Sprintf("%s/%s/%s/%d", pathProjects, projectID, pathTasks, taskID), &res, task)
	if err != nil {
		return model.TaskResponse{}, err
	}
	return res, apiError(resp)
}

func (c *TasksService) Delete(ctx context.Context, projectID string, taskID int64) (model.TaskDeleteResponse, error) {
	var res model.TaskDeleteResponse
	resp, err := c.client.delete(ctx, fmt.Sprintf("%s/%s/%s/%d", pathProjects, projectID, pathTasks, taskID), &res)
	if err != nil {
		return model.TaskDeleteResponse{}, err
	}
	return res, apiError(resp)
}
