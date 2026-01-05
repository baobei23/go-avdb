package store

import (
	"net/http"
	"strconv"
)

type PaginationQuery struct {
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
	Page   int    `json:"page"`
	Search string `json:"search"`
}

func (fq PaginationQuery) Parse(r *http.Request) (PaginationQuery, error) {
	qs := r.URL.Query()

	limit := qs.Get("limit")
	if limit != "" {
		l, err := strconv.Atoi(limit)
		if err != nil {
			return fq, nil
		}
		fq.Limit = l
	}

	if fq.Limit < 1 || fq.Limit > 100 {
		fq.Limit = 20 // Default limit
	}

	page := qs.Get("page")
	if page != "" {
		p, err := strconv.Atoi(page)
		if err != nil {
			return fq, nil
		}
		fq.Page = p
	}

	if fq.Page < 1 {
		fq.Page = 1
	}

	fq.Offset = (fq.Page - 1) * fq.Limit

	search := qs.Get("search")
	if search != "" {
		fq.Search = search
	}

	return fq, nil
}
