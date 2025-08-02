package routes

import (
	"digital-library/backend/internal/controllers"
	"digital-library/backend/internal/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupLoanRoutes(r *gin.Engine, db *gorm.DB) {
	loanCtrl := &controllers.LoanController{DB: db}

	// All loan routes require JWT
	loanRoutes := r.Group("/loans")
	loanRoutes.Use(middleware.JWTAuth())
	{
		loanRoutes.POST("/", middleware.HasPermission("checkout_book"), loanCtrl.CheckoutBook)
		loanRoutes.PUT("/:id/return", middleware.HasPermission("return_book"), loanCtrl.ReturnBook)
		loanRoutes.GET("/user/:user_id", loanCtrl.GetUserLoans) // Requires auth
	}
}