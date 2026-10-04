package filter

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/MohammadrezaNadirkhanloo/pkg/apperror"
)

const (
	maxFilters     = 25
	maxSorts       = 5
	maxSearchRunes = 200
	maxInValues    = 100
)

const (
	paramPage     = "page"
	paramPageSize = "pageSize"
	paramSort     = "sort"
	paramFilter   = "filter"
	paramSearch   = "search"
)

func ParseAdvanced(values url.Values, spec Spec, limits PageLimits) (Query, error) {
	q := Query{}

	if err := parsePagination(values, &q, limits); err != nil {
		return q, err
	}
	if err := parseSearch(values, &q, spec); err != nil {
		return q, err
	}
	if err := parseSorts(values[paramSort], &q, spec); err != nil {
		return q, err
	}
	if err := parseFilters(values[paramFilter], &q, spec); err != nil {
		return q, err
	}

	return q, nil
}

func ParseSimple(values url.Values, spec Spec, limits PageLimits) (Query, error) {
	q := Query{}

	if err := parsePagination(values, &q, limits); err != nil {
		return q, err
	}
	if err := parseSearch(values, &q, spec); err != nil {
		return q, err
	}

	if raw := values.Get(paramSort); raw != "" {
		if err := parseSorts([]string{raw}, &q, spec); err != nil {
			return q, err
		}
	}

	for name, raw := range values {
		switch name {
		case paramPage, paramPageSize, paramSort, paramSearch:
			continue
		case paramFilter:
			return q, apperror.InvalidInput(
				"this endpoint does not accept advanced filters (filter=...). " +
					"Use the simple form field=value instead.")
		}

		field, ok := spec.Field(name)
		if !ok {
			return q, apperror.InvalidInput(fmt.Sprintf(
				"unknown parameter %q. Allowed fields: %s",
				name, spec.NamesText()))
		}
		if !field.Filterable {
			return q, apperror.InvalidInput(fmt.Sprintf(
				"field %q is not filterable.", name))
		}
		if len(raw) == 0 || raw[0] == "" {
			continue
		}

		op := simpleOperator(field.Variant)
		cond, err := buildCondition(name, field, op, raw[0])
		if err != nil {
			return q, err
		}
		q.Filters = append(q.Filters, cond)

		if len(q.Filters) > maxFilters {
			return q, tooManyFilters()
		}
	}

	return q, nil
}

func simpleOperator(v Variant) Operator {
	switch v {
	case VariantText, VariantLink:
		return OpContains
	case VariantMultiSelect, VariantTags:
		return OpIn
	case VariantBoolean:
		return OpIsTrue
	default:
		return OpEquals
	}
}

func parsePagination(values url.Values, q *Query, limits PageLimits) error {
	if raw := values.Get(paramPage); raw != "" {
		n, err := strconv.Atoi(raw)
		if err != nil || n < 1 {
			return apperror.InvalidInput(
				"parameter page must be a positive integer.")
		}
		q.Page = n
	}

	if raw := values.Get(paramPageSize); raw != "" {
		n, err := strconv.Atoi(raw)
		if err != nil || n < 1 {
			return apperror.InvalidInput(
				"parameter pageSize must be a positive integer.")
		}
		q.PageSize = n
	}

	q.Normalize(limits)
	return nil
}

func parseSearch(values url.Values, q *Query, spec Spec) error {
	raw := strings.TrimSpace(values.Get(paramSearch))
	if raw == "" {
		return nil
	}
	if len([]rune(raw)) > maxSearchRunes {
		return apperror.InvalidInput(fmt.Sprintf(
			"search term must not exceed %d characters.", maxSearchRunes))
	}
	if !spec.HasSearchable() {
		return apperror.InvalidInput(
			"global search is not defined for this resource.")
	}
	q.Search = raw
	return nil
}

