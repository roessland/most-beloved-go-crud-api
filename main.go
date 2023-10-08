package main

import (
	"context"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/roessland/most-beloved-go-crud-api/handlers"
)

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/", handlers.Health)
	r.POST("/", handlers.Create)

	return r
}

func main() {
	ctx := context.Background()

	conn, err := pgx.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}
	defer conn.Close(ctx)

	dbHealthCheck(conn)

	r := setupRouter()
	r.Run()
}

func dbHealthCheck(conn *pgx.Conn) {
	healthCheckCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var result int
	err := conn.QueryRow(healthCheckCtx, "SELECT 1 + 1;").Scan(&result)
	if err != nil || result != 2 {
		panic(err)
	}
}
