package lokalise

import (
	"net/http"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/google/go-querystring/query"
)

const (
	PaginationOffset = "offset"
	PaginationCursor = "cursor"
)

type PageCounter interface {
	NumberOfPages() int64
	CurrentPage() int64
}

type CursorPager interface {
	HasNextCursor() bool
	NextCursor() string
}

type Paged struct {
	TotalCount int64  `json:"-"`
	PageCount  int64  `json:"-"`
	Limit      int64  `json:"-"`
	Page       int64  `json:"-"`
	Cursor     string `json:"-"`
}

func (p Paged) NumberOfPages() int64 {
	return p.PageCount
}

func (p Paged) CurrentPage() int64 {
	return p.Page
}

func (p Paged) NextCursor() string { return p.Cursor }

func (p Paged) HasNextCursor() bool { return p.Cursor != "" }

type OptionsApplier interface {
	Apply(req *resty.Request)
}

type PageOptions struct {
	Pagination string `url:"pagination,omitempty"`
	Limit      uint   `url:"limit,omitempty"`
	Page       uint   `url:"page,omitempty"`
	Cursor     string `url:"cursor,omitempty"`
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
	headerNextCursor = "X-Pagination-Next-Cursor"
)

func applyPaged(res *resty.Response, paged *Paged) {
	headers := res.Header()
	paged.Limit = headerInt64(headers, headerLimit)
	paged.TotalCount = headerInt64(headers, headerTotalCount)
	paged.PageCount = headerInt64(headers, headerPageCount)
	paged.Page = headerInt64(headers, headerPage)
	paged.Cursor = headers.Get(headerNextCursor)
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
