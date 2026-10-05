package database

import (
	"fmt"
	"strings"

	"github.com/MohammadrezaNadirkhanloo/internal/domain/authz"
	"github.com/MohammadrezaNadirkhanloo/internal/domain/filter"
	"github.com/MohammadrezaNadirkhanloo/pkg/apperror"
	"gorm.io/gorm"
)

func ApplyScope(db *gorm.DB, scope authz.Scope, spec filter.Spec) (*gorm.DB, error) {
	if !scope.Allowed {
		return nil, forbidden(scope.Reason)
	}

	clause, args, err := BuildScopeWhere(scope, spec)
	if err != nil {
		return nil, err
	}
	if clause == "" {
		return db, nil
	}
	return db.Where(clause, args...), nil
}

func BuildScopeWhere(scope authz.Scope, spec filter.Spec) (string, []any, error) {
	if !scope.Allowed {
		return "", nil, forbidden(scope.Reason)
	}

	var (
		parts []string
		args  []any
	)

	if !scope.Unrestricted && len(scope.Allow) > 0 {
		clause, a, err := orGroup(scope.Allow, spec)
		if err != nil {
			return "", nil, err
		}
		parts = append(parts, clause)
		args = append(args, a...)
	}

	if len(scope.Deny) > 0 {
		clause, a, err := orGroup(scope.Deny, spec)
		if err != nil {
			return "", nil, err
		}
		parts = append(parts, "NOT "+clause)
		args = append(args, a...)
	}

	if len(parts) == 0 {
		return "", nil, nil
	}
	return strings.Join(parts, " AND "), args, nil
}

func orGroup(sets []authz.Conditions, spec filter.Spec) (string, []any, error) {
	parts := make([]string, 0, len(sets))
	args := make([]any, 0, len(sets))

	for _, set := range sets {
		clause, a, err := conditionsToSQL(set, spec)
		if err != nil {
			return "", nil, err
		}
		if clause == "" {
			return "", nil, nil
		}
		parts = append(parts, clause)
		args = append(args, a...)
	}

	if len(parts) == 0 {
		return "", nil, nil
	}
	return "(" + strings.Join(parts, " OR ") + ")", args, nil
}

func conditionsToSQL(c authz.Conditions, spec filter.Spec) (string, []any, error) {
	if len(c) == 0 {
		return "", nil, nil
	}

	keys := sortedKeys(c)

	parts := make([]string, 0, len(keys))
	args := make([]any, 0, len(keys))

	for _, key := range keys {
		field, ok := spec.Field(key)
		if !ok {
			return "", nil, apperror.Forbidden()
		}

		clause, a, err := conditionToSQL(field.Column, c[key])
		if err != nil {
			return "", nil, err
		}
		parts = append(parts, clause)
		args = append(args, a...)
	}

	if len(parts) == 1 {
		return parts[0], args, nil
	}
	return "(" + strings.Join(parts, " AND ") + ")", args, nil
}

func conditionToSQL(col string, expected any) (string, []any, error) {
	ops, isOperatorMap := operatorMap(expected)
	if !isOperatorMap {
		if expected == nil {
			return col + " IS NULL", nil, nil
		}
		return col + " = ?", []any{expected}, nil
	}

	parts := make([]string, 0, len(ops))
	args := make([]any, 0, len(ops))

	for _, op := range sortedKeys(ops) {
		operand := ops[op]

		switch op {
		case authz.OpNe:
			parts = append(parts, col+" IS DISTINCT FROM ?")
			args = append(args, operand)

		case authz.OpIn, authz.OpNin:
			values, err := toList(operand)
			if err != nil {
				return "", nil, err
			}
			if len(values) == 0 {
				if op == authz.OpIn {
					parts = append(parts, "FALSE")
				} else {
					parts = append(parts, "TRUE")
				}
				continue
			}
			ph := placeholders(len(values))
			if op == authz.OpIn {
				parts = append(parts, col+" IN ("+ph+")")
			} else {
				parts = append(parts, "("+col+" IS NULL OR "+col+" NOT IN ("+ph+"))")
			}
			args = append(args, values...)

		case authz.OpGt:
			parts = append(parts, col+" > ?")
			args = append(args, operand)
		case authz.OpGte:
			parts = append(parts, col+" >= ?")
			args = append(args, operand)
		case authz.OpLt:
			parts = append(parts, col+" < ?")
			args = append(args, operand)
		case authz.OpLte:
			parts = append(parts, col+" <= ?")
			args = append(args, operand)

		case authz.OpExists:
			exists, ok := operand.(bool)
			if !ok {
				return "", nil, apperror.Internal(fmt.Errorf(
					"authz: مقدار $exists باید بولی باشد"))
			}
			if exists {
				parts = append(parts, col+" IS NOT NULL")
			} else {
				parts = append(parts, col+" IS NULL")
			}

		default:
			return "", nil, apperror.Forbidden()
		}
	}

	if len(parts) == 1 {
		return parts[0], args, nil
	}
	return "(" + strings.Join(parts, " AND ") + ")", args, nil
}

func operatorMap(v any) (map[string]any, bool) {
	m, ok := v.(map[string]any)
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

func toList(v any) ([]any, error) {
	switch typed := v.(type) {
	case []any:
		return typed, nil
	case []string:
		out := make([]any, len(typed))
		for i, s := range typed {
			out[i] = s
		}
		return out, nil
	case []int64:
		out := make([]any, len(typed))
		for i, n := range typed {
			out[i] = n
		}
		return out, nil
	case []int:
		out := make([]any, len(typed))
		for i, n := range typed {
			out[i] = n
		}
		return out, nil
	default:
		return nil, apperror.Internal(fmt.Errorf(
			"authz: مقدار $in/$nin باید آرایه باشد"))
	}
}

func sortedKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	for i := 1; i < len(keys); i++ {
		for j := i; j > 0 && keys[j] < keys[j-1]; j-- {
			keys[j], keys[j-1] = keys[j-1], keys[j]
		}
	}
	return keys
}

func forbidden(reason string) error {
	if reason != "" {
		return apperror.New(apperror.CodeForbidden, reason)
	}
	return apperror.Forbidden()
}