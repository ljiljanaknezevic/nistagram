package repository

import (
	"github.com/jinzhu/gorm"
	"search-service/model"
)

type SearchRepository struct {
	Database *gorm.DB
}
func (repo *SearchRepository) GetUserByUsername(username string) model.User {
	var user model.User
	repo.Database.Where("username = ? ", username).First(&user)
	return user
}
func (repo *SearchRepository) GetAllUsers() []model.User{
	var users []model.User
	repo.Database.Preload("Followers").Find(&users)
	return users
}

func (repo *SearchRepository) GetAllUsersExceptLogging(username string) []model.User{
	var users []model.User
	repo.Database.Where("email != ?", username).Preload("Followers").Find(&users)
	return users
}
