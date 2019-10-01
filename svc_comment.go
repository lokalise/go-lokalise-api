package lokalise

import (
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
	WithProjectID
	Comments []Comment `json:"comments"`
}

type CommentResponse struct {
	WithProjectID
	Comment Comment `json:"comment"`
}

type DeleteCommentResponse struct {
	WithProjectID
	IsDeleted bool `json:"comment_deleted"`
}

type CommentService struct {
	BaseService
}

func (c *CommentService) ListProject(projectID string, pageOptions PageOptions) (r ListCommentsResponse, err error) {
	resp, err := c.getList(c.Ctx(), pathComments(projectID), &r, &pageOptions)

	if err != nil {
		return
	}
	applyPaged(resp, &r.Paged)
	return r, apiError(resp)
}

func (c *CommentService) ListByKey(projectID string, keyID int64, pageOptions PageOptions) (r ListCommentsResponse, err error) {
	resp, err := c.getList(c.Ctx(), pathCommentsByKey(projectID, keyID), &r, &pageOptions)

	if err != nil {
		return
	}
	applyPaged(resp, &r.Paged)
	return r, apiError(resp)
}

func (c *CommentService) Create(projectID string, keyID int64, comments []NewComment) (r ListCommentsResponse, err error) {
	resp, err := c.post(c.Ctx(), pathCommentsByKey(projectID, keyID), &r, map[string]interface{}{"comments": comments})

	if err != nil {
		return
	}
	applyPaged(resp, &r.Paged)
	return r, apiError(resp)
}

func (c *CommentService) Retrieve(projectID string, keyID, commentID int64) (r CommentResponse, err error) {
	url := path.Join(pathCommentsByKey(projectID, keyID), strconv.FormatInt(commentID, 10))
	resp, err := c.get(c.Ctx(), url, &r)

	if err != nil {
		return
	}
	return r, apiError(resp)
}

func (c *CommentService) Delete(projectID string, keyID, commentID int64) (r DeleteCommentResponse, err error) {
	url := path.Join(pathCommentsByKey(projectID, keyID), strconv.FormatInt(commentID, 10))
	resp, err := c.delete(c.Ctx(), url, &r)

	if err != nil {
		return
	}
	return r, apiError(resp)
}
