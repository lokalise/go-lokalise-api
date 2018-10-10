package lokalise

import (
	"errors"
	"fmt"

	"github.com/go-resty/resty"
)

// RequestError is the API error model.
type RequestError struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

func (r RequestError) Error() string {
	return fmt.Sprintf("API request error %d %s", r.Code, r.Message)
}

// apiError indentifies whether the response contains an API error.
func apiError(res *resty.Response) error {
	if !res.IsError() {
		return nil
	}
	responseError := res.Error()
	if responseError == nil {
		return errors.New("lokalise: response marked as error but no data returned")
	}
	responseErrorModel, ok := responseError.(*RequestError)
	if !ok {
		return errors.New("lokalise: response error model unknown")
	}
	return responseErrorModel
}
