package repository

import (
	"fmt"
	"search-service/model"

	"github.com/jinzhu/gorm"
)

type CommentRepository struct {
	Database *gorm.DB
}

func (repo *CommentRepository) CreateComment(comment *model.Comment) error {
	result := repo.Database.Create(comment)
	fmt.Println(result.RowsAffected)
	return nil
}
func (repo *CommentRepository) GetAllCommentsByPostsID(postID string) []model.Comment {
	var comments []model.Comment
	fmt.Println(postID)
	repo.Database.Where("post_id = ? ", postID).Find(&comments)
	fmt.Println(comments)
	return comments
}
