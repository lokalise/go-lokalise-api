package model

type Comment struct {
	CommentID    int64  `json:"comment_id,omitempty"`
	KeyID        int64  `json:"key_id,omitempty"`
	Comment      string `json:"comment,omitempty"`
	AddedBy      int64  `json:"added_by,omitempty"`
	AddedByEmail string `json:"added_by_email,omitempty"`
	AddedAt      string `json:"added_at,omitempty"`
}
