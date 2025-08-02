package main

import (
	"digital-library/backend/internal/models"
	"digital-library/backend/internal/routes"
	"digital-library/backend/internal/utils"
	"digital-library/backend/pkg/database"
	"digital-library/backend/pkg/email"
	"log"
	"time"
)

func main() {
	// Initialize database
	database.ConnectDB()
	
	// Initialize email service
	emailService := email.NewService(email.Config{
		From:     "your gmail",
		Password: "your app password",
		SmtpHost: "smtp.gmail.com",
		SmtpPort: "587",
	})
	// Seed initial data
	seedData(emailService)
	
	// Set up routes with email service
	r := routes.SetupRoutes(database.DB, emailService)
	r.Static("/profile-photos", "./uploads/profile_photos")

	
	// Start server
	log.Println("Server starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}

func seedData(emailService *email.Service) {
	// Create permissions
	permissions := []string{
		"create_book",
		"edit_book",
		"delete_book",
		"manage_users",
		"view_books",
		"checkout_book",
		"return_book",
		"view_loans",
		"manage_overdue",
		"verify_email",
	}
	
	for _, p := range permissions {
		database.DB.FirstOrCreate(&models.Permission{Name: p}, models.Permission{Name: p})
	}
	
	// Create roles
	adminRole := models.Role{Name: "admin"}
	database.DB.FirstOrCreate(&adminRole, models.Role{Name: "admin"})
	
	// Assign all permissions to admin
	var allPermissions []models.Permission
	database.DB.Find(&allPermissions)
	database.DB.Model(&adminRole).Association("Permissions").Append(allPermissions)
	
	userRole := models.Role{Name: "user"}
	database.DB.FirstOrCreate(&userRole, models.Role{Name: "user"})
	
	// Assign permissions to user
	basicPermissions := []models.Permission{}
	database.DB.Where("name IN ?", []string{
		"view_books",
		"checkout_book",
		"return_book",
		"view_loans",
	}).Find(&basicPermissions)
	database.DB.Model(&userRole).Association("Permissions").Append(basicPermissions)
	
	// Create initial admin user (pre-verified)
	hashedPassword, err := utils.HashPassword("admin123")
	if err != nil {
		log.Println("Failed to hash admin password:", err)
		return
	}
	
	adminUser := models.User{
		Username:    "admin",
		Email:       "admin@library.com",
		Password:    hashedPassword,
		IsVerified:  true,
		VerifyToken: "",
	}
	
	if err := database.DB.FirstOrCreate(&adminUser, models.User{Username: "admin"}).Error; err != nil {
		log.Println("Failed to create admin user:", err)
		return
	}
	
	// Assign admin role
	if err := database.DB.Model(&adminUser).Association("Roles").Append(&adminRole); err != nil {
		log.Println("Failed to assign admin role:", err)
	}
	
	// Create test user (unverified initially)
	userPassword, _ := utils.HashPassword("user123")
	token, err := utils.GenerateVerificationToken()
	if err != nil {
		log.Println("Failed to generate verification token:", err)
		return
	}
	
	testUser := models.User{
		Username:     "testuser",
		Email:        "user@library.com",
		Password:     userPassword,
		IsVerified:   false,
		VerifyToken:  token,
		VerifyExpiry: time.Now().Add(24 * time.Hour),
	}
	
	if err := database.DB.FirstOrCreate(&testUser, models.User{Username: "testuser"}).Error; err == nil {
		// Assign user role
		database.DB.Model(&testUser).Association("Roles").Append(&userRole)
		
		// Send verification email for test user
		go func() {
			if err := emailService.SendVerificationEmail(testUser.Email, testUser.VerifyToken); err != nil {
				log.Printf("Failed to send test user verification: %v", err)
			}
		}()
	}
	
	log.Println("Database seeding completed successfully")
}