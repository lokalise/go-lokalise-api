package lokalise

import (
	"net/http"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/google/go-querystring/query"
)

type PageCounter interface {
	NumberOfPages() int64
	CurrentPage() int64
}

type Paged struct {
	TotalCount int64 `json:"-"`
	PageCount  int64 `json:"-"`
	Limit      int64 `json:"-"`
	Page       int64 `json:"-"`
}

func (p Paged) NumberOfPages() int64 {
	return p.PageCount
}

func (p Paged) CurrentPage() int64 {
	return p.Page
}

type OptionsApplier interface {
	Apply(req *resty.Request)
}

type PageOptions struct {
	Limit uint `url:"limit,omitempty"`
	Page  uint `url:"page,omitempty"`
}

func (options PageOptions) Apply(req *resty.Request) {
	v, _ := query.Values(options)
	req.SetQueryString(v.Encode())
}

const (
	headerTotalCount = "X-Pagination-Total-Count"
	headerPageCount  = "X-Pagination-Page-Count"
	headerLimit      = "X-Pagination-Limit"
	headerPage       = "X-Pagination-Page"
)

func applyPaged(res *resty.Response, paged *Paged) {
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
