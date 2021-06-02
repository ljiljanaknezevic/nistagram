package repository

import (
	"fmt"
	"user-service-mod/model"

	"github.com/jinzhu/gorm"
)

type UserRepository struct {
	Database *gorm.DB
}

func (repo *UserRepository) CreateUser(user *model.User) error {
	result := repo.Database.Create(user)
	fmt.Println(result.RowsAffected)
	return nil
}
func (repo *UserRepository) UpdateUser(user *model.User) error {
	result := repo.Database.Save(user)
	fmt.Println(result.RowsAffected)
	return nil
}

func (repo *UserRepository) UserExists(email string, username string) bool {
	var count int64
	repo.Database.Where("email = ? or username = ?", email, username).Find(&model.User{}).Count(&count)
	return count != 0
}

func (repo *UserRepository) GetUserByEmail(email string) bool {
	var count int64
	repo.Database.Where("email = ? ", email).Find(&model.User{}).Count(&count)
	return count != 0
}

func (repo *UserRepository) GetUserByEmailAddress(email string) model.User {
	var user model.User
	repo.Database.Where("email = ? ", email).First(&user)
	return user
}

func (repo *UserRepository) UserForLogin(email string) model.User {
	var authUser model.User
	repo.Database.Where("email = ?", email).First(&authUser)
	return authUser
}
