package repository

import (
	"golang-sosmed-gin/dto"
	"golang-sosmed-gin/entity"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserRepository interface {
	GetAllUser(req *dto.UserRequest) *[]entity.User
	GetMyProfile(ID int) *entity.User
	UserExist(id int) bool
	EditProfile(user *entity.User) error
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
func (r *userRepository) UserExist(id int) bool {
	var user entity.User

	err := r.db.First(&user, "id = ?", id).Error

	return err == nil
}

func (r *userRepository) EditProfile(user *entity.User) error {

	logrus.Info(user)

	var userData entity.User

	r.db.First(&userData, user.ID)

	if user.Name != "" {
		userData.Name = user.Name
	}

	if user.Profile != "" {
		userData.Profile = user.Profile
	}

	if userData.Name != "" || userData.Profile != "" {
		result := r.db.Save(&userData).Error

		return result
	}

	return nil

	// result := r.db.First(&user, "id", user.ID).Updates(entity.User{Name: user.Name, Profile: user.Profile}).Error

}
