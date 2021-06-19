package repository

import (
	"fmt"
	"story-service-mod/model"

	"github.com/jinzhu/gorm"
)

type StoryRepository struct {
	Database *gorm.DB
}

func (repo *StoryRepository) CreateStory(story *model.Story) error {
	result := repo.Database.Create(story)
	fmt.Println(result.RowsAffected)
	return nil
}

func (repo *StoryRepository) GetAllStoriesByEmail(email string) []model.Story {
	var stories []model.Story
	repo.Database.Where("email = ? ", email).Find(&stories)
	return stories
}

