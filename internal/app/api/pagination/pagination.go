package pagination

import (
	"net/http"
	"strconv"
)

const (
	DefaultPage    int64 = 1
	DefaultPerPage int64 = 25
	MaxPerPage     int64 = 1000
)

// Pagination is a struct for pagination
type Pagination struct {
	Page int64
	Per  int64
}

// NewPagination creates a new pagination instance
func NewPagination(r *http.Request) *Pagination {
	page := parseQueryParam(r, "page", DefaultPage)
	per := parseQueryParam(r, "per", DefaultPerPage)

	if page < 1 {
		page = DefaultPage
	}

	if per < 1 {
		per = DefaultPerPage
	}

	if per > MaxPerPage {
		per = MaxPerPage
	}

	return &Pagination{
		Page: page,
		Per:  per,
	}
}

func parseQueryParam(r *http.Request, key string, defaultValue int64) int64 {
	param := r.URL.Query().Get(key)
	if param == "" {
		return defaultValue
	}

	value, err := strconv.ParseInt(param, 10, 32)
	if err != nil {
		return defaultValue
	}

	return value
}

// Offset returns the pagination offset
func (p *Pagination) Offset() int64 {
	return (p.Page - 1) * p.Per
}
