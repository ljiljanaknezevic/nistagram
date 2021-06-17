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
	Likes       []Like `gorm:"many2many:post_likes; json:"likes"`
}
type Like struct {
	gorm.Model
	Username string `json:"username"`
}
