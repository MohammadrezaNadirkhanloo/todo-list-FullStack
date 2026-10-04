package database

import (
	"fmt"
	"strings"
	"time"

	"github.com/MohammadrezaNadirkhanloo/internal/domain/filter"
	"github.com/MohammadrezaNadirkhanloo/pkg/apperror"
	"gorm.io/gorm"
)

func BuildWhere(q filter.Query, spec filter.Spec) (string, []any, error) {
	clauses := make([]string, 0, len(q.Filters)+1)
	args := make([]any, 0, len(q.Filters)+1)

	for _, cond := range q.Filters {
		field, ok := spec.Field(cond.Field)
		if !ok || !field.Filterable {
			return "", nil, apperror.InvalidInput(fmt.Sprintf(
				"filtering on field %q is not allowed.", cond.Field))
		}
		if !cond.Op.AllowedFor(field.Variant) {
			return "", nil, apperror.InvalidInput(fmt.Sprintf(
				"operator %q is not allowed on field %q.", cond.Op, cond.Field))
		}

		clause, condArgs, err := buildCondition(field, cond)
		if err != nil {
			return "", nil, err
		}
		if clause == "" {
			continue
		}
		clauses = append(clauses, clause)
		args = append(args, condArgs...)
	}

	if searchClause, searchArgs := buildSearch(q.Search, spec); searchClause != "" {
		clauses = append(clauses, searchClause)
		args = append(args, searchArgs...)
	}

	return strings.Join(clauses, " AND "), args, nil
}

func buildCondition(field filter.Field, cond filter.Condition) (string, []any, error) {
	col := field.Column

	switch field.Variant {

	case filter.VariantText, filter.VariantLink, filter.VariantSelect:
		return textCondition(col, cond)

	case filter.VariantNumber, filter.VariantCurrency,
		filter.VariantPercentage, filter.VariantRating:
		return comparableCondition(col, cond)

	case filter.VariantMultiSelect, filter.VariantTags:
		return setCondition(col, field.Array, cond)

	case filter.VariantBoolean:
		return boolCondition(col, cond)

	case filter.VariantDate:
		return dateCondition(col, cond)

	case filter.VariantDateTime:
		return comparableCondition(col, cond)

	case filter.VariantTime:
		return timeCondition(col, cond)

	default:
		return "", nil, apperror.InvalidInput(fmt.Sprintf(
			"field type %q is not supported.", field.Variant))
	}
}

func textCondition(col string, cond filter.Condition) (string, []any, error) {
	switch cond.Op {
	case filter.OpEquals:
		v, err := one(cond)
		return col + " = ?", []any{v}, err

	case filter.OpNotEquals:
		v, err := one(cond)
		return col + " IS DISTINCT FROM ?", []any{v}, err

	case filter.OpContains:
		v, err := oneString(cond)
		return likeClause(col, false), []any{"%" + escapeLike(v) + "%"}, err

	case filter.OpNotContains:
		v, err := oneString(cond)
		return "(" + col + " IS NULL OR " + likeClause(col, true) + ")",
			[]any{"%" + escapeLike(v) + "%"}, err

	case filter.OpStartsWith:
		v, err := oneString(cond)
		return likeClause(col, false), []any{escapeLike(v) + "%"}, err

	case filter.OpEndsWith:
		v, err := oneString(cond)
		return likeClause(col, false), []any{"%" + escapeLike(v)}, err

	case filter.OpIsEmpty:
		return "(" + col + " IS NULL OR " + col + " = '')", nil, nil

	case filter.OpIsNotEmpty:
		return "(" + col + " IS NOT NULL AND " + col + " <> '')", nil, nil

	default:
		return "", nil, unsupported(cond.Op, col)
	}
}