func parseSorts(raws []string, q *Query, spec Spec) error {
	if len(raws) > maxSorts {
		return apperror.InvalidInput(fmt.Sprintf(
			"at most %d sort columns are allowed.", maxSorts))
	}

	seen := make(map[string]bool, len(raws))

	for _, raw := range raws {
		raw = strings.TrimSpace(raw)
		if raw == "" {
			continue
		}

		name, dir, hasDir := strings.Cut(raw, ":")
		name = strings.TrimSpace(name)
		dir = strings.ToLower(strings.TrimSpace(dir))

		field, ok := spec.Field(name)
		if !ok {
			return apperror.InvalidInput(fmt.Sprintf(
				"sorting by field %q is not allowed. Allowed fields: %s",
				name, spec.NamesText()))
		}
		if !field.Sortable {
			return apperror.InvalidInput(fmt.Sprintf(
				"field %q is not sortable.", name))
		}

		desc := false
		switch {
		case !hasDir || dir == "" || dir == "asc":
			desc = false
		case dir == "desc":
			desc = true
		default:
			return apperror.InvalidInput(fmt.Sprintf(
				"invalid sort direction %q; only asc or desc.", dir))
		}

		if seen[name] {
			continue
		}
		seen[name] = true

		q.Sorts = append(q.Sorts, Sort{Field: name, Desc: desc})
	}

	return nil
}

func parseFilters(raws []string, q *Query, spec Spec) error {
	if len(raws) > maxFilters {
		return tooManyFilters()
	}

	for _, raw := range raws {
		raw = strings.TrimSpace(raw)
		if raw == "" {
			continue
		}

		name, rest, ok := strings.Cut(raw, ":")
		if !ok {
			return apperror.InvalidInput(fmt.Sprintf(
				"invalid filter format %q; expected column:operator:value", raw))
		}
		opRaw, value, _ := strings.Cut(rest, ":")

		name = strings.TrimSpace(name)
		op := Operator(strings.TrimSpace(opRaw))

		field, ok := spec.Field(name)
		if !ok {
			return apperror.InvalidInput(fmt.Sprintf(
				"filtering on field %q is not allowed. Allowed fields: %s",
				name, spec.NamesText()))
		}
		if !field.Filterable {
			return apperror.InvalidInput(fmt.Sprintf(
				"field %q is not filterable.", name))
		}
		if !op.IsKnown() {
			return apperror.InvalidInput(fmt.Sprintf(
				"unknown operator %q.", op))
		}
		if !op.AllowedFor(field.Variant) {
			return apperror.InvalidInput(fmt.Sprintf(
				"operator %q is not allowed on field %q (type %s). Allowed operators: %s",
				op, name, field.Variant, joinOperators(OperatorsFor(field.Variant))))
		}

		cond, err := buildCondition(name, field, op, value)
		if err != nil {
			return err
		}
		q.Filters = append(q.Filters, cond)
	}

	return nil
}

func tooManyFilters() error {
	return apperror.InvalidInput(fmt.Sprintf(
		"at most %d filters are allowed at the same time.", maxFilters))
}

func joinOperators(ops []Operator) string {
	parts := make([]string, 0, len(ops))
	for _, o := range ops {
		parts = append(parts, string(o))
	}
	return strings.Join(parts, ", ")
}

func buildCondition(name string, field Field, op Operator, raw string) (Condition, error) {
	cond := Condition{Field: name, Op: op, Raw: raw}

	if !op.NeedsValue() {
		return cond, nil
	}

	if field.Variant == VariantBoolean {
		switch strings.ToLower(strings.TrimSpace(raw)) {
		case "true", "1":
			cond.Op = OpIsTrue
		case "false", "0":
			cond.Op = OpIsFalse
		default:
			return cond, apperror.InvalidInput(fmt.Sprintf(
				"value of field %q must be true or false.", name))
		}
		return cond, nil
	}

	if strings.TrimSpace(raw) == "" {
		return cond, apperror.InvalidInput(fmt.Sprintf(
			"no filter value provided for field %q.", name))
	}

	if op.IsMultiValue() {
		parts := strings.Split(raw, ",")
		if len(parts) > maxInValues {
			return cond, apperror.InvalidInput(fmt.Sprintf(
				"at most %d values are allowed for operator %q.", maxInValues, op))
		}
		values := make([]any, 0, len(parts))
		for _, p := range parts {
			p = strings.TrimSpace(p)
			if p == "" {
				continue
			}
			v, err := coerce(name, field, p)
			if err != nil {
				return cond, err
			}
			values = append(values, v)
		}
		if len(values) == 0 {
			return cond, apperror.InvalidInput(fmt.Sprintf(
				"value list for field %q is empty.", name))
		}
		cond.Values = values
		return cond, nil
	}

	v, err := coerce(name, field, strings.TrimSpace(raw))
	if err != nil {
		return cond, err
	}
	cond.Values = []any{v}
	return cond, nil
}

