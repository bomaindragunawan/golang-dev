package models

import "gorm.io/gorm"

// Struktur request untuk registrasi
type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role"` // Bisa kosong atau diisi
}

type UserResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type User struct {
	gorm.Model
	Name         string `json:"name" binding:"required"`
	Email        string `json:"email" gorm:"unique" binding:"required"`
	Password     string `json:"password" binding:"required"`
	Role         string `json:"role" gorm:"default:user"`
	RefreshToken string `json:"refresh_token,omitempty"`
}
