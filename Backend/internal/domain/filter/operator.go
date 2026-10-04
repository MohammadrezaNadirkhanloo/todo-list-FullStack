package filter

type Operator string

const (
	OpContains     Operator = "contains"
	OpNotContains  Operator = "does_not_contain"
	OpStartsWith   Operator = "starts_with"
	OpEndsWith     Operator = "ends_with"

	OpEquals       Operator = "equals"
	OpNotEquals    Operator = "not_equals"

	OpGreaterThan    Operator = "gt"
	OpGreaterOrEqual Operator = "gte"
	OpLessThan       Operator = "lt"
	OpLessOrEqual    Operator = "lte"

	OpIn    Operator = "in"
	OpNotIn Operator = "not_in"

	OpIsEmpty    Operator = "is_empty"
	OpIsNotEmpty Operator = "is_not_empty"
	OpIsTrue     Operator = "is_true"
	OpIsFalse    Operator = "is_false"
)

var valuelessOperators = map[Operator]bool{
	OpIsEmpty:    true,
	OpIsNotEmpty: true,
	OpIsTrue:     true,
	OpIsFalse:    true,
}

var multiValueOperators = map[Operator]bool{
	OpIn:    true,
	OpNotIn: true,
}

func (o Operator) NeedsValue() bool { return !valuelessOperators[o] }

func (o Operator) IsMultiValue() bool { return multiValueOperators[o] }

var allowedOperators = map[Variant]map[Operator]bool{
	VariantText: set(
		OpContains, OpNotContains, OpEquals, OpNotEquals,
		OpStartsWith, OpEndsWith, OpIsEmpty, OpIsNotEmpty,
	),
	VariantLink: set(
		OpContains, OpEquals, OpIsEmpty, OpIsNotEmpty,
	),
	VariantNumber: set(
		OpEquals, OpNotEquals, OpGreaterThan, OpLessThan,
		OpGreaterOrEqual, OpLessOrEqual,
	),
	VariantCurrency: set(
		OpEquals, OpNotEquals, OpGreaterThan, OpLessThan,
		OpGreaterOrEqual, OpLessOrEqual,
	),
	VariantPercentage: set(
		OpEquals, OpGreaterOrEqual, OpLessOrEqual,
	),
	VariantRating: set(
		OpEquals, OpGreaterOrEqual, OpLessOrEqual,
	),
	VariantSelect: set(
		OpEquals, OpNotEquals,
	),
	VariantMultiSelect: set(
		OpIn, OpNotIn,
	),
	VariantTags: set(
		OpIn, OpNotIn, OpIsEmpty, OpIsNotEmpty,
	),
	VariantBoolean: set(
		OpIsTrue, OpIsFalse,
	),
	VariantDate: set(
		OpEquals, OpNotEquals, OpLessThan, OpGreaterThan,
		OpLessOrEqual, OpGreaterOrEqual, OpIsEmpty, OpIsNotEmpty,
	),
	VariantDateTime: set(
		OpEquals, OpNotEquals, OpLessThan, OpGreaterThan,
		OpLessOrEqual, OpGreaterOrEqual, OpIsEmpty, OpIsNotEmpty,
	),
	VariantTime: set(
		OpEquals, OpNotEquals, OpLessThan, OpGreaterThan,
		OpIsEmpty, OpIsNotEmpty,
	),
}

func set(ops ...Operator) map[Operator]bool {
	m := make(map[Operator]bool, len(ops))
	for _, o := range ops {
		m[o] = true
	}
	return m
}

func (o Operator) IsKnown() bool {
	for _, ops := range allowedOperators {
		if ops[o] {
			return true
		}
	}
	return false
}

func (o Operator) AllowedFor(v Variant) bool {
	ops, ok := allowedOperators[v]
	if !ok {
		return false
	}
	return ops[o]
}

func OperatorsFor(v Variant) []Operator {
	ops, ok := allowedOperators[v]
	if !ok {
		return nil
	}
	out := make([]Operator, 0, len(ops))
	for _, o := range operatorOrder {
		if ops[o] {
			out = append(out, o)
		}
	}
	return out
}

var operatorOrder = []Operator{
	OpEquals, OpNotEquals,
	OpContains, OpNotContains, OpStartsWith, OpEndsWith,
	OpGreaterThan, OpGreaterOrEqual, OpLessThan, OpLessOrEqual,
	OpIn, OpNotIn,
	OpIsEmpty, OpIsNotEmpty, OpIsTrue, OpIsFalse,
}