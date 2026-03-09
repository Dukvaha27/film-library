package handler

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/Dukvaha27/film-library/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CreateReviewRequest struct {
	Author  string `json:"author" binding:"required"`
	Rating  int    `json:"rating" binding:"required"`
	Comment string `json:"comment"`
}

func (h *Handler) CreateReview(c *gin.Context) {
	movieIDParam := c.Param("id")
	movieID, err := strconv.ParseUint(movieIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid movie id"})
		return
	}

	var req CreateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	req.Author = strings.TrimSpace(req.Author)
	req.Comment = strings.TrimSpace(req.Comment)

	if req.Author == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "author is required"})
		return
	}

	if req.Rating < 1 || req.Rating > 10 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "rating must be between 1 and 10"})
		return
	}

	var review models.Review

	err = h.db.Transaction(func(tx *gorm.DB) error {
		var movie models.Movie
		if err := tx.First(&movie, uint(movieID)).Error; err != nil {
			return err
		}

		review = models.Review{
			MovieID: uint(movieID),
			Author:  req.Author,
			Rating:  req.Rating,
			Comment: req.Comment,
		}

		if err := tx.Create(&review).Error; err != nil {
			return err
		}

		var avgRating float64
		if err := tx.Model(&models.Review{}).
			Where("movie_id = ?", uint(movieID)).
			Select("AVG(rating)").
			Scan(&avgRating).Error; err != nil {
			return err
		}

		if err := tx.Model(&movie).
			Update("avg_rating", avgRating).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "movie not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create review"})
		return
	}

	c.JSON(http.StatusCreated, review)
}

func (h *Handler) GetMovieReviews(c *gin.Context) {
	movieIDParam := c.Param("id")
	movieID, err := strconv.ParseUint(movieIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid movie id"})
		return
	}

	var movie models.Movie
	if err := h.db.First(&movie, uint(movieID)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "movie not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get movie"})
		return
	}

	var reviews []models.Review
	if err := h.db.
		Where("movie_id = ?", uint(movieID)).
		Order("created_at desc").
		Find(&reviews).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get reviews"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"movie_id": movie.ID,
		"reviews":  reviews,
	})
}
