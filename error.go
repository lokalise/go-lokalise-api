package lokalise

import (
	"errors"
	"fmt"

	"github.com/go-resty/resty/v2"
)

// Error is an API error.
type Error struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

func (r Error) Error() string {
	return fmt.Sprintf("API request error %d %s", r.Code, r.Message)
}

type errorResponse struct {
	Error Error `json:"error"`
}

// apiError identifies whether the response contains an API error.
func apiError(res *resty.Response) error {
	if !res.IsError() {
		return nil
	}
	responseError := res.Error()
	if responseError == nil {
		return errors.New("lokalise: response marked as error but no data returned")
	}
	responseErrorModel, ok := responseError.(*errorResponse)
	if !ok {
		return errors.New("lokalise: response error model unknown")
	}
	return responseErrorModel.Error
}
