package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name         string `json:"name" binding:"required"`
	Email        string `json:"email" gorm:"unique" binding:"required"`
	Password     string `json:"password" binding:"required"`
	RefreshToken string `json:"refresh_token,omitempty"`
}
