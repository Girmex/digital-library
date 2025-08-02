package routes

import (
	"github.com/gin-gonic/gin"
	"digital-library/backend/pkg/email"
	"gorm.io/gorm"
)

func SetupRoutes(db *gorm.DB, emailService *email.Service) *gin.Engine {
	r := gin.Default()

	// Setup auth routes with email service
	SetupAuthRoutes(r, db, emailService)
	
	// Setup other routes without email service
	SetupBookRoutes(r, db)
	SetupLoanRoutes(r, db)
	return r
}