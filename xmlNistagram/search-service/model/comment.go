package model

import (
	"github.com/jinzhu/gorm"
)

type Comment struct {
	gorm.Model
	Text   string `json:"text"`
	Email  string `json:"email"`
	PostID string `json:"postID"`
}
