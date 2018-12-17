package model

type Screenshot struct {
	ScreenshotID   int64    `json:"screenshot_id,omitempty"`
	Title          string   `json:"title,omitempty"`
	Description    string   `json:"description,omitempty"`
	ScreenshotTags []string `json:"screenshot_tags,omitempty"`
	URL            string   `json:"url,omitempty"`
}
