package api

import (
	"net/http"
	"roomsync/models"
	"roomsync/repository"
	"roomsync/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type RegisterReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
	Email    string `json:"email" binding:"required,email"`
}

type LoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Register(c *gin.Context) {
	var req RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	// Check if username or email exists
	var count int64
	repository.DB.Model(&models.User{}).Where("username = ? OR email = ?", req.Username, req.Email).Count(&count)
	if count > 0 {
		utils.Error(c, http.StatusBadRequest, "Username or email already exists")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	user := models.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Email:    req.Email,
		Role:     "employee", // Default role
	}

	// First user becomes admin automatically for testing purposes
	var totalUsers int64
	repository.DB.Model(&models.User{}).Count(&totalUsers)
	if totalUsers == 0 {
		user.Role = "admin"
	}

	if err := repository.DB.Create(&user).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to create user")
		return
	}

	utils.Success(c, user)
}

func Login(c *gin.Context) {
	var req LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	var user models.User
	if err := repository.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		utils.Error(c, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		utils.Error(c, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	token, err := utils.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	utils.Success(c, gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"role":     user.Role,
		},
	})
}

func GetCurrentUser(c *gin.Context) {
	userID, _ := c.Get("userID")
	var user models.User
	if err := repository.DB.First(&user, userID).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "User not found")
		return
	}
	utils.Success(c, user)
}
