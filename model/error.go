package model

import "fmt"

// Error is an API error.
type Error struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

func (r Error) Error() string {
	return fmt.Sprintf("API request error %d %s", r.Code, r.Message)
}
