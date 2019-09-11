package lokalise

import (
	"context"
	"fmt"
	"path"
	"strconv"
)

type Comment struct {
	CommentID    int64  `json:"comment_id,omitempty"`
	KeyID        int64  `json:"key_id,omitempty"`
	Comment      string `json:"comment,omitempty"`
	AddedBy      int64  `json:"added_by,omitempty"`
	AddedByEmail string `json:"added_by_email,omitempty"`
	AddedAt      string `json:"added_at,omitempty"`
	AddedAtTs    int64  `json:"added_at_timestamp,omitempty"`
}

func pathComments(projectID string) string {
	return fmt.Sprintf("%s/%s/comments", pathProjects, projectID)
}

func pathCommentsByKey(projectID string, keyID int64) string {
	return fmt.Sprintf("%s/%s/%s/%d/comments", pathProjects, projectID, pathKeys, keyID)
}

type NewComment struct {
	Comment string `json:"comment"`
}

type ListCommentsResponse struct {
	Paged
	ProjectID string    `json:"project_id"`
	Comments  []Comment `json:"comments"`
}

type CommentResponse struct {
	ProjectID string  `json:"project_id"`
	Comment   Comment `json:"comment"`
}

type DeleteCommentResponse struct {
	ProjectID string `json:"project_id"`
	IsDeleted bool   `json:"comment_deleted"`
}

type CommentService struct {
	client *Client
}

func (c *CommentService) ListProject(ctx context.Context, projectID string, pageOptions PageOptions) (result ListCommentsResponse, err error) {
	var res ListCommentsResponse
	resp, err := c.client.getList(ctx, pathComments(projectID), &res, &pageOptions)
	if err != nil {
		return result, err
	}

	applyPaged(resp, &result.Paged)
	return result, apiError(resp)
}

func (c *CommentService) ListByKey(ctx context.Context, projectID string, keyID int64, pageOptions PageOptions) (result ListCommentsResponse, err error) {
	resp, err := c.client.getList(ctx, pathCommentsByKey(projectID, keyID), &result, &pageOptions)
	if err != nil {
		return result, err
	}

	applyPaged(resp, &result.Paged)
	return result, apiError(resp)
}

func (c *CommentService) Create(ctx context.Context, projectID string, keyID int64, comments []NewComment) (result ListCommentsResponse, err error) {
	resp, err := c.client.post(ctx, pathCommentsByKey(projectID, keyID), &result, map[string]interface{}{"comments": comments})
	if err != nil {
		return result, err
	}

	applyPaged(resp, &result.Paged)
	return result, apiError(resp)
}

func (c *CommentService) Retrieve(ctx context.Context, projectID string, keyID, commentID int64) (result CommentResponse, err error) {
	url := path.Join(pathCommentsByKey(projectID, keyID), strconv.FormatInt(commentID, 10))

	resp, err := c.client.get(ctx, url, &result)
	if err != nil {
		return result, err
	}
	return result, apiError(resp)
}

func (c *CommentService) Delete(ctx context.Context, projectID string, keyID, commentID int64) (result DeleteCommentResponse, err error) {
	url := path.Join(pathCommentsByKey(projectID, keyID), strconv.FormatInt(commentID, 10))

	resp, err := c.client.delete(ctx, url, &result)
	if err != nil {
		return result, err
	}
	return result, apiError(resp)
}
