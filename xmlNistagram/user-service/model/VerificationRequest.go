package model

import "github.com/jinzhu/gorm"

type VerificationRequest struct {
	gorm.Model
	Email    string `gorm:"unique" json:"email"`
	FullName string `json:"fullName"`
	Category string `json:"category"`
	Image    string `json:"image"`
}