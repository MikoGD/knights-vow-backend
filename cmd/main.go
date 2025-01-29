package main

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"

	"knights-vow/internal/database"
	"knights-vow/internal/middleware"
	"knights-vow/internal/resources/files"
	"knights-vow/internal/resources/users"
)

func main() {
	database.InitDatabase()

	defer database.CloseDatabase()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://192.168.68.100:5173", "http://192.168.68.103:5173"}, // Specify exact origin
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           1 * time.Hour,
	}))

	r.Use(middleware.AuthenticateUser)

	v1 := r.Group("api/v1")

	users.CreateRouterGroup(v1)
	files.CreateRouterGroup(v1)

	r.Run("0.0.0.0:8080")
}
