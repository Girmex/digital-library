package routes

import (
	"digital-library/backend/internal/controllers"
	"digital-library/backend/pkg/email"
	"digital-library/backend/internal/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupAuthRoutes(r *gin.Engine, db *gorm.DB, emailService *email.Service) {
    authCtrl := &controllers.AuthController{
        DB:    db,
        Email: emailService,
    }

    auth := r.Group("/auth")
    {
        auth.POST("/register", authCtrl.Register)
        auth.GET("/verify-email", authCtrl.VerifyEmail)
        auth.POST("/login", authCtrl.Login)
        auth.POST("/forgot-password", authCtrl.ForgotPassword)
        auth.POST("/reset-password", authCtrl.ResetPassword)
        
        // Protected routes
        auth.Use(middleware.JWTAuth())
        {
            auth.POST("/upload-profile-photo", authCtrl.UploadProfilePhoto)
            auth.GET("/users/:userID/profile-photo", authCtrl.GetProfilePhoto)
        }
    }
}