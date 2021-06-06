package service

import (
	"fmt"
	"post-service-mod/model"
	"post-service-mod/repository"
)

type PostService struct {
	Repo *repository.PostRepository
}

func (service *PostService) SavePost(post *model.Post) error {
	fmt.Println("'usoooooooooooooooooooooooo'")
	return nil
}
