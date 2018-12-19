package lokalise

import (
	"errors"

	"github.com/17media/go-lokalise-api/model"
	"github.com/go-resty/resty"
)

type errorResponse struct {
	Error model.Error `json:"error"`
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
