package model

import (
	"github.com/jinzhu/gorm"
)

type File struct {
	gorm.Model
	Path string `json:"path"`
	Type string `json:"type"`
}
