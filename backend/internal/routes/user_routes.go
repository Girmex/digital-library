package routes

import (
	"digital-library/backend/internal/controllers"
	"digital-library/backend/internal/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupUserRoutes(r *gin.Engine, db *gorm.DB) {
	userCtrl := &controllers.UserController{DB: db}
	
	// All user routes require JWT
	userRoutes := r.Group("/users")
	userRoutes.Use(middleware.JWTAuth())
	{	
		userRoutes.GET("/", userCtrl.GetUsers)
		userRoutes.GET("/:id", userCtrl.GetUser)
		userRoutes.PUT("/:id", middleware.HasPermission("edit_user"), userCtrl.UpdateUser)
		userRoutes.DELETE("/:id", middleware.HasPermission("delete_user"), userCtrl.DeleteUser)
	}
}