package model

import (
	"github.com/jinzhu/gorm"
)

type Highlight struct {
	gorm.Model
	Title   string   `json:"title"`
	Stories []string `gorm:"many2many:highligh_stories; json:"stories"`
}
