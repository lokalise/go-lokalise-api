package model

type Screenshot struct {
	ScreenshotID   int64    `json:"screenshot_id,omitempty"`
	KeyIDs         []int64  `json:"key_ids,omitempty"`
	URL            string   `json:"url,omitempty"`
	Title          string   `json:"title,omitempty"`
	Description    string   `json:"description,omitempty"`
	ScreenshotTags []string `json:"screenshot_tags,omitempty"`
	Width          int64    `json:"width,omitempty"`
	Height         int64    `json:"height,omitempty"`
}
