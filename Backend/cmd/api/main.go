package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/MohammadrezaNadirkhanloo/todo-list-FullStack/internal/config"
	"github.com/MohammadrezaNadirkhanloo/todo-list-FullStack/internal/database"
	"github.com/MohammadrezaNadirkhanloo/todo-list-FullStack/internal/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	if err := database.InitDB(); err != nil {
		log.Fatal(err)
	}
	defer database.CloseDB()

	router := gin.Default()
	router.SetTrustedProxies(nil)
	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "hello",
		})
	})

	router.POST("/todos", handlers.CreateTodoHandler())

	router.Run(fmt.Sprintf(":%s", config.Config.Port))
}
