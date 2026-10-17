package validation

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"unicode"

	"github.com/MohammadrezaNadirkhanloo/internal/api/response"
	"github.com/MohammadrezaNadirkhanloo/internal/config"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

var (
	iranianMobileRe = regexp.MustCompile(`^09(1[0-9]|2[0-2]|3[0-9]|9[0-9])[0-9]{7}$`)
	slugRe          = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

	registerOnce sync.Once
)

func Register(cfg config.PasswordConfig) error {
	var regErr error

	registerOnce.Do(func() {
		v, ok := binding.Validator.Engine().(*validator.Validate)
		if !ok {
			regErr = errors.New("validation: gin validation engine is not available")
			return
		}

		validators := map[string]validator.Func{
			"mobile":   iranianMobile,
			"slug":     slug,
			"password": passwordValidator(cfg),
		}

		for tag, fn := range validators {
			if err := v.RegisterValidation(tag, fn); err != nil {
				regErr = fmt.Errorf("validation: failed to register validator %q: %w", tag, err)
				return
			}
		}
	})

	return regErr
}

func iranianMobile(fl validator.FieldLevel) bool {
	value, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	return iranianMobileRe.MatchString(value)
}

func slug(fl validator.FieldLevel) bool {
	value, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	return slugRe.MatchString(value)
}

func passwordValidator(cfg config.PasswordConfig) validator.Func {
	return func(fl validator.FieldLevel) bool {
		value, ok := fl.Field().Interface().(string)
		if !ok {
			return false
		}
		return CheckPassword(value, cfg) == nil
	}
}

func CheckPassword(pw string, cfg config.PasswordConfig) error {
	length := len([]rune(pw))

	if length < cfg.MinLength {
		return fmt.Errorf("password must be at least %d characters long", cfg.MinLength)
	}
	if length > cfg.MaxLength {
		return fmt.Errorf("password must not exceed %d characters", cfg.MaxLength)
	}
	if cfg.IncludeDigits && !containsFunc(pw, unicode.IsDigit) {
		return errors.New("password must contain at least one digit")
	}
	if cfg.IncludeUppercase && !containsFunc(pw, unicode.IsUpper) {
		return errors.New("password must contain at least one uppercase letter")
	}
	if cfg.IncludeLowercase && !containsFunc(pw, unicode.IsLower) {
		return errors.New("password must contain at least one lowercase letter")
	}
	if cfg.IncludeSymbols && !containsFunc(pw, func(r rune) bool {
		return unicode.IsPunct(r) || unicode.IsSymbol(r)
	}) {
		return errors.New("password must contain at least one symbol")
	}
	return nil
}

func containsFunc(s string, fn func(rune) bool) bool {
	for _, r := range s {
		if fn(r) {
			return true
		}
	}
	return false
}

func Translate(err error) []response.FieldError {
	var ve validator.ValidationErrors
	if !errors.As(err, &ve) {
		return nil
	}

	out := make([]response.FieldError, 0, len(ve))
	for _, fe := range ve {
		out = append(out, response.FieldError{
			Field:   fieldName(fe),
			Message: message(fe),
		})
	}
	return out
}

func fieldName(fe validator.FieldError) string {
	name := fe.Field()
	if name == "" {
		return name
	}
	r := []rune(name)
	r[0] = unicode.ToLower(r[0])
	return string(r)
}

func message(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "this field is required."
	case "email":
		return "invalid email format."
	case "min":
		return fmt.Sprintf("value must be at least %s.", fe.Param())
	case "max":
		return fmt.Sprintf("value must be at most %s.", fe.Param())
	case "len":
		return fmt.Sprintf("length must be exactly %s.", fe.Param())
	case "numeric":
		return "only digits are allowed."
	case "alphanum":
		return "only English letters and digits are allowed."
	case "mobile":
		return "invalid mobile number (example: 09123456789)."
	case "slug":
		return "slug may only contain lowercase English letters, digits, and hyphens."
	case "password":
		return "password does not meet the requirements (length, uppercase, lowercase, digit)."
	case "oneof":
		return fmt.Sprintf("value must be one of: %s",
			strings.ReplaceAll(fe.Param(), " ", ", "))
	case "gte":
		return fmt.Sprintf("value must be greater than or equal to %s.", fe.Param())
	case "lte":
		return fmt.Sprintf("value must be less than or equal to %s.", fe.Param())
	default:
		return fmt.Sprintf("value of this field is invalid (%s).", fe.Tag())
	}
}
