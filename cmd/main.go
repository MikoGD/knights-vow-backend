package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"

	"knights-vow/internal/database"
	"knights-vow/internal/middleware"
	"knights-vow/internal/users"
	// "knights-vow/uploads"
)

func main() {
	database.InitDatabase()

	defer database.CloseDatabase()

	r := gin.Default()

	r.Use(middleware.AuthenticateUser)

	v1 := r.Group("api/v1")

	users.CreateRouterGroup(v1)

	r.Run()
}
