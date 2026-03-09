package handler

import (
	"net/http"
	"strconv"

	"github.com/Dukvaha27/film-library/internal/models"
	"github.com/gin-gonic/gin"
)

func (h Handler) CreateMovie(ctx *gin.Context) {
	var movie models.Movie

	err := ctx.ShouldBindJSON(&movie)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	h.db.Create(&movie)
	ctx.JSON(http.StatusOK, movie)
}

func (h Handler) GetMovies(ctx *gin.Context) {
	query := h.db.Model(&models.Movie{})

	genre_id := ctx.Query("genre_id")
	year := ctx.Query("year")

	if genre_id != "" {
		if gId, gIdErr := strconv.Atoi(genre_id); gIdErr == nil {
			query = query.Where("genre_id = ? ", gId)
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": gIdErr.Error()})
		}
	}

	if year != "" {
		if yId, yIdErr := strconv.Atoi(genre_id); yIdErr == nil {
			query = query.Where("year = ? ", yId)
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": yIdErr.Error()})
		}
	}

	movies := []models.Movie{}
	result := query.Find(&movies)

	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	if result.RowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Фильмы не найдены",
		})
	}

	ctx.JSON(http.StatusOK, movies)
}

func (h Handler) GetMovieById(ctx *gin.Context) {
	id, errId := strconv.Atoi(ctx.Param("id"))

	if errId != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": errId.Error(), "value": "ID nust be string"})
	}

	var movie models.Movie

	result := h.db.First(&movie, id)

	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	if result.RowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "фильм не найден"})
		return
	}

	ctx.JSON(http.StatusOK, movie)
}

func (h Handler) UpdateMovie(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID must be number"})
		return
	}

	var updateMovie models.UpdateMovie
	if err := ctx.ShouldBindJSON(&updateMovie); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var updatedMovie models.Movie
	result := h.db.Model(&models.Movie{}).Where("id = ?", id).Updates(updateMovie)

	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	if result.RowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Фильм не найден"})
		return
	}

	h.db.First(&updatedMovie, id)

	ctx.JSON(http.StatusOK, updatedMovie)
}

func (h Handler) RemoveMovie(ctx *gin.Context) {
	id, errIdd := strconv.Atoi(ctx.Param("id"))

	if errIdd != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "ID must be number"})
		return
	}

	result := h.db.Delete(&models.Movie{}, id)

	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": " an occured when deleting"})
		return
	}

	if result.RowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "movie is not found"})
		return
	}

	ctx.JSON(http.StatusOK, "record was deleted")
}
