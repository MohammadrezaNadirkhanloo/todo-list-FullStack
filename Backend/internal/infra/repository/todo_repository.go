package repository

import (
	"github.com/MohammadrezaNadirkhanloo/internal/domain/model"
	"github.com/MohammadrezaNadirkhanloo/internal/domain/repository"
	"github.com/MohammadrezaNadirkhanloo/internal/infra/database"
)

type todoRepository struct {
	*BaseRepository[model.Todo]
}

var _ repository.TodoRepository = (*todoRepository)(nil)

func NewTodoRepository(db *database.DB) repository.TodoRepository {
	return &todoRepository{
		BaseRepository: NewBaseRepository[model.Todo](db, Options{
			Entity: "Todo",
			Spec:   model.TodoSpec,
		}),
	}
}
