package apperror

import "net/http"

type Code string

const (
	CodeInternal        Code = "INTERNAL"
	CodeInvalidInput    Code = "INVALID_INPUT"
	CodeValidation      Code = "VALIDATION_FAILED"
	CodeNotFound        Code = "NOT_FOUND"
	CodeConflict        Code = "CONFLICT"
	CodeUnauthorized    Code = "UNAUTHORIZED"
	CodeForbidden       Code = "FORBIDDEN"
	CodeTokenExpired    Code = "TOKEN_EXPIRED"
	CodeTokenInvalid    Code = "TOKEN_INVALID"
	CodeTokenRequired   Code = "TOKEN_REQUIRED"
	CodeCSRFFailed      Code = "CSRF_FAILED"
	CodeRateLimited     Code = "RATE_LIMITED"
	CodeTimeout         Code = "TIMEOUT"
	CodeUnavailable     Code = "SERVICE_UNAVAILABLE"
	CodeRequestTooLarge Code = "REQUEST_TOO_LARGE"
)

var httpStatusByCode = map[Code]int{
	CodeInternal:        http.StatusInternalServerError,
	CodeInvalidInput:    http.StatusBadRequest,
	CodeValidation:      http.StatusUnprocessableEntity,
	CodeNotFound:        http.StatusNotFound,
	CodeConflict:        http.StatusConflict,
	CodeUnauthorized:    http.StatusUnauthorized,
	CodeForbidden:       http.StatusForbidden,
	CodeTokenExpired:    http.StatusUnauthorized,
	CodeTokenInvalid:    http.StatusUnauthorized,
	CodeTokenRequired:   http.StatusUnauthorized,
	CodeCSRFFailed:      http.StatusForbidden,
	CodeRateLimited:     http.StatusTooManyRequests,
	CodeTimeout:         http.StatusGatewayTimeout,
	CodeUnavailable:     http.StatusServiceUnavailable,
	CodeRequestTooLarge: http.StatusRequestEntityTooLarge,
}

func HTTPStatus(err error) int { // 500
	status, ok := httpStatusByCode[CodeOf(err)]
	if !ok {
		return http.StatusInternalServerError
	}
	return status
}
