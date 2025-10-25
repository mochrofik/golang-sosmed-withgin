package service

import (
	"golang-sosmed-gin/dto"
	"golang-sosmed-gin/entity"
	"golang-sosmed-gin/repository"
)

type UserService interface {
	GetAllUser(req *dto.UserRequest) *[]entity.User
	GetMyProfile(ID int) *entity.User
}

type userService struct {
	repository repository.UserRepository
}

func NewUserService(r repository.UserRepository) *userService {
	return &userService{
		repository: r,
	}
}

func (s *userService) GetAllUser(req *dto.UserRequest) *[]entity.User {

	user := s.repository.GetAllUser(req)

	return user
}

func (s *userService) GetMyProfile(ID int) *entity.User {
	user := s.repository.GetMyProfile(ID)

	return user
}
