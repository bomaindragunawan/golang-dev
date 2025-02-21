package models

import (
	"time"

	"gorm.io/gorm"
)

type PasswordReset struct {
	gorm.Model
	Email     string `json:"email" gorm:"unique"`
	OTP       string `json:"otp"`
	ExpiresAt time.Time
}
