package middleware

import (
	"context"
	"errors"
	"net/http"
)

const (
	accessCookieBase = "access_token"
	csrfCookieBase   = "csrf_token"
	hostCookiePrefix = "__Host-"
)

const CSRFHeaderName = "X-CSRF-Token"

func AccessCookieName(insecure bool) string {
	if insecure {
		return accessCookieBase
	}
	return hostCookiePrefix + accessCookieBase
}

func CSRFCookieName(insecure bool) string {
	if insecure {
		return csrfCookieBase
	}
	return hostCookiePrefix + csrfCookieBase
}

var (
	ErrCookieMissing   = errors.New("middleware: cookie is missing")
	ErrCookieDuplicate = errors.New("middleware: duplicate cookie")
)

func SingleCookie(r *http.Request, name string) (string, error) {
	var (
		value string
		found int
	)
	for _, ck := range r.Cookies() {
		if ck.Name != name {
			continue
		}
		found++
		value = ck.Value
	}
	switch found {
	case 0:
		return "", ErrCookieMissing
	case 1:
		return value, nil
	default:
		return "", ErrCookieDuplicate
	}
}

type tokenIDKey struct{}

func withTokenID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, tokenIDKey{}, id)
}

func tokenIDFrom(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(tokenIDKey{}).(string)
	return id, ok && id != ""
}