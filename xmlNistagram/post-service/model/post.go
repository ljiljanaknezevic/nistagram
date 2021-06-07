package model

import (
	"github.com/jinzhu/gorm"
)

type Post struct {
	gorm.Model
	Description string `json:"description"`
	Username    string `gorm: "unique" json:"username"`
	Tags        string `json:"tags"`
	ImageID     uint   `json:imageID`
	Location    string `json:location`
}
