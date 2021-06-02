package model

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	Name        string `json:"name"`
	Email       string `gorm:"unique" json:"email"`
	Password    string `json:"password"`
	Role        string `json:"role"`
	PhoneNumber string `json:"phone"`
	Gender      string `json:"gender"`
	Birhtday    string `json:"birthday"`
	Username    string `gorm:"unique" json:"username"`
	Website     string `json:"website"`
	Biography   string `json:"biography"`
}