package repository

import (
	"golang-sosmed-gin/entity"

	"gorm.io/gorm"
)

type PostRepository interface {
	Posting(req *entity.Post) (*entity.Post, error)
	UserExist(id int) bool
	UploadFiles(req *entity.UploadPosting) error
	MyPost(userId int) *[]entity.Post
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

	r.db.Preload("UploadPostings").Joins("User").Find(&posting, "user_id", userId)

	return &posting
}
