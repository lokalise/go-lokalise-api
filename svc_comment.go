package lokalise

import (
	"fmt"
)

const (
	pathComments = "comments"
)

// The Comment service
type CommentService struct {
	BaseService
}

// ‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// Service entity objects
// _____________________________________________________________________________________________________________________

type Comment struct {
	CommentID    int64  `json:"comment_id"`
	KeyID        int64  `json:"key_id"`
	Comment      string `json:"comment"`
	AddedBy      int64  `json:"added_by"`
	AddedByEmail string `json:"added_by_email"`
	AddedAt      string `json:"added_at"`
	AddedAtTs    int64  `json:"added_at_timestamp"`
}

// ‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// Service request/response objects
// _____________________________________________________________________________________________________________________

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

// ‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// Service methods and functions
// _____________________________________________________________________________________________________________________

// Retrieves a list of all comments in the project
func (c *CommentService) ListProject(projectID string) (r ListCommentsResponse, err error) {
	url := fmt.Sprintf("%s/%s/%s", pathProjects, projectID, pathComments)
	resp, err := c.getList(c.Ctx(), url, &r, c.PageOpts())

	if err != nil {
		return
	}
	applyPaged(resp, &r.Paged)
	return r, apiError(resp)
}

// Retrieves a list of all comments for a key
func (c *CommentService) ListByKey(projectID string, keyID int64) (r ListCommentsResponse, err error) {
	resp, err := c.getList(c.Ctx(), pathCommentsByKey(projectID, keyID), &r, c.PageOpts())

	if err != nil {
		return
	}
	applyPaged(resp, &r.Paged)
	return r, apiError(resp)
}

// Adds a set of comments to the key
func (c *CommentService) Create(projectID string, keyID int64, comments []NewComment) (r ListCommentsResponse, err error) {
	resp, err := c.post(c.Ctx(), pathCommentsByKey(projectID, keyID), &r, map[string]interface{}{"comments": comments})

	if err != nil {
		return
	}
	applyPaged(resp, &r.Paged)
	return r, apiError(resp)
}

// Retrieves a Comment
func (c *CommentService) Retrieve(projectID string, keyID, commentID int64) (r CommentResponse, err error) {
	resp, err := c.get(c.Ctx(), pathCommentByKeyAndID(projectID, keyID, commentID), &r)

	if err != nil {
		return
	}
	return r, apiError(resp)
}

// Deletes a comment from the project. Authenticated user can only delete own comments
func (c *CommentService) Delete(projectID string, keyID, commentID int64) (r DeleteCommentResponse, err error) {
	resp, err := c.delete(c.Ctx(), pathCommentByKeyAndID(projectID, keyID, commentID), &r)

	if err != nil {
		return
	}
	return r, apiError(resp)
}

func pathCommentsByKey(projectID string, keyID int64) string {
	return fmt.Sprintf("%s/%s/%s/%d/%s", pathProjects, projectID, pathKeys, keyID, pathComments)
}

func pathCommentByKeyAndID(projectID string, keyID, commentID int64) string {
	return fmt.Sprintf("%s/%s/%s/%d/%s/%d", pathProjects, projectID, pathKeys, keyID, pathComments, commentID)
}
