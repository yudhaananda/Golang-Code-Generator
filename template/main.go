package main

import (
	"fmt"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"[project]/entity"
	"[project]/handler"
	"[project]/repository"
	"[project]/service"
	"log"
)

func main() {

	dsn := "root:@tcp(127.0.0.1:3306)/[project]?charset=utf8mb4&parseTime=true&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err.Error())
	}
	db.AutoMigrate([migrateArea])
	//Repository Region
	[repoArea]

	//Service Region
	[serviceArea]

	//Handler Region
	[handlerArea]

	//Router Region
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"PUT", "PATCH", "GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "https://github.com"
		},
		MaxAge: 12 * time.Hour,
	}))

	router.GET("/", testHandler.TestAPI)
	api := router.Group("/api/v1")
	//auth
	[apiArea]

	router.Run()

}
