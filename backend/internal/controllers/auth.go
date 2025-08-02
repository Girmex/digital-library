package controllers

import (
	"net/http"
	"time"
	"digital-library/backend/internal/models"
	"digital-library/backend/internal/utils"
	"digital-library/backend/pkg/email"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"errors"
	"os"
	"path/filepath"
	"mime/multipart"


)

type AuthController struct {
	DB    *gorm.DB
	Email *email.Service
}

func (ac *AuthController) Register(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=8"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if email already exists
	var existingUser models.User
	if err := ac.DB.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
		return
	}

	// Generate verification token
	token, err := utils.GenerateVerificationToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate verification token"})
		return
	}

	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user := models.User{
		Username:     input.Username,
		Email:        input.Email,
		Password:     hashedPassword,
		IsVerified:   false,
		VerifyToken:  token,
		VerifyExpiry: time.Now().Add(24 * time.Hour),
	}

	// Assign default role
	var userRole models.Role
	if err := ac.DB.Where("name = ?", "user").First(&userRole).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Default role not found"})
		return
	}
	user.Roles = append(user.Roles, userRole)

	if err := ac.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	// Send verification email
	go func() {
		if err := ac.Email.SendVerificationEmail(user.Email, token); err != nil {
			log.Printf("Failed to send verification email: %v", err)
		}
	}()

	c.JSON(http.StatusCreated, gin.H{
		"message": "Registration successful. Please check your email to verify your account.",
	})
}

func (ac *AuthController) VerifyEmail(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Verification token is required"})
		return
	}

	var user models.User
	if err := ac.DB.Where("verify_token = ? AND verify_expiry > ?", token, time.Now()).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired verification token"})
		return
	}

	// Mark user as verified
	user.IsVerified = true
	user.VerifyToken = ""
	if err := ac.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email verified successfully"})
}

func (ac *AuthController) Login(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := ac.DB.Where("username = ?", input.Username).Preload("Roles").First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if !user.IsVerified {
		c.JSON(http.StatusForbidden, gin.H{"error": "Please verify your email first"})
		return
	}

	if !utils.CheckPasswordHash(input.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Get role names
	var roleNames []string
	for _, role := range user.Roles {
		roleNames = append(roleNames, role.Name)
	}

	token, err := utils.GenerateJWT(user.ID, roleNames)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// ForgotPassword handles password reset requests
func (ac *AuthController) ForgotPassword(c *gin.Context) {
	var input struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find user by email
	var user models.User
	if err := ac.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Return success even if email not found to prevent email enumeration
			c.JSON(http.StatusOK, gin.H{"message": "If this email exists, a reset link has been sent"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Generate reset token
	token, err := utils.GenerateVerificationToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate reset token"})
		return
	}

	// Set reset token and expiry (1 hour)
	user.ResetToken = token
	user.ResetExpiry = time.Now().Add(1 * time.Hour)
	if err := ac.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save reset token"})
		return
	}

	// Send reset email
	go func() {
		if err := ac.Email.SendResetPasswordEmail(user.Email, token); err != nil {
			log.Printf("Failed to send verification email: %v", err)
		}
	}()

		
}

// ResetPassword handles password reset
func (ac *AuthController) ResetPassword(c *gin.Context) {
	var input struct {
		Token       string `json:"token" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=8"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find user by valid reset token
	var user models.User
	if err := ac.DB.Where("reset_token = ? AND reset_expiry > ?", input.Token, time.Now()).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired reset token"})
		return
	}

	// Hash new password
	hashedPassword, err := utils.HashPassword(input.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Update password and clear reset token
	user.Password = hashedPassword
	user.ResetToken = ""
	if err := ac.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset successfully"})
}

// Add to AuthController struct
const uploadDir = "./uploads/profile_photos"

// UploadProfilePhoto handles profile photo uploads
func (ac *AuthController) UploadProfilePhoto(c *gin.Context) {
    // 1. Authentication
    userID, exists := c.Get("userID")
    if !exists {
        c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    // 2. File Validation
    file, err := c.FormFile("profile_photo")
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Valid image file required (JPEG/PNG)"})
        return
    }

    // 3. File Type Check
    if !isValidImageType(file) {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Only JPEG/PNG allowed"})
        return
    }

    // 4. Process Upload
    filePath, err := utils.SaveUploadedFile(file, uploadDir)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "File upload failed"})
        return
    }

    // 5. Update Database
    var user models.User
    if err := ac.DB.First(&user, userID).Error; err != nil {
        os.Remove(filePath) // Clean up new file if user not found
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    // 6. Delete Old Photo (if exists)
    if user.ProfilePhoto != "" {
        if err := os.Remove(user.ProfilePhoto); err != nil {
            log.Printf("Warning: Failed to delete old photo: %v", err)
        }
    }

    // 7. Save Changes
    user.ProfilePhoto = filePath
    if err := ac.DB.Save(&user).Error; err != nil {
        os.Remove(filePath) // Clean up if save fails
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
        return
    }

    // 8. Return Response
    c.JSON(http.StatusOK, gin.H{
        "message":  "Profile photo updated",
        "photo_url": "/profile-photos/" + filepath.Base(filePath),
    })
}

// Helper function
func isValidImageType(file *multipart.FileHeader) bool {
    allowedTypes := map[string]bool{
        "image/jpeg": true,
        "image/png":  true,
    }
    return allowedTypes[file.Header.Get("Content-Type")]
}
// GetProfilePhoto serves the profile photo
func (ac *AuthController) GetProfilePhoto(c *gin.Context) {
    userID := c.Param("userID")
    
    var user models.User
    if err := ac.DB.Select("profile_photo").First(&user, userID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    if user.ProfilePhoto == "" {
        c.JSON(http.StatusNotFound, gin.H{"error": "No profile photo"})
        return
    }

    c.File(user.ProfilePhoto)
}