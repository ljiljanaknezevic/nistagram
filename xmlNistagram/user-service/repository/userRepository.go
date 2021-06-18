package repository

import (
	"fmt"
	"html/template"
	"user-service-mod/model"

	"github.com/jinzhu/gorm"
)

type UserRepository struct {
	Database *gorm.DB
}

func (repo *UserRepository) CreateUser(user *model.User) bool {
	if !repo.UserExists(user.Email, user.Username) {
		result := repo.Database.Create(user)
		fmt.Println(result.RowsAffected)
		return true
	}
	return false
}
func (repo *UserRepository) CreateRequest(request *model.VerificationRequest) bool {
	fmt.Println(request.Email)
	if !repo.RequestExists(request.Email) {
		fmt.Println("usao u if za request")
		result := repo.Database.Create(request)
		fmt.Println(result.RowsAffected)
		return true
	}
	return false
}
func (repo *UserRepository) GetAllUsersExceptLogging(email string) []model.User {
	var users []model.User
	repo.Database.Where("email != ?", email).Preload("Following").Preload("WaitingFollowers").Preload("Followers").Preload("Blocked").Preload("UsersWhoBlocked").Find(&users)
	return users
}
func (repo *UserRepository) GetAllUsersExceptLoggingForTag(email string) []model.User {
	var users []model.User
	var isTrue bool
	isTrue = true
	repo.Database.Where("email != ? and can_tag = ?", email, isTrue).Preload("Following").Preload("WaitingFollowers").Preload("Blocked").Preload("Followers").Preload("UsersWhoBlocked").Find(&users)
	return users
}

func (repo *UserRepository) GetAllRequests() []model.VerificationRequest {
	var requests []model.VerificationRequest
	repo.Database.Find(&requests)
	return requests
}

func (repo *UserRepository) UpdateUser(user *model.User) error {
	result := repo.Database.Preload("Following").Preload("WaitingFollowers").Preload("Followers").Preload("Blocked").Preload("UsersWhoBlocked").Save(user)
	fmt.Println(result.RowsAffected)
	return nil
}
func (repo *UserRepository) DeleteFromWaitingList(ID uint) error {
	repo.Database.Where("ID = ?", ID).Preload("Following").Preload("WaitingFollowers").Preload("Followers").Preload("Blocked").Preload("UsersWhoBlocked").Delete(&model.WaitingFollower{})
	return nil
}

func (repo *UserRepository) DeleteVerificationRequest(email string) error {
	//repo.Database.Where("email = ?", email).Delete(&model.VerificationRequest{})
	repo.Database.Exec("DELETE FROM verification_requests WHERE email=$1;", email)
	return nil
}

func (repo *UserRepository) UserExists(email string, username string) bool {
	var count int64
	repo.Database.Where("email = ? or username = ?", email, username).Find(&model.User{}).Count(&count)
	return count != 0
}
func (repo *UserRepository) RequestExists(email string) bool {
	var count int64
	repo.Database.Where("email = ?", email).Find(&model.VerificationRequest{}).Count(&count)
	return count != 0
}

func (repo *UserRepository) GetUserByEmail(email string) bool {
	var count int64
	repo.Database.Where("email = ? ", email).Find(&model.User{}).Count(&count)
	return count != 0
}
func (repo *UserRepository) GetWaitingUser(username string) model.WaitingFollower {
	var user model.WaitingFollower
	repo.Database.Where("username = ? ", username).First(&user)
	return user
}

func (repo *UserRepository) GetUserByEmailAddress(email string) model.User {
	var user model.User
	repo.Database.Where("email = ? ", email).Preload("Following").Preload("WaitingFollowers").Preload("Followers").Preload("Blocked").Preload("UsersWhoBlocked").First(&user)

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
func (repo *UserRepository) GetUserByUsername(username string) model.User {
	var user model.User
	repo.Database.Where("username = ? ", username).Preload("Following").Preload("WaitingFollowers").Preload("Followers").Preload("Blocked").Preload("UsersWhoBlocked").First(&user)
	return user
}

func (repo *UserRepository) UserForLogin(email string) model.User {
	var authUser model.User
	repo.Database.Where("email = ?", email).Preload("Following").Preload("WaitingFollowers").Preload("Followers").Preload("Blocked").Preload("UsersWhoBlocked").First(&authUser)
	return authUser
}
