package models

import (
	"time"

	"gorm.io/gorm"
)
type User struct {
    gorm.Model
    Username       string    `gorm:"unique;not null"`
    Email          string    `gorm:"unique;not null"`
    Password       string    `gorm:"not null"`
    ProfilePhoto   string    // Stores the file path or URL
    IsVerified     bool      `gorm:"default:false"`
    VerifyToken    string    
    VerifyExpiry   time.Time 
    ResetToken     string    
    ResetExpiry    time.Time 
    Roles          []Role    `gorm:"many2many:user_roles;"`
}

type Role struct {
	gorm.Model
	Name string `gorm:"unique;not null"`
	Permissions []Permission `gorm:"many2many:role_permissions;"`

}

type Permission struct {
	gorm.Model
	Name string `gorm:"unique;not null"`

}
type Book struct {
	gorm.Model
	Title string `gorm:"unique;not null"`
	Author string `gorm:"not null"`
	ISBN string `gorm:"unique;not null"`
	Description string `gorm:"not null"`
	CategoryID uint
	Category Category
	Status string // Available, CheckedOut, Reserved	
}
type Category struct {
	gorm.Model
	Name string `gorm:"unique;not null"`
}
type Loan struct {
    gorm.Model
    UserID      uint      `gorm:"not null"`
    User        User      `gorm:"foreignKey:UserID"`
    BookID      uint      `gorm:"not null"`
    Book        Book      `gorm:"foreignKey:BookID"`
    CheckoutDate time.Time `gorm:"not null"`
    DueDate     time.Time `gorm:"not null"`
    ReturnDate  *time.Time // Nullable for unreturned books
    Status      string    `gorm:"type:varchar(20);not null"` // ACTIVE, OVERDUE, RETURNED
}