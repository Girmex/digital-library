package routes

import (
	"digital-library/backend/internal/controllers"
	"digital-library/backend/internal/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupBookRoutes(r *gin.Engine, db *gorm.DB) {
	bookCtrl := &controllers.BookController{DB: db}

	// All book routes require JWT
	bookRoutes := r.Group("/books")
	bookRoutes.Use(middleware.JWTAuth())
	{
		// Viewingwsl
		
		bookRoutes.GET("/", bookCtrl.GetBooks)  
		bookRoutes.GET("/:id", bookCtrl.GetBook) 

		// Modification (with extra permissions)
		bookRoutes.POST("/", middleware.HasPermission("create_book"), bookCtrl.CreateBook)
		bookRoutes.PUT("/:id", middleware.HasPermission("edit_book"), bookCtrl.UpdateBook)
		bookRoutes.DELETE("/:id", middleware.HasPermission("delete_book"), bookCtrl.DeleteBook)
	}
}