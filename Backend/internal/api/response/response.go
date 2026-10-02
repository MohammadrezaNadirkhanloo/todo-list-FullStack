package response

import (
	"net/http"

	"github.com/MohammadrezaNadirkhanloo/pkg/appctx"
	"github.com/MohammadrezaNadirkhanloo/pkg/apperror"
	"github.com/gin-gonic/gin"
)

type Envelope struct {
	Success bool       `json:"success"`
	Result  any        `json:"result,omitempty"`
	Message string     `json:"message,omitempty"`
	Error   *ErrorBody `json:"error,omitempty"`
}

type ErrorBody struct {
	Code       apperror.Code `json:"code"`
	Message    string        `json:"message"`
	Validation []FieldError  `json:"validation,omitempty"`
}

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// //////
type ListResponse[T any] struct {
	Data       []T   `json:"data"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"totalPages"`
	Page       int   `json:"page"`
	PageSize   int   `json:"pageSize"`
}

// func NewListResponse[T any](p *filter.PagedList[T]) ListResponse[T] {
// 	if p == nil {
// 		return ListResponse[T]{Data: []T{}, TotalPages: 1, Page: 1}
// 	}
// 	return ListResponse[T]{
// 		Data:       p.Items,
// 		Total:      p.TotalRows,
// 		TotalPages: p.TotalPages,
// 		Page:       p.Page,
// 		PageSize:   p.PageSize,
// 	}
// }

// success
// func List[T any](c *gin.Context, p *filter.PagedList[T]) {
// 	c.JSON(http.StatusOK, NewListResponse(p))
// }

func Raw(c *gin.Context, status int, body any) {
	c.JSON(status, body)
}

func OK(c *gin.Context, result any) {
	c.JSON(http.StatusOK, Envelope{Success: true, Result: result})
}

func Created(c *gin.Context, result any) {
	c.JSON(http.StatusCreated, Envelope{Success: true, Result: result})
}

func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// error
func Fail(c *gin.Context, err error) {
	message := apperror.UserMessage(err)

	body := &ErrorBody{
		Code:    apperror.CodeOf(err),
		Message: message,
	}

	//Validation
	var fieldErrors []FieldError
	if fe, ok := err.(interface{ Fields() []FieldError }); ok {
		fieldErrors = fe.Fields()
	}
	body.Validation = fieldErrors

	c.AbortWithStatusJSON(apperror.HTTPStatus(err), Envelope{
		Success:   false,
		Message:   message,
		Error:     body,
	})
	//return
}

// json error
func FailValidation(c *gin.Context, message string, fields []FieldError) {
	c.AbortWithStatusJSON(http.StatusUnprocessableEntity, Envelope{
		Success: false,
		Message: message,
		Error: &ErrorBody{
			Code:       apperror.CodeValidation,
			Message:    message,
			Validation: fields,
		},
	})
}

func requestID(c *gin.Context) string {
	if id, ok := appctx.RequestID(c.Request.Context()); ok {
		return id
	}
	return ""
}
