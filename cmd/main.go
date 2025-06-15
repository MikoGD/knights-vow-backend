package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"www.github.com/mikogd/hextech/env"

	"knights-vow/internal/database"
	"knights-vow/internal/middleware"
	"knights-vow/internal/resources/files"
	"knights-vow/internal/resources/users"
)

func main() {
	if err := env.LoadEnv("./.env"); err != nil {
		log.Fatalf("Failed to load env: %s\n", err)
	}

	db := database.InitDatabase()
	defer database.CloseDatabase(db)

	r := gin.Default()

	if err := env.LoadEnv("./.env"); err != nil {
		log.Fatalf("Failed to load env\n%s", err)
	}

	// allowedOriginsValue := os.Getenv("ALLOWED_ORIGINS")
	// if allowedOriginsValue == "" {
	// 	log.Fatalln("ALLOWED_ORIGINS not set")
	// }
	// allowedOrigins := strings.Split(allowedOriginsValue, ",")

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           1 * time.Hour,
	}))

	r.Use(middleware.AuthenticateUser)

	v1 := r.Group("api/v1")

	users.CreateRouterGroup(v1, db)
	files.CreateRouterGroup(v1, db)

	URL := os.Getenv("URL")
	port := os.Getenv("PORT")
	if URL == "" || port == "" {
		log.Fatalln("URL or PORT not set")
	}

	r.Run(fmt.Sprintf("%s:%s", URL, port))
	// r.Run("0.0.0.0:8080")
}
