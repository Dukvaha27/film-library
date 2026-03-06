package handler

import (
	"net/http"

	"github.com/Dukvaha27/film-library/internal/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)



func (h Handler) RegisterUser(c *gin.Context) {
	var req models.AuthenticationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Заполните данные для логина и пароля",
			"error":   err.Error(),
		})
		return
	}

	// чисто для себя полезно будет задуматься о том, почему password это []byte
	password, errPassword := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if errPassword != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Ошибка при генерации пароля",
			"error":   errPassword.Error(),
		})

		return
	}

	user := models.User{
		Username: req.Username,
		Password: string(password),
	}

	if errDB := h.db.Create(&user).Error; errDB != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Ошибка при добавлении пользователя в базу данных",
			"error":   errDB.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":       user.ID,
		"Username": user.Username,
	})
}

func  (h Handler) Login(c *gin.Context) {
	var req models.AuthenticationRequest
	var user models.User

	if errReq := c.ShouldBindJSON(&req); errReq != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Заполните данные логина и пароля",
			"error":   errReq.Error(),
		})
		return
	}

	if errFind := h.db.Where("username = ?", req.Username).First(&user).Error; errFind != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Not found User",
			"error":   errFind.Error(),
		})
		return
	}

	if errCompare := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); errCompare != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Ошибка при проверке пароля",
			"error":   errCompare.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, "Правильный пароль")
}
