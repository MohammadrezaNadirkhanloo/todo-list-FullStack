package appctx

import "context"

type ctxKey int

const (
	keyRequestID ctxKey = iota
	keyUserID
	keyUsername
	keyRoles
)

func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, keyRequestID, id)
}

func RequestID(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(keyRequestID).(string)
	return v, ok
}

func WithUser(ctx context.Context, id int64, username string, roles []string) context.Context {
	ctx = context.WithValue(ctx, keyUserID, id)
	ctx = context.WithValue(ctx, keyUsername, username)
	return context.WithValue(ctx, keyRoles, roles)
}

func UserID(ctx context.Context) (int64, bool) {
	v, ok := ctx.Value(keyUserID).(int64)
	return v, ok
}

func Username(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(keyUsername).(string)
	return v, ok
}

func Roles(ctx context.Context) ([]string, bool) {
	v, ok := ctx.Value(keyRoles).([]string)
	return v, ok
}

func HasRole(ctx context.Context, allowed ...string) bool {
	roles, ok := Roles(ctx)
	if !ok {
		return false
	}
	for _, role := range roles {
		for _, want := range allowed {
			if role == want {
				return true
			}
		}
	}
	return false
}