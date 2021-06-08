package service

import (
	"story-service-mod/model"
	"story-service-mod/repository"
)

type StoryService struct {
	Repo     *repository.StoryRepository
	FileRepo *repository.FileRepository
}

func (service *StoryService) SaveStory(story *model.Story) error {
	service.Repo.CreateStory(story)
	return nil
}
func (service *StoryService) SaveFile(file *model.File) error {
	service.FileRepo.CreateFile(file)
	return nil
}

func (service *StoryService) FindFileIdByPath(path string) uint {
	return service.FileRepo.FindIdByPath(path)
}
func (service *StoryService) FindFilePathById(imageID uint) string {
	return service.FileRepo.FindFilePathById(imageID)
}
func (service *StoryService) GetAllStoriesByEmail(email string) []model.Story {
	return service.Repo.GetAllStoriesByEmail(email)
}
