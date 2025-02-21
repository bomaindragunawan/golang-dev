package models

import "gorm.io/gorm"

type BlacklistToken struct {
	gorm.Model
	Token string `json:"token" gorm:"unique"`
}
