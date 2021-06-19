package service

import (
	"search-service/model"
	"search-service/repository"
)

type SearchService struct {
	Repo        *repository.SearchRepository
	CommentRepo *repository.CommentRepository
}

func (service *SearchService) SaveComment(comment *model.Comment) error {
	service.CommentRepo.CreateComment(comment)
	return nil
}
func (service *SearchService) GetUserByUsername(username string) model.User {
	user := service.Repo.GetUserByUsername(username)
	return user
}
func (service *SearchService) GetAllUsers() []model.User {
	users := service.Repo.GetAllUsers()
	return users
}
func (service *SearchService) GetAllUsersExceptLogging(username string) []model.User {
	users := service.Repo.GetAllUsersExceptLogging(username)
	return users
}
func (service *SearchService) GetUserByEmailAddress(email string) model.User {
	user := service.Repo.GetUserByEmailAddress(email)
	return user
}

func (service *SearchService) GetAllPosts() []model.Post {
	posts := service.Repo.GetAllPosts()
	return posts
}

func (service *SearchService) GetAllStories() []model.Story {
	stories := service.Repo.GetAllStories()
	return stories
}
func (service *SearchService) GetPostsForSearchedUser(email string) []model.Post {
	posts := service.Repo.GetPostsForSearchedUser(email)
	return posts
}
func (service *SearchService) FindFileById(id uint) model.File {
	file := service.Repo.FindFileById(id)
	return file
}

func (service *SearchService) GetAllCommentsByPostsID(postID string) []model.Comment {
	return service.CommentRepo.GetAllCommentsByPostsID(postID)
}
