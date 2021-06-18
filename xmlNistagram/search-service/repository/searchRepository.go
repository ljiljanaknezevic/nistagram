package repository

import (
	"search-service/model"

	"github.com/jinzhu/gorm"
)

type SearchRepository struct {
	Database *gorm.DB
}

func (repo *SearchRepository) GetUserByUsername(username string) model.User {
	var user model.User
	repo.Database.Where("username = ? ", username).Preload("Muted").Preload("Following").Preload("WaitingFollowers").Preload("Followers").Preload("UsersWhoBlocked").Preload("Blocked").First(&user)
	return user
}
func (repo *SearchRepository) GetAllUsers() []model.User {
	var users []model.User

	repo.Database.Preload("Muted").Preload("Following").Preload("WaitingFollowers").Preload("Followers").Preload("Blocked").Preload("UsersWhoBlocked").Find(&users)
	return users
}

func (repo *SearchRepository) GetAllUsersExceptLogging(username string) []model.User {
	var users []model.User
	repo.Database.Where("email != ?", username).Preload("UsersWhoBlocked").Preload("Muted").Preload("Following").Preload("WaitingFollowers").Preload("Followers").Preload("Blocked").Find(&users)

	return users
}
func (repo *SearchRepository) GetUserByEmailAddress(email string) model.User {
	var user model.User
	repo.Database.Where("email = ? ", email).Preload("UsersWhoBlocked").Preload("Muted").Preload("Following").Preload("WaitingFollowers").Preload("Followers").Preload("Blocked").First(&user)
	return user
}

func (repo *SearchRepository) GetAllPosts() []model.Post {
	var posts []model.Post
	repo.Database.Find(&posts)
	return posts
}

func (repo *SearchRepository) GetAllStories() []model.Story {
	var stories []model.Story
	repo.Database.Find(&stories)
	return stories
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
