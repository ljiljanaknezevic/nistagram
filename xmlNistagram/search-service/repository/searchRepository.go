package repository

import (
	"html/template"
	"search-service/model"

	"github.com/jinzhu/gorm"
)

type SearchRepository struct {
	Database *gorm.DB
}

func (repo *SearchRepository) GetUserByUsername(username string) model.User {
	var user model.User
	repo.Database.Where("username = ? ", username).Preload("Following").Preload("WaitingFollowers").Preload("Followers").First(&user)
	user.Name = template.HTMLEscapeString(user.Name)
	user.Password = template.HTMLEscapeString(user.Password)
	user.PhoneNumber = template.HTMLEscapeString(user.PhoneNumber)
	user.Gender = template.HTMLEscapeString(user.Gender)
	user.Birhtday = template.HTMLEscapeString(user.Birhtday)
	user.Username = template.HTMLEscapeString(user.Username)
	user.Website = template.HTMLEscapeString(user.Website)
	user.Biography = template.HTMLEscapeString(user.Biography)
	return user
}
func (repo *SearchRepository) GetAllUsers() []model.User {
	var users []model.User
	repo.Database.Preload("Following").Preload("WaitingFollowers").Preload("Followers").Find(&users)
	return users
}

func (repo *SearchRepository) GetAllUsersExceptLogging(username string) []model.User {
	var users []model.User
	repo.Database.Where("email != ?", username).Preload("Following").Preload("WaitingFollowers").Preload("Followers").Find(&users)
	return users
}
func (repo *SearchRepository) GetUserByEmailAddress(email string) model.User {
	var user model.User
	repo.Database.Where("email = ? ", email).Preload("Following").Preload("WaitingFollowers").Preload("Followers").First(&user)
	return user
}

func (repo *SearchRepository) GetAllPosts() []model.Post {
	var posts []model.Post
	repo.Database.Find(&posts)
	return posts
}
func (repo *SearchRepository) GetPostsForSearchedUser(email string) []model.Post {
	var posts []model.Post
	repo.Database.Where("email = ?", email).Find(&posts)
	return posts
}

func (repo *SearchRepository) FindFileById(id uint) model.File {
	var file model.File
	repo.Database.Where("ID = ? ", id).First(&file)
	return file
}
