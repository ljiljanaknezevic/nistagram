package service

import (
	"post-service-mod/model"
	"post-service-mod/repository"
)

type PostService struct {
	Repo     *repository.PostRepository
	FileRepo *repository.FileRepository
}

func (service *PostService) SavePost(post *model.Post) error {
	service.Repo.CreatePost(post)
	return nil
}
func (service *PostService) SaveFile(file *model.File) error {
	service.FileRepo.CreateFile(file)
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
