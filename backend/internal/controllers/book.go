package controllers
import (
	"net/http"
	"digital-library/backend/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

)

type BookController struct {
	DB *gorm.DB
}
func (bc *BookController) CreateBook(c *gin.Context) {
	var book models.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := bc.DB.Create(&book).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not create book"})
		return
	}

	c.JSON(http.StatusCreated, book)
}

func (bc *BookController) GetBooks(c *gin.Context) {
	var books []models.Book
	if err := bc.DB.Preload("Category").Find(&books).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch books"})
		return
	}

	c.JSON(http.StatusOK, books)
}


func (bc *BookController) GetBook(c *gin.Context) {
	var book models.Book
	bookID := c.Param("id")
	if err := bc.DB.Preload("Category").First(&book, bookID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}
	c.JSON(http.StatusOK, book)
}

func (bc *BookController) UpdateBook(c *gin.Context) {
	var book models.Book
	bookID := c.Param("id")
	if err := bc.DB.First(&book, bookID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := bc.DB.Save(&book).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update book"})
		return
	}
	c.JSON(http.StatusOK, book)
}
func (bc *BookController) DeleteBook(c *gin.Context) {
	var book models.Book
	bookID := c.Param("id")
	if err := bc.DB.First(&book, bookID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}
	if err := bc.DB.Delete(&book).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete book"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Book deleted successfully"})
}