func comparableCondition(col string, cond filter.Condition) (string, []any, error) {
	switch cond.Op {
	case filter.OpEquals:
		v, err := one(cond)
		return col + " = ?", []any{v}, err
	case filter.OpNotEquals:
		v, err := one(cond)
		return col + " IS DISTINCT FROM ?", []any{v}, err
	case filter.OpGreaterThan:
		v, err := one(cond)
		return col + " > ?", []any{v}, err
	case filter.OpGreaterOrEqual:
		v, err := one(cond)
		return col + " >= ?", []any{v}, err
	case filter.OpLessThan:
		v, err := one(cond)
		return col + " < ?", []any{v}, err
	case filter.OpLessOrEqual:
		v, err := one(cond)
		return col + " <= ?", []any{v}, err
	case filter.OpIsEmpty:
		return col + " IS NULL", nil, nil
	case filter.OpIsNotEmpty:
		return col + " IS NOT NULL", nil, nil
	default:
		return "", nil, unsupported(cond.Op, col)
	}
}

func setCondition(col string, array bool, cond filter.Condition) (string, []any, error) {
	switch cond.Op {
	case filter.OpIn:
		if len(cond.Values) == 0 {
			return "", nil, missingValue(col)
		}
		if array {
			return col + " && " + arrayLiteral(len(cond.Values)), cond.Values, nil
		}
		return col + " IN (" + placeholders(len(cond.Values)) + ")", cond.Values, nil

	case filter.OpNotIn:
		if len(cond.Values) == 0 {
			return "", nil, missingValue(col)
		}
		if array {
			return "(" + col + " IS NULL OR NOT (" + col + " && " +
				arrayLiteral(len(cond.Values)) + "))", cond.Values, nil
		}
		return "(" + col + " IS NULL OR " + col + " NOT IN (" +
			placeholders(len(cond.Values)) + "))", cond.Values, nil

	case filter.OpIsEmpty:
		if array {
			return "(" + col + " IS NULL OR cardinality(" + col + ") = 0)", nil, nil
		}
		return "(" + col + " IS NULL OR " + col + " = '')", nil, nil

	case filter.OpIsNotEmpty:
		if array {
			return "(" + col + " IS NOT NULL AND cardinality(" + col + ") > 0)", nil, nil
		}
		return "(" + col + " IS NOT NULL AND " + col + " <> '')", nil, nil

	default:
		return "", nil, unsupported(cond.Op, col)
	}
}

func boolCondition(col string, cond filter.Condition) (string, []any, error) {
	switch cond.Op {
	case filter.OpIsTrue:
		return col + " IS TRUE", nil, nil
	case filter.OpIsFalse:
		return col + " IS FALSE", nil, nil
	default:
		return "", nil, unsupported(cond.Op, col)
	}
}

func dateCondition(col string, cond filter.Condition) (string, []any, error) {
	if cond.Op == filter.OpIsEmpty {
		return col + " IS NULL", nil, nil
	}
	if cond.Op == filter.OpIsNotEmpty {
		return col + " IS NOT NULL", nil, nil
	}

	v, err := one(cond)
	if err != nil {
		return "", nil, err
	}
	day, ok := v.(time.Time)
	if !ok {
		return "", nil, apperror.InvalidInput(fmt.Sprintf(
			"value of field %q must be a date.", cond.Field))
	}
	next := day.AddDate(0, 0, 1)

	switch cond.Op {
	case filter.OpEquals:
		return "(" + col + " >= ? AND " + col + " < ?)", []any{day, next}, nil
	case filter.OpNotEquals:
		return "(" + col + " IS NULL OR " + col + " < ? OR " + col + " >= ?)",
			[]any{day, next}, nil
	case filter.OpLessThan:
		return col + " < ?", []any{day}, nil
	case filter.OpGreaterOrEqual:
		return col + " >= ?", []any{day}, nil
	case filter.OpGreaterThan:
		return col + " >= ?", []any{next}, nil
	case filter.OpLessOrEqual:
		return col + " < ?", []any{next}, nil
	default:
		return "", nil, unsupported(cond.Op, col)
	}
}

func timeCondition(col string, cond filter.Condition) (string, []any, error) {
	if cond.Op == filter.OpIsEmpty {
		return col + " IS NULL", nil, nil
	}
	if cond.Op == filter.OpIsNotEmpty {
		return col + " IS NOT NULL", nil, nil
	}

	v, err := one(cond)
	if err != nil {
		return "", nil, err
	}
	lhs := "CAST(" + col + " AS time)"

	switch cond.Op {
	case filter.OpEquals:
		return lhs + " = CAST(? AS time)", []any{v}, nil
	case filter.OpNotEquals:
		return lhs + " IS DISTINCT FROM CAST(? AS time)", []any{v}, nil
	case filter.OpLessThan:
		return lhs + " < CAST(? AS time)", []any{v}, nil
	case filter.OpGreaterThan:
		return lhs + " > CAST(? AS time)", []any{v}, nil
	default:
		return "", nil, unsupported(cond.Op, col)
	}
}

func buildSearch(term string, spec filter.Spec) (string, []any) {
	term = strings.TrimSpace(term)
	if term == "" {
		return "", nil
	}

	fields := spec.SearchableFields()
	if len(fields) == 0 {
		return "", nil
	}

	pattern := "%" + escapeLike(term) + "%"
	parts := make([]string, 0, len(fields))
	args := make([]any, 0, len(fields))

	for _, f := range fields {
		lhs := f.Column
		switch f.Variant {
		case filter.VariantText, filter.VariantLink, filter.VariantSelect:
		default:
			lhs = "CAST(" + f.Column + " AS text)"
		}
		parts = append(parts, likeClause(lhs, false))
		args = append(args, pattern)
	}

	return "(" + strings.Join(parts, " OR ") + ")", args
}

func BuildOrder(q filter.Query, spec filter.Spec) string {
	parts := make([]string, 0, len(q.Sorts)+1)

	for _, s := range q.Sorts {
		field, ok := spec.Field(s.Field)
		if !ok || !field.Sortable {
			continue
		}
		parts = append(parts, field.Column+" "+s.Direction())
	}

	if len(parts) == 0 {
		return spec.DefaultSort()
	}
	return strings.Join(parts, ", ")
}

func ApplyFilters(db *gorm.DB, q filter.Query, spec filter.Spec) (*gorm.DB, error) {
	clause, args, err := BuildWhere(q, spec)
	if err != nil {
		return nil, err
	}
	if clause == "" {
		return db, nil
	}
	return db.Where(clause, args...), nil
}

func ApplyQuery(db *gorm.DB, q filter.Query, spec filter.Spec) (*gorm.DB, error) {
	db, err := ApplyFilters(db, q, spec)
	if err != nil {
		return nil, err
	}
	return db.
		Order(BuildOrder(q, spec)).
		Offset(q.Offset()).
		Limit(q.Limit()), nil
}

func one(cond filter.Condition) (any, error) {
	if len(cond.Values) != 1 {
		return nil, missingValue(cond.Field)
	}
	return cond.Values[0], nil
}

func oneString(cond filter.Condition) (string, error) {
	v, err := one(cond)
	if err != nil {
		return "", err
	}
	s, ok := v.(string)
	if !ok {
		return "", apperror.InvalidInput(fmt.Sprintf(
			"value of field %q must be a string.", cond.Field))
	}
	return s, nil
}

func missingValue(field string) error {
	return apperror.InvalidInput(fmt.Sprintf(
		"no filter value provided for field %q.", field))
}

func unsupported(op filter.Operator, col string) error {
	return apperror.InvalidInput(fmt.Sprintf(
		"operator %q is not supported on column %q.", op, col))
}

func placeholders(n int) string {
	if n <= 0 {
		return ""
	}
	return strings.TrimSuffix(strings.Repeat("?, ", n), ", ")
}

func arrayLiteral(n int) string {
	return "ARRAY[" + placeholders(n) + "]::text[]"
}

func likeClause(col string, negate bool) string {
	op := "ILIKE"
	if negate {
		op = "NOT ILIKE"
	}
	return fmt.Sprintf(`%s %s ? ESCAPE '\'`, col, op)
}

func escapeLike(s string) string {
	replacer := strings.NewReplacer(
		`\`, `\\`,
		`%`, `\%`,
		`_`, `\_`,
	)
	return replacer.Replace(s)
}
