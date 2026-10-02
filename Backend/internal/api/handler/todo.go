package handler

import "github.com/gin-gonic/gin"

type Todo struct{}

func NewTodo() *Todo{
	return &Todo{}
}

//read

//create
func (t *Todo)CreatTodo(c *gin.Context){}

//update

//delete
