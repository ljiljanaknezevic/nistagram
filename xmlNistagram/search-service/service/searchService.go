package service

import (
	"search-service/model"
	"search-service/repository"
)

type SearchService struct {
	Repo *repository.SearchRepository
}
func (service *SearchService) GetUserByUsername(username string) model.User {
	user := service.Repo.GetUserByUsername(username)
	return user
}
func (service *SearchService) GetAllUsers() []model.User{
	users:= service.Repo.GetAllUsers()
	return users
}
func (service *SearchService) GetAllUsersExceptLogging(username string) []model.User{
	users:= service.Repo.GetAllUsersExceptLogging(username)
	return users
}

func (service *SearchService) GetAllPosts() []model.Post{
	posts:= service.Repo.GetAllPosts()
	return posts
}
func (service *SearchService) FindFileById(id uint) model.File{
	file := service.Repo.FindFileById(id)
	return file
}

