package repository

import (
	"fmt"
	"post-service-mod/model"

	"github.com/jinzhu/gorm"
)

type PostRepository struct {
	Database *gorm.DB
}

func (repo *PostRepository) CreatePost(post *model.Post) error {
	result := repo.Database.Create(post)
	fmt.Println(result.RowsAffected)
	return nil
}

func (repo *PostRepository) GetAllPostsByEmail(email string) []model.Post {
	var posts []model.Post
	repo.Database.Where("email = ? ", email).Find(&posts)
	return posts
}
