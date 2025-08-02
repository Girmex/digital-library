package controllers

import (
	"net/http"
	"digital-library/backend/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

)
type UserController struct {
	DB *gorm.DB
}
func (uc *UserController) GetUsers(c *gin.Context) {
	var users []models.User
	if err := uc.DB.Preload("Roles").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch users"})
		return
	}
	c.JSON(http.StatusOK, users)
	
}
func (uc *UserController) GetUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	if err := uc.DB.Preload("Roles").First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (uc *UserController) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	if err := uc.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := uc.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update user"})
		return
	}
	c.JSON(http.StatusOK, user)
} 

func (uc *UserController) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	if err := uc.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	if err := uc.DB.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}



