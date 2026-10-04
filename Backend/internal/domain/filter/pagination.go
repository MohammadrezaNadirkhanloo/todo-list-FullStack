package filter

import "math"

type PagedList[T any] struct {
	Items           []T   `json:"items"`
	Page            int   `json:"page"`
	PageSize        int   `json:"pageSize"`
	TotalRows       int64 `json:"totalRows"`
	TotalPages      int   `json:"totalPages"`
	HasPreviousPage bool  `json:"hasPreviousPage"`
	HasNextPage     bool  `json:"hasNextPage"`
}

func NewPagedList[T any](items []T, totalRows int64, page, pageSize int) *PagedList[T] {
	if items == nil {
		items = []T{}
	}

	totalPages := 1
	if pageSize > 0 {
		if n := int(math.Ceil(float64(totalRows) / float64(pageSize))); n > 1 {
			totalPages = n
		}
	}

	return &PagedList[T]{
		Items:           items,
		Page:            page,
		PageSize:        pageSize,
		TotalRows:       totalRows,
		TotalPages:      totalPages,
		HasPreviousPage: page > 1,
		HasNextPage:     page < totalPages,
	}
}

func MapPagedList[TIn any, TOut any](in *PagedList[TIn], mapFn func(TIn) TOut) *PagedList[TOut] {
	if in == nil {
		return nil
	}

	items := make([]TOut, 0, len(in.Items))
	for _, item := range in.Items {
		items = append(items, mapFn(item))
	}

	return &PagedList[TOut]{
		Items:           items,
		Page:            in.Page,
		PageSize:        in.PageSize,
		TotalRows:       in.TotalRows,
		TotalPages:      in.TotalPages,
		HasPreviousPage: in.HasPreviousPage,
		HasNextPage:     in.HasNextPage,
	}
}