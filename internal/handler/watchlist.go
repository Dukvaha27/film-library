package handler

import (
	"net/http"
	"strconv"

	"github.com/Dukvaha27/film-library/internal/models"
	"github.com/gin-gonic/gin"
)

// Watchlist (POST/DELETE /watchlist, GET /users/:id/watchlist)

// POST   /watchlist/:movied_id/:user_id
// DELETE /watchlist/:mob
func (h Handler) AddToWatchlist(c *gin.Context) {
	userID, errUserID := strconv.Atoi(c.Param("user_id"))
	if errUserID != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Введите число а не строку",
		})
		return
	}

	movieID, errMovieID := strconv.Atoi(c.Param("movie_id"))
	if errMovieID != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Введите число а не строку",
		})
		return
	}

	newWatchlist := models.Watchlist{
		MovieID: uint(movieID),
		UserID:  uint(userID),
	}

	if errDB := h.db.Create(&newWatchlist).Error; errDB != nil {
		c.JSON(http.StatusBadRequest, errDB.Error())
		return
	}

	c.JSON(http.StatusCreated, "Watchlist added")
}

func (h Handler) DeleteFromWatchlist(c *gin.Context) {
	userID, errUserID := strconv.Atoi(c.Param("user_id"))
	if errUserID != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Введите число а не строку",
		})
		return
	}

	movieID, errMovieID := strconv.Atoi(c.Param("movie_id"))
	if errMovieID != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Введите число а не строку",
		})
		return
	}

	if err := h.db.Exec("DELETE FROM watchlist where user_id =$1 and movie_id=$2", userID, movieID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (h Handler) GetUserWatchList(c *gin.Context) {
	var watchlist []models.Watchlist
	userID, errType := strconv.Atoi(c.Param("id"))
	if errType != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errType.Error(),
		})
		return
	}
	if errDB := h.db.Where("user_id = ?", userID).Preload("Movie").Find(&watchlist).Error; errDB != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": errDB.Error(),
		})
		return
	}
	movies := []models.Movie{}
	for _, v := range watchlist {
		movies = append(movies, v.Movie)
	}

	c.JSON(http.StatusOK, movies)
}
