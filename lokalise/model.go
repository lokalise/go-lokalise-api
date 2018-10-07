package lokalise

// RequestError is the API error model.
type RequestError struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}
