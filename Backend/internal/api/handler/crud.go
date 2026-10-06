package handler

import (
	"context"

	"github.com/MohammadrezaNadirkhanloo/internal/api/response"
	"github.com/MohammadrezaNadirkhanloo/internal/domain/authz"
	"github.com/MohammadrezaNadirkhanloo/internal/domain/filter"
	"github.com/gin-gonic/gin"
)

type CRUDUsecase[TCreate any, TUpdate any, TOut any] interface {
	Create(ctx context.Context, in TCreate) (TOut, error)
	Update(ctx context.Context, id int64, in TUpdate) (TOut, error)
	Delete(ctx context.Context, id int64) error
	GetByID(ctx context.Context, id int64) (TOut, error)
	GetByQuery(ctx context.Context, q filter.Query) (*filter.PagedList[TOut], error)

	Spec() filter.Spec
	PageLimits() filter.PageLimits
}

type ListMode int

const (
	ListSimple ListMode = iota

	ListAdvanced
)

type CRUD[TCreate any, TUpdate any, TOut any] struct {
	usecase  CRUDUsecase[TCreate, TUpdate, TOut]
	basePath string
	mode     ListMode
}

func NewCRUD[TCreate any, TUpdate any, TOut any](
	uc CRUDUsecase[TCreate, TUpdate, TOut],
	basePath string,
	mode ListMode,
) *CRUD[TCreate, TUpdate, TOut] {
	return &CRUD[TCreate, TUpdate, TOut]{usecase: uc, basePath: basePath, mode: mode}
}

func (h *CRUD[TCreate, TUpdate, TOut]) Create(c *gin.Context) {
	var in TCreate
	if !BindJSON(c, &in) {
		return
	}

	out, err := h.usecase.Create(c.Request.Context(), in)
	if err != nil {
		response.Fail(c, err)
		return
	}

	respondCreated(c, "", out)
}

func (h *CRUD[TCreate, TUpdate, TOut]) Update(c *gin.Context) {
	id, ok := PathID(c)
	if !ok {
		return
	}

	var in TUpdate
	if !BindJSON(c, &in) {
		return
	}

	out, err := h.usecase.Update(c.Request.Context(), id, in)
	if err != nil {
		response.Fail(c, err)
		return
	}

	response.OK(c, out)
}

func (h *CRUD[TCreate, TUpdate, TOut]) Delete(c *gin.Context) {
	id, ok := PathID(c)
	if !ok {
		return
	}

	if err := h.usecase.Delete(c.Request.Context(), id); err != nil {
		response.Fail(c, err)
		return
	}

	response.NoContent(c)
}

func (h *CRUD[TCreate, TUpdate, TOut]) GetByID(c *gin.Context) {
	id, ok := PathID(c)
	if !ok {
		return
	}

	out, err := h.usecase.GetByID(c.Request.Context(), id)
	if err != nil {
		response.Fail(c, err)
		return
	}

	response.OK(c, out)
}

func (h *CRUD[TCreate, TUpdate, TOut]) List(c *gin.Context) {
	q, ok := ParseListQuery(c, h.usecase.Spec(), h.usecase.PageLimits(), h.mode)
	if !ok {
		return
	}

	page, err := h.usecase.GetByQuery(c.Request.Context(), q)
	if err != nil {
		response.Fail(c, err)
		return
	}

	response.List(c, page)
}

func RegisterCRUD[TCreate any, TUpdate any, TOut any](
	rg *gin.RouterGroup,
	h *CRUD[TCreate, TUpdate, TOut],
	subject string,
	authorize func(action, subject string) gin.HandlerFunc,
) {
	if subject == "" {
		panic("handler: RegisterCRUD requires a subject")
	}
	if authorize == nil {
		panic("handler: RegisterCRUD requires an authorize function")
	}

	rg.GET("", h.List)
	rg.GET("/:id", h.GetByID)

	rg.POST("", authorize(authz.ActionCreate, subject), h.Create)
	rg.PATCH("/:id", authorize(authz.ActionUpdate, subject), h.Update)
	rg.DELETE("/:id", authorize(authz.ActionDelete, subject), h.Delete)
}
