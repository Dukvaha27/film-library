package main

import (
	"fmt"

	"github.com/Dukvaha27/film-library/internal/database"
	"github.com/Dukvaha27/film-library/internal/handler"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	r := gin.Default()

	if err := godotenv.Load(); err != nil {
		fmt.Println("Ошибка загрузки .env файла")
		return
	}
	db, err := database.Connect()
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := database.Migrate(db); err != nil {
		fmt.Println(err)
		return
	}

	h := handler.New(db)

	r.POST("/register", h.RegisterUser)
	r.POST("/login", h.Login)
	r.GET("/users/:user_id", h.GetUserInfo)

	watchlist := r.Group("/watchlist")
	{
		watchlist.GET("/:user_id", h.GetUserWatchList)
		watchlist.POST("/:user_id/:movie_id", h.AddToWatchlist)
		watchlist.DELETE("/:user_id/:movie_id", h.DeleteFromWatchlist)
	}

	movies := r.Group("/movies")
	{
		movies.POST("/:id/reviews", h.CreateReview)
		movies.GET("/:id/reviews", h.GetMovieReviews)
		movies.GET("/:id", h.GetMovieById)
		movies.POST("", h.CreateMovie)
		movies.PATCH("/:id", h.UpdateMovie)
		movies.DELETE("/:id", h.RemoveMovie)
	}

	genres := r.Group("/genres")
	{
		genres.GET("", h.GetGenres)
		genres.POST("", h.CreateGenre)
		genres.PUT("/:id", h.UpdateGenre)
		genres.DELETE("/:id", h.DeleteGenre)
	}

	r.Run(":3000")
}
