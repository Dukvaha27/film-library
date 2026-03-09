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

type CreateGenreRequest struct {
	Name string `json:"name" binding:"required"`
}

type UpdateGenreRequest struct {
	Name string `json:"name" binding:"required"`
}

func (h *Handler) CreateGenre(c *gin.Context) {
	var req CreateGenreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "genre name is required"})
		return
	}

	var existingGenre models.Genre
	err := h.db.Where("LOWER(name) = LOWER(?)", req.Name).First(&existingGenre).Error

	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "genre already exists"})
		return
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check genre"})
		return
	}

	genre := models.Genre{Name: req.Name}
	if err := h.db.Create(&genre).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create genre"})
		return
	}

	c.JSON(http.StatusCreated, genre)
}

func (h *Handler) GetGenres(c *gin.Context) {
	var genres []models.Genre

	if err := h.db.Order("name ASC").Find(&genres).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get genres"})
		return
	}

	c.JSON(http.StatusOK, genres)
}

func (h *Handler) UpdateGenre(c *gin.Context) {
	genreIDParam := c.Param("id")
	genreID, err := strconv.ParseUint(genreIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid genre id"})
		return
	}

	var req UpdateGenreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "genre name is required"})
		return
	}

	var genre models.Genre
	if err := h.db.First(&genre, uint(genreID)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "genre not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get genre"})
		return
	}

	var existingGenre models.Genre
	err = h.db.Where("LOWER(name) = LOWER(?) AND id <> ?", req.Name, uint(genreID)).First(&existingGenre).Error
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "genre already exists"})
		return
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check genre"})
		return
	}

	genre.Name = req.Name
	if err := h.db.Save(&genre).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update genre"})
		return
	}

	c.JSON(http.StatusOK, genre)
}

func (h *Handler) DeleteGenre(c *gin.Context) {
	genreIDParam := c.Param("id")
	genreID, err := strconv.ParseUint(genreIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid genre id"})
		return
	}

	var genre models.Genre
	if err := h.db.First(&genre, uint(genreID)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "genre not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get genre"})
		return
	}

	err = h.db.Transaction(func(tx *gorm.DB) error {
		if err := h.db.
			Where("genre_id = ?", genre.ID).
			Delete(&models.Movie{}).Error; err != nil {
			return err
		}

		if err := h.db.Delete(&genre).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "genre deleted"})
}
