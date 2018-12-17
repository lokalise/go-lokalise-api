package model

import "fmt"

// ErrorKey is error for key api
type ErrorKey struct {
	KeyName string `json:"key_name,omitempty"`
}

// Error is an API error.
type Error struct {
	Code    int      `json:"code,omitempty"`
	Message string   `json:"message,omitempty"`
	Key     ErrorKey `json:"key,omitempty"`
}

func (r Error) Error() string {
	return fmt.Sprintf("API request error %d %s", r.Code, r.Message)
}
