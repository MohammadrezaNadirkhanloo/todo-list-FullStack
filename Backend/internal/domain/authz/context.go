package authz

import "context"

type ctxKey int

const keyAbility ctxKey = iota

func NewContext(ctx context.Context, a *Ability) context.Context {
	return context.WithValue(ctx, keyAbility, a)
}

func FromContext(ctx context.Context) (*Ability, bool) {
	a, ok := ctx.Value(keyAbility).(*Ability)
	return a, ok && a != nil
}

func MustFromContext(ctx context.Context) *Ability {
	if a, ok := FromContext(ctx); ok {
		return a
	}
	return NewAbility(nil)
}