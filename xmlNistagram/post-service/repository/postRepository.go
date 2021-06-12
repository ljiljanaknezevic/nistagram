package repository

import (
	"fmt"
	"html/template"
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

	var newPosts []model.Post
	//newPosts := []*model.Post{}

	for _, element := range posts {
		element.Tags = template.HTMLEscapeString(element.Tags)
		element.Location = template.HTMLEscapeString(element.Location)
		element.Description = template.HTMLEscapeString(element.Description)
		newPosts = append(newPosts, element)
	}
	return newPosts
}
