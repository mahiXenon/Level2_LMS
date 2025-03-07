package models

import (
	"time"
)

type User struct {
	ID            uint   `json:"id" gorm:"primaryKey"`
	Name          string `json:"name"`
	Password      string `json:"password"`
	Email         string `json:"email" gorm:"unique"`
	ContactNumber string `json:"contact_number" gorm:"unique"`
	Role          string `json:"role"`
	// Libraries     []Library `gorm:"many2many:library_users;"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Library struct {
	ID   uint   `json:"id" gorm:"primaryKey"`
	Name string `json:"name"`
	// Users []User `gorm:"many2many:library_users;"`
}

type LibraryUser struct {
	UserId    uint `json:"user_id" gorm:"primaryKey;foreignKey:UserId;references:ID"`
	LibraryId uint `json:"library_id" gorm:"primaryKey;foreignKey:LibraryId;references:ID"`
}

type AuthLibrary struct {
	Name string `json:"name" binding:"required"`
}

type AuthInput struct {
	Name          string `json:"name" binding:"required"`
	Email         string `json:"email" binding:"required"`
	Password      string `json:"password" binding:"required"`
	ContactNumber string `json:"contact_number" binding:"required" gorm:"unique"`
}

type AuthLogin struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type Admin struct {
	ID uint `json:"id" binding:"required"`
}

type RegisterLibrary struct {
	Name string `json:"name" binding:"required"`
}
