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
	repo.Database.Find(&users)
	return users
}
