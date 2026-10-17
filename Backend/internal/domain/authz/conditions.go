package authz

import (
	"strings"
	"time"
)

type Conditions map[string]any

const (
	OpNe     = "$ne"
	OpIn     = "$in"
	OpNin    = "$nin"
	OpGt     = "$gt"
	OpGte    = "$gte"
	OpLt     = "$lt"
	OpLte    = "$lte"
	OpExists = "$exists"
)

type Resource map[string]any

func (c Conditions) Match(res Resource) bool {
	if len(c) == 0 {
		return true
	}
	if res == nil {
		return false
	}

	for path, expected := range c {
		actual, found := lookup(res, path)
		if !matchOne(expected, actual, found) {
			return false
		}
	}
	return true
}

func lookup(res Resource, path string) (any, bool) {
	if !strings.Contains(path, ".") {
		v, ok := res[path]
		return v, ok
	}

	var current any = map[string]any(res)
	for _, part := range strings.Split(path, ".") {
		m, ok := toMap(current)
		if !ok {
			return nil, false
		}
		current, ok = m[part]
		if !ok {
			return nil, false
		}
	}
	return current, true
}

func toMap(v any) (map[string]any, bool) {
	switch m := v.(type) {
	case map[string]any:
		return m, true
	case Resource:
		return m, true
	default:
		return nil, false
	}
}

func matchOne(expected, actual any, found bool) bool {
	ops, isOperatorMap := asOperatorMap(expected)
	if !isOperatorMap {
		return found && equal(expected, actual)
	}

	for op, operand := range ops {
		switch op {
		case OpExists:
			want, ok := operand.(bool)
			if !ok || want != found {
				return false
			}

		case OpNe:
			if found && equal(operand, actual) {
				return false
			}

		case OpIn:
			if !found || !inList(operand, actual) {
				return false
			}

		case OpNin:
			if found && inList(operand, actual) {
				return false
			}

		case OpGt, OpGte, OpLt, OpLte:
			if !found || !compareOp(op, actual, operand) {
				return false
			}

		default:
			return false
		}
	}
	return true
}

func asOperatorMap(v any) (map[string]any, bool) {
	m, ok := toMap(v)
	if !ok || len(m) == 0 {
		return nil, false
	}
	for k := range m {
		if !strings.HasPrefix(k, "$") {
			return nil, false
		}
	}
	return m, true
}

func inList(list, actual any) bool {
	items, ok := list.([]any)
	if !ok {
		switch typed := list.(type) {
		case []string:
			items = make([]any, len(typed))
			for i, s := range typed {
				items[i] = s
			}
		case []int64:
			items = make([]any, len(typed))
			for i, n := range typed {
				items[i] = n
			}
		case []int:
			items = make([]any, len(typed))
			for i, n := range typed {
				items[i] = n
			}
		default:
			return false
		}
	}

	for _, item := range items {
		if equal(item, actual) {
			return true
		}
	}
	return false
}

func equal(a, b any) bool {
	if a == nil || b == nil {
		return a == nil && b == nil
	}

	if af, aok := toFloat(a); aok {
		if bf, bok := toFloat(b); bok {
			return af == bf
		}
		return false
	}

	if at, aok := toTime(a); aok {
		if bt, bok := toTime(b); bok {
			return at.Equal(bt)
		}
		return false
	}

	as, aok := a.(string)
	bs, bok := b.(string)
	if aok && bok {
		return as == bs
	}

	ab, aok := a.(bool)
	bb, bok := b.(bool)
	if aok && bok {
		return ab == bb
	}

	return false
}

func compareOp(op string, actual, operand any) bool {
	if af, aok := toFloat(actual); aok {
		if bf, bok := toFloat(operand); bok {
			return compareFloat(op, af, bf)
		}
		return false
	}

	if at, aok := toTime(actual); aok {
		if bt, bok := toTime(operand); bok {
			switch op {
			case OpGt:
				return at.After(bt)
			case OpGte:
				return at.After(bt) || at.Equal(bt)
			case OpLt:
				return at.Before(bt)
			case OpLte:
				return at.Before(bt) || at.Equal(bt)
			}
		}
		return false
	}

	as, aok := actual.(string)
	bs, bok := operand.(string)
	if aok && bok {
		switch op {
		case OpGt:
			return as > bs
		case OpGte:
			return as >= bs
		case OpLt:
			return as < bs
		case OpLte:
			return as <= bs
		}
	}

	return false
}

func compareFloat(op string, a, b float64) bool {
	switch op {
	case OpGt:
		return a > b
	case OpGte:
		return a >= b
	case OpLt:
		return a < b
	case OpLte:
		return a <= b
	default:
		return false
	}
}

func toFloat(v any) (float64, bool) {
	switch n := v.(type) {
	case int:
		return float64(n), true
	case int8:
		return float64(n), true
	case int16:
		return float64(n), true
	case int32:
		return float64(n), true
	case int64:
		return float64(n), true
	case uint:
		return float64(n), true
	case uint8:
		return float64(n), true
	case uint16:
		return float64(n), true
	case uint32:
		return float64(n), true
	case uint64:
		return float64(n), true
	case float32:
		return float64(n), true
	case float64:
		return n, true
	default:
		return 0, false
	}
}

func toTime(v any) (time.Time, bool) {
	switch t := v.(type) {
	case time.Time:
		return t, true
	case *time.Time:
		if t == nil {
			return time.Time{}, false
		}
		return *t, true
	default:
		return time.Time{}, false
	}
}