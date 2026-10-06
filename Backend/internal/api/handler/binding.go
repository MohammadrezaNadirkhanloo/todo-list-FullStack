package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/MohammadrezaNadirkhanloo/internal/api/response"
	"github.com/MohammadrezaNadirkhanloo/internal/api/validation"
	"github.com/MohammadrezaNadirkhanloo/internal/domain/filter"
	"github.com/MohammadrezaNadirkhanloo/pkg/apperror"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func BindJSON[T any](c *gin.Context, target *T) bool {
	if err := c.ShouldBindJSON(target); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			response.FailValidation(c, "some fields are invalid.", validation.Translate(err))
			return false
		}

		response.Fail(c, apperror.InvalidInput(
			"request body is invalid. Please check the JSON format and field types."))
		return false
	}
	return true
}

func BindQuery[T any](c *gin.Context, target *T) bool {
	if err := c.ShouldBindQuery(target); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			response.FailValidation(c, "request parameters are invalid.", validation.Translate(err))
			return false
		}
		response.Fail(c, apperror.InvalidInput("request parameters are invalid."))
		return false
	}
	return true
}

func PathID(c *gin.Context) (int64, bool) {
	raw := c.Param("id")

	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || id <= 0 {
		response.Fail(c, apperror.InvalidInput("the provided ID is invalid."))
		return 0, false
	}
	return id, true
}

func ParseListQuery(
	c *gin.Context,
	spec filter.Spec,
	limits filter.PageLimits,
	mode ListMode,
) (filter.Query, bool) {
	values := c.Request.URL.Query()

	var (
		q   filter.Query
		err error
	)
	if mode == ListAdvanced {
		q, err = filter.ParseAdvanced(values, spec, limits)
	} else {
		q, err = filter.ParseSimple(values, spec, limits)
	}

	if err != nil {
		response.Fail(c, err)
		return q, false
	}
	return q, true
}

func respondCreated(c *gin.Context, location string, body any) {
	if location != "" {
		c.Header("Location", location)
	}
	c.JSON(http.StatusCreated, response.Envelope{Success: true, Result: body})
}
