package handler

import (
	"github.com/MohammadrezaNadirkhanloo/internal/usecase"
	"github.com/MohammadrezaNadirkhanloo/internal/usecase/dto"
)

type TodoHandler struct {
	*CRUD[dto.CreateTodoInput, dto.UpdateTodoInput, dto.TodoOutput]
}

var _ CRUDUsecase[dto.CreateTodoInput, dto.UpdateTodoInput, dto.TodoOutput] = (*usecase.TodoUsecase)(nil)

func NewTodoHandler(uc *usecase.TodoUsecase) *TodoHandler {
	return &TodoHandler{
		CRUD: NewCRUD(uc, "/todos", ListSimple),
	}
}
