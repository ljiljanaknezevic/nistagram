package repository

import (
	"fmt"
	"story-service-mod/model"

	"github.com/jinzhu/gorm"
)

type FileRepository struct {
	Database *gorm.DB
}

func (repo *FileRepository) CreateFile(file *model.File) error {
	result := repo.Database.Create(file)
	fmt.Println(result.RowsAffected)
	return nil
}

func (repo *FileRepository) FindIdByPath(path string) uint {
	var file model.File
	repo.Database.Where("path = ? ", path).First(&file)
	fmt.Println(file.ID)
	fmt.Println(file)
	return file.ID
}
func (repo *FileRepository) FindFilePathById(imageID uint) string {
	var file model.File
	repo.Database.Where("id = ? ", imageID).First(&file)
	return file.Path
}
