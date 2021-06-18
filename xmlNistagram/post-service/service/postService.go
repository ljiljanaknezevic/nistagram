package service

import (
	"post-service-mod/model"
	"post-service-mod/repository"
)

type PostService struct {
	Repo        *repository.PostRepository
	FileRepo    *repository.FileRepository
	CommentRepo *repository.CommentRepository
}

func (service *PostService) SavePost(post *model.Post) error {
	service.Repo.CreatePost(post)
	return nil
}
func (service *PostService) CreateSpam(spam *model.Spam) bool {
	return service.Repo.CreateSpam(spam)
}
func (service *PostService) SaveFile(file *model.File) error {
	service.FileRepo.CreateFile(file)
	return nil
}
func (service *PostService) SaveComment(comment *model.Comment) error {
	service.CommentRepo.CreateComment(comment)
	return nil
}

func (service *PostService) FindFileIdByPath(path string) uint {
	return service.FileRepo.FindIdByPath(path)
}
func (service *PostService) FindFilePathById(imageID uint) string {
	return service.FileRepo.FindFilePathById(imageID)
}
func (service *PostService) GetAllPostsByEmail(email string) []model.Post {
	return service.Repo.GetAllPostsByEmail(email)
}
func (service *PostService) GetAllCommentsByPostsID(postID string) []model.Comment {
	return service.CommentRepo.GetAllCommentsByPostsID(postID)
}
func (service *PostService) GetPostById(postId string) model.Post {
	return service.Repo.GetPostById(postId)
}
func (service *PostService) UpdatePost(post *model.Post) error {
	service.Repo.UpdatePost(post)
	return nil
}

func (service *PostService) GetAllPosts() []model.Post {
	return service.Repo.GetAllPosts()

}
func (service *PostService) Dislike(ID uint) error {
	service.Repo.Dislike(ID)
	return nil
}
