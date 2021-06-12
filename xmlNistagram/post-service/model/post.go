package model

import (
	"github.com/jinzhu/gorm"
)

type Post struct {
	gorm.Model
	Description string `json:"description"`
	Email       string `json:"email"`
	Tags        string `json:"tags"`
	ImageID     uint   `json:imageID`
	Location    string `json:location`
}
