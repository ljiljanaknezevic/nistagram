package model

import "github.com/jinzhu/gorm"

type Spam struct {
	gorm.Model
	PostId string `json:"postId"`
	Email  string `json:"email"`
	Reason string `json:"reason"`
}
