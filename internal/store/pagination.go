package store

import (
	"net/http"
	"strconv"
)

type PaginationQuery struct {
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
	Sort   string `json:"sort"`
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

	if fq.Limit == 0 {
		fq.Limit = 20 // Default limit
	}

	offset := qs.Get("offset")
	if offset != "" {
		l, err := strconv.Atoi(offset)
		if err != nil {
			return fq, nil
		}
		fq.Offset = l
	}

	sort := qs.Get("sort")
	if sort != "" {
		fq.Sort = sort
	} else {
		fq.Sort = "desc" // Default sort
	}

	search := qs.Get("search")
	if search != "" {
		fq.Search = search
	}

	return fq, nil
}
