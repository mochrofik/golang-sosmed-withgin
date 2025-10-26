package repository

import (
	"fmt"
	"golang-sosmed-gin/entity"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type PostRepository interface {
	Posting(req *entity.Post) (*entity.Post, error)
	UserExist(id int) bool
	UploadFiles(req *entity.UploadPosting) error
	MyPost(userId int) *[]entity.Post
	DeletePost(ID int) error
	LikePost(PostID int, UserID int, isLike bool, to int) error
	CheckLike(PostID int, UserID int) (int, error)
	CheckLikeActive(PostID int, UserID int, like int) error
}

type postRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) *postRepository {
	return &postRepository{
		db: db,
	}
}

func (r *postRepository) Posting(post *entity.Post) (*entity.Post, error) {

	err := r.db.Create(&post).Error

	return post, err
}
func (r *postRepository) UploadFiles(upload *entity.UploadPosting) error {

	err := r.db.Create(&upload).Error

	return err
}

func (r *postRepository) UserExist(id int) bool {
	var user entity.User

	err := r.db.First(&user, "id = ?", id).Error

	return err == nil
}

func (r *postRepository) MyPost(userId int) *[]entity.Post {

	var posting []entity.Post

	r.db.Preload("UploadPostings").Preload("LikePostings").Joins("User").Find(&posting, "user_id", userId)

	return &posting
}

func (r *postRepository) DeletePost(ID int) error {
	var post *entity.Post
	var uploads []entity.UploadPosting

	err3 := r.db.Find(&uploads, "post_id", ID).Error
	if err3 != nil {
		return err3
	}

	logrus.Println("uploads")

	for _, v := range uploads {
		logrus.Println(*v.FileUrl)

		filePath := "./" + *v.FileUrl

		err1 := os.Remove(filePath)
		if err1 != nil {
			if os.IsNotExist(err1) {
				return fmt.Errorf("file tidak ditemukan di path: %s", filePath)
			}
			return fmt.Errorf("gagal menghapus file %s: %w", filePath, err1)
		}
	}

	resultUpload := r.db.Where("post_id = ?", ID).Delete(&entity.UploadPosting{})
	resultComments := r.db.Where("post_id = ?", ID).Delete(&entity.LikePosting{})

	if resultUpload.Error != nil {
		return fmt.Errorf("gagal menghapus file: %w", resultUpload.Error)
	}
	if resultComments.Error != nil {
		return fmt.Errorf("gagal menghapus file: %w", resultComments.Error)
	}

	err2 := r.db.Delete(&post, ID).Error

	return err2

}

func (r *postRepository) LikePost(PostID int, UserID int, isLike bool, to int) error {
	var likePost entity.LikePosting

	if !isLike {

		var existingLike entity.LikePosting

		result := r.db.Where("post_id = ? AND user_id = ?", PostID, UserID).First(&existingLike)

		if result.Error == nil {
			newLikeStatus := 1
			if existingLike.Like == 1 {
				newLikeStatus = 0
			}
			return r.db.Model(&existingLike).Update("like", newLikeStatus).Error
		}
		return result.Error

	} else {
		likePost = entity.LikePosting{
			PostID:    uint(PostID),
			UserID:    UserID,
			Like:      1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		err := r.db.Create(&likePost).Error

		return err
	}

}

func (r *postRepository) CheckLike(PostID int, UserID int) (int, error) {
	var likePost entity.LikePosting
	err := r.db.Where("post_id = ?", PostID).Where("user_id = ?", UserID).First(&likePost).Error

	return likePost.Like, err
}
func (r *postRepository) CheckLikeActive(PostID int, UserID int, like int) error {
	var likePost entity.LikePosting

	err := r.db.Where("post_id = ?", PostID).Where("user_id = ?", UserID).Where("like = ?", like).First(&likePost).Error
	return err
}
