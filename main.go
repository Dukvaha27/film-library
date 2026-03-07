package main

import (
	"fmt"
	"os"

	"github.com/Dukvaha27/film-library/internal/handler"
	"github.com/Dukvaha27/film-library/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	r := gin.Default()
err := godotenv.Load()
if err != nil {
		fmt.Println("Ошибка загрузки .env файла")
		return
	}
	// вызови функцию Load из библиотеки godotenv

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"))

	db, errDB := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if errDB != nil {
		fmt.Printf("Ошибка при инициализации DB %v", errDB)
		return
	}
	if errMigrating := db.AutoMigrate(&models.User{}, &models.Movie{}, &models.Watchlist{}); errMigrating != nil {
		fmt.Printf("Ошибка при миграции таблиц %v", errMigrating)
		return
	}

	h := handler.New(db)

	r.POST("/register",h.RegisterUser) //
	r.POST("/login", h.Login)
	r.GET("/users/:user_id",h.GetUserInfo)

	group := r.Group("watchlist", h.GetUserWatchList)
	group.POST("/:user_id/:movie_id",h.AddToWatchlist)
	group.DELETE("/:user_id/:movie_id", h.DeleteFromWatchlist)

	r.Run(":3000")
}
