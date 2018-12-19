package lokalise

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/17media/go-lokalise-api/model"
	"github.com/go-resty/resty"
)

type OptionsApplier interface {
	Apply(req *resty.Request)
}

type PageOptions struct {
	Limit int64
	Page  int64
}

func (options *PageOptions) Apply(req *resty.Request) {
	if options.Limit != 0 {
		req.SetQueryParam("limit", fmt.Sprintf("%d", options.Limit))
	}
	if options.Page != 0 {
		req.SetQueryParam("page", fmt.Sprintf("%d", options.Page))
	}
}

const (
	headerTotalCount = "X-Pagination-Total-Count"
	headerPageCount  = "X-Pagination-Page-Count"
	headerLimit      = "X-Pagination-Limit"
	headerPage       = "X-Pagination-Page"
)

func applyPaged(res *resty.Response, paged *model.Paged) {
	headers := res.Header()
	paged.TotalCount = headerInt64(headers, headerTotalCount)
	paged.PageCount = headerInt64(headers, headerPageCount)
	paged.Limit = headerInt64(headers, headerLimit)
	paged.Page = headerInt64(headers, headerPage)
}

func headerInt64(headers http.Header, headerKey string) int64 {
	headerValue := headers.Get(headerKey)
	if headerValue == "" {
		return -1
	}
	value, err := strconv.ParseInt(headers.Get(headerKey), 10, 64)
	if err != nil {
		return -1
	}
	return value
}
