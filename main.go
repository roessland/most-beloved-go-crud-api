package main

import (
	"github.com/gin-gonic/gin"
	"github.com/roessland/most-beloved-go-crud-api/handlers"
)

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/", handlers.Health)
	r.POST("/", handlers.Create)

	return r
}

func main() {
	r := setupRouter()
	// Listen on PORT env var, or on 3000 if not set
	r.Run()
}