func coerce(name string, field Field, raw string) (any, error) {
	switch field.Variant {

	case VariantText, VariantLink:
		return raw, nil

	case VariantSelect, VariantMultiSelect, VariantTags:
		if !field.AllowsOption(raw) {
			return nil, apperror.InvalidInput(fmt.Sprintf(
				"value %q is not allowed for field %q. Allowed values: %s",
				raw, name, strings.Join(field.Options, ", ")))
		}
		return raw, nil

	case VariantNumber, VariantCurrency:
		return parseNumber(name, raw)

	case VariantPercentage:
		v, err := parseFloat(name, raw)
		if err != nil {
			return nil, err
		}
		if v < 0 || v > 100 {
			return nil, apperror.InvalidInput(fmt.Sprintf(
				"value of field %q must be between 0 and 100.", name))
		}
		return v, nil

	case VariantRating:
		v, err := strconv.Atoi(raw)
		if err != nil {
			return nil, apperror.InvalidInput(fmt.Sprintf(
				"value of field %q must be an integer.", name))
		}
		if v < 0 || v > field.RatingMax() {
			return nil, apperror.InvalidInput(fmt.Sprintf(
				"value of field %q must be between 0 and %d.", name, field.RatingMax()))
		}
		return int64(v), nil

	case VariantDate:
		t, err := time.ParseInLocation("2006-01-02", raw, time.UTC)
		if err != nil {
			return nil, apperror.InvalidInput(fmt.Sprintf(
				"value of field %q must be a date in YYYY-MM-DD format.", name))
		}
		return t, nil

	case VariantDateTime:
		t, err := parseDateTime(raw)
		if err != nil {
			return nil, apperror.InvalidInput(fmt.Sprintf(
				"value of field %q must be an ISO datetime (e.g. 2024-01-15T10:30:00).", name))
		}
		return t, nil

	case VariantTime:
		s, err := parseClock(raw)
		if err != nil {
			return nil, apperror.InvalidInput(fmt.Sprintf(
				"value of field %q must be a time in HH:mm or HH:mm:ss format.", name))
		}
		return s, nil

	case VariantBoolean:
		return nil, apperror.InvalidInput(fmt.Sprintf(
			"boolean field %q does not take a value.", name))

	default:
		return nil, apperror.InvalidInput(fmt.Sprintf(
			"field type %q is not supported.", name))
	}
}

func parseNumber(name, raw string) (any, error) {
	if i, err := strconv.ParseInt(raw, 10, 64); err == nil {
		return i, nil
	}
	return parseFloat(name, raw)
}

func parseFloat(name, raw string) (float64, error) {
	f, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return 0, apperror.InvalidInput(fmt.Sprintf(
			"value of field %q must be a number.", name))
	}
	return f, nil
}

var dateTimeLayouts = []string{
	time.RFC3339Nano,
	time.RFC3339,
	"2006-01-02T15:04:05",
	"2006-01-02T15:04",
	"2006-01-02 15:04:05",
	"2006-01-02",
}

func parseDateTime(raw string) (time.Time, error) {
	for _, layout := range dateTimeLayouts {
		if t, err := time.ParseInLocation(layout, raw, time.UTC); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("invalid datetime format")
}

func parseClock(raw string) (string, error) {
	for _, layout := range []string{"15:04:05", "15:04"} {
		if t, err := time.Parse(layout, raw); err == nil {
			return t.Format("15:04:05"), nil
		}
	}
	return "", fmt.Errorf("invalid time format")
}
