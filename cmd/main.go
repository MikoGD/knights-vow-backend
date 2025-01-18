package main

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"

	"knights-vow/internal/database"
	"knights-vow/internal/middleware"
	"knights-vow/internal/resources/uploads"
	"knights-vow/internal/resources/users"
)

func main() {
	database.InitDatabase()

	defer database.CloseDatabase()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           2 * time.Hour,
	}))
	r.Use(middleware.AuthenticateUser)

	v1 := r.Group("api/v1")

	users.CreateRouterGroup(v1)
	uploads.CreateRouterGroup(v1)

	r.Run()
}
