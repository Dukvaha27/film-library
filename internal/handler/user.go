package handler

import (
	"net/http"
	"strconv"

	"github.com/Dukvaha27/film-library/internal/models"
	"github.com/gin-gonic/gin"
)

func (h Handler) GetUserInfo(c *gin.Context) {
	var user models.User
	userID, errID := strconv.Atoi(c.Param("user_id"))
	if errID != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "ошибка преобразования строки в число",
			"error":   errID.Error(),
		})
		return
	}
	if err := h.db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Ошибка при поиске юзера, похоже его нет!",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, user)

}
