package controllers

import (
	"net/http"
	"time"
	"digital-library/backend/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type LoanController struct {
	DB *gorm.DB
}

// Checkout a book
func (lc *LoanController) CheckoutBook(c *gin.Context) {
	var input struct {
		BookID uint   `json:"book_id" binding:"required"`
		UserID uint   `json:"user_id" binding:"required"`
		Days   int    `json:"loan_days" binding:"required"` // Loan duration
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify book exists and is available
	var book models.Book
	if err := lc.DB.First(&book, input.BookID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	if book.Status != "AVAILABLE" {
		c.JSON(http.StatusConflict, gin.H{"error": "Book is not available"})
		return
	}

	// Create loan
	loan := models.Loan{
		UserID:      input.UserID,
		BookID:      input.BookID,
		CheckoutDate: time.Now(),
		DueDate:     time.Now().AddDate(0, 0, input.Days),
		Status:      "ACTIVE",
	}

	if err := lc.DB.Create(&loan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create loan"})
		return
	}

	// Update book status
	lc.DB.Model(&book).Update("status", "CHECKED_OUT")

	c.JSON(http.StatusCreated, loan)
}

// Return a book
func (lc *LoanController) ReturnBook(c *gin.Context) {
	loanID := c.Param("id")

	var loan models.Loan
	if err := lc.DB.Preload("Book").First(&loan, loanID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Loan not found"})
		return
	}

	// Update loan
	returnTime := time.Now()
	loan.ReturnDate = &returnTime
	loan.Status = "RETURNED"

	if err := lc.DB.Save(&loan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update loan"})
		return
	}

	// Update book status
	lc.DB.Model(&loan.Book).Update("status", "AVAILABLE")

	c.JSON(http.StatusOK, loan)
}

// Get user's active loans
func (lc *LoanController) GetUserLoans(c *gin.Context) {
	userID := c.Param("user_id")

	var loans []models.Loan
	if err := lc.DB.Preload("Book").
		Where("user_id = ? AND status = ?", userID, "ACTIVE").
		Find(&loans).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch loans"})
		return
	}

	c.JSON(http.StatusOK, loans)
}

// Get overdue loans
func (lc *LoanController) GetOverdueLoans(c *gin.Context) {
	var loans []models.Loan
	if err := lc.DB.Preload("Book").Preload("User").
		Where("due_date < ? AND status = ?", time.Now(), "ACTIVE").
		Find(&loans).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch overdue loans"})
		return
	}

	c.JSON(http.StatusOK, loans)
}