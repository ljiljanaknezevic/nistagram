package repository

import (
	"fmt"
	"html/template"
	"post-service-mod/model"
	"strconv"

	"github.com/jinzhu/gorm"
)

type PostRepository struct {
	Database *gorm.DB
}

func (repo *PostRepository) CreatePost(post *model.Post) error {
	repo.Database.Create(post)
	return nil
}

func (repo PostRepository) CreateSpam(spam model.Spam) bool {
	if !repo.SpamExists(spam.Email, spam.PostId) {
		fmt.Println("usao u if za request")
		result := repo.Database.Create(spam)
		fmt.Println(result.RowsAffected)
		return true
	}
	return false
}
func (repo *PostRepository) SpamExists(email string, postId string) bool {
	var count int64
	repo.Database.Where("email = ? and postId = ?", email, postId).Find(&model.Spam{}).Count(&count)
	return count != 0
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

func (repo *PostRepository) GetPostById(postId string) model.Post {
	var post model.Post

	u64, err := strconv.ParseUint(postId, 10, 32)
	if err != nil {
		fmt.Println(err)
	}
	wd := uint(u64)

	repo.Database.Where("ID = ? ", wd).Preload("Likes").Find(&post)

	return post
}
func (repo *PostRepository) UpdatePost(post *model.Post) error {
	repo.Database.Preload("Likes").Save(post)
	return nil
}
func (repo *PostRepository) GetAllPosts() []model.Post {
	var posts []model.Post
	repo.Database.Preload("Likes").Find(&posts)
	return posts
}
func (repo *PostRepository) Dislike(ID uint) error {
	repo.Database.Where("ID = ?", ID).Preload("Likes").Delete(&model.Like{})
	return nil
}
