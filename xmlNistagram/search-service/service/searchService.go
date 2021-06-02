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
