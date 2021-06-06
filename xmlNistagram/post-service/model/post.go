package model

import (
	"github.com/jinzhu/gorm"
)

type Post struct {
	gorm.Model
	Description string `json:"description"`
	Username    string `gorm: "unique" json:"username"`
	Tags        string `json:"tags"`
	Image       string `json:image`
}
