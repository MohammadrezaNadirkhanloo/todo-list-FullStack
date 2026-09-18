package handlers

import (
	"net/http"

	"github.com/MohammadrezaNadirkhanloo/todo-list-FullStack/internal/repository"
	"github.com/gin-gonic/gin"
)

type CreateTodoInput struct {
	Title     string `json:"title" binding:"required"`
	Completed bool   `json:"completed"`
}

func CreateTodoHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// userIDInterface, exists := c.Get("user_id")

		// if !exists {
		// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "user_id not found in context"})
		// 	return
		// }

		// userID := userIDInterface.(string)

		var input CreateTodoInput

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		todo, err := repository.CreateTodo(input.Title, input.Completed)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return 
		}

		c.JSON(http.StatusCreated, todo)
	}
}
