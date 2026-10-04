package filter

import (
	"fmt"
	"strings"

	"github.com/MohammadrezaNadirkhanloo/pkg/apperror"
)

type Condition struct {
	Field  string
	Op     Operator
	Raw    string
	Values []any
}

type Sort struct {
	Field string
	Desc  bool
}

func (s Sort) Direction() string {
	if s.Desc {
		return "desc"
	}
	return "asc"
}

type Query struct {
	Page     int
	PageSize int
	Search   string
	Filters  []Condition
	Sorts    []Sort
}

type PageLimits struct {
	Default int
	Max     int
}

func (p PageLimits) sanitize() PageLimits {
	if p.Default <= 0 {
		p.Default = 20
	}
	if p.Max <= 0 {
		p.Max = 100
	}
	if p.Default > p.Max {
		p.Default = p.Max
	}
	return p
}

func (q *Query) Normalize(limits PageLimits) {
	limits = limits.sanitize()

	if q.Page < 1 {
		q.Page = 1
	}
	if q.PageSize <= 0 {
		q.PageSize = limits.Default
	}
	if q.PageSize > limits.Max {
		q.PageSize = limits.Max
	}
}

func (q Query) Offset() int { return (q.Page - 1) * q.PageSize }

func (q Query) Limit() int { return q.PageSize }

func (q Query) Validate(spec Spec) error {
	for _, c := range q.Filters {
		f, ok := spec.Field(c.Field)
		if !ok {
			return apperror.InvalidInput(fmt.Sprintf(
				"filtering on field %q is not allowed. Allowed fields: %s",
				c.Field, spec.NamesText()))
		}
		if !f.Filterable {
			return apperror.InvalidInput(fmt.Sprintf(
				"field %q is not filterable.", c.Field))
		}
		if !c.Op.AllowedFor(f.Variant) {
			return apperror.InvalidInput(fmt.Sprintf(
				"operator %q is not allowed on field %q (type %s).",
				c.Op, c.Field, f.Variant))
		}
	}

	for _, s := range q.Sorts {
		f, ok := spec.Field(s.Field)
		if !ok {
			return apperror.InvalidInput(fmt.Sprintf(
				"sorting by field %q is not allowed. Allowed fields: %s",
				s.Field, spec.NamesText()))
		}
		if !f.Sortable {
			return apperror.InvalidInput(fmt.Sprintf(
				"field %q is not sortable.", s.Field))
		}
	}

	if q.Search != "" && !spec.HasSearchable() {
		return apperror.InvalidInput(
			"global search is not defined for this resource.")
	}

	return nil
}

func (q Query) String() string {
	parts := make([]string, 0, len(q.Filters)+len(q.Sorts)+1)
	for _, c := range q.Filters {
		parts = append(parts, string(c.Op)+"("+c.Field+")")
	}
	for _, s := range q.Sorts {
		parts = append(parts, "sort("+s.Field+" "+s.Direction()+")")
	}
	if q.Search != "" {
		parts = append(parts, "search(*)")
	}
	return fmt.Sprintf("page=%d size=%d %s",
		q.Page, q.PageSize, strings.Join(parts, " "))
}
