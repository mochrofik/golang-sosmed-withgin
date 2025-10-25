package repository

import (
	"golang-sosmed-gin/dto"
	"golang-sosmed-gin/entity"

	"gorm.io/gorm"
)

type UserRepository interface {
	GetAllUser(req *dto.UserRequest) *[]entity.User
	GetMyProfile(ID int) *entity.User
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) GetAllUser(req *dto.UserRequest) *[]entity.User {

	var users []entity.User

	query := r.db
	if req.Search != nil && *req.Search != "" {
		search := *req.Search
		query = query.Where("name LIKE ?", "%"+search+"%").Or("email LIKE ?", "%"+search+"%")
	}
	query.Find(&users)

	return &users
}

func (r *userRepository) GetMyProfile(ID int) *entity.User {
	var user entity.User

	r.db.First(&user, "id = ?", ID)

	return &user
}
