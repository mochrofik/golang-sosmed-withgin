package service

import (
	"golang-sosmed-gin/dto"
	"golang-sosmed-gin/entity"
	"golang-sosmed-gin/errorhandler"
	"golang-sosmed-gin/helper"
	"golang-sosmed-gin/repository"
)

type AuthService interface {
	Register(req *dto.RegisterRequest) error
	Login(req *dto.LoginRequest) (*dto.LoginResponse, error)
}

type authService struct {
	repository repository.AuthRepository
}

func NewAuthService(r repository.AuthRepository) *authService {
	return &authService{
		repository: r,
	}
}

func (s *authService) Register(req *dto.RegisterRequest) error {

	if req.Name == "" || req.Email == "" || req.Gender == "" || req.Password == "" {
		return &errorhandler.BadRequestError{Message: "incorrectly data"}

	}
	if emailExist := s.repository.EmailExist(req.Email); emailExist {
		return &errorhandler.BadRequestError{Message: "email already exist"}
	}

	if req.Password != req.PasswordConfirm {
		return &errorhandler.BadRequestError{Message: "password not match"}

	}

	passwordHash, err := helper.HashPassword(req.Password)

	if err != nil {
		return &errorhandler.InternalServerError{
			Message: err.Error(),
		}
	}

	user := entity.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: passwordHash,
		Gender:   req.Gender,
	}

	if err := s.repository.Register(&user); err != nil {
		return &errorhandler.InternalServerError{Message: err.Error()}
	}

	return nil
}

func (s *authService) Login(req *dto.LoginRequest) (*dto.LoginResponse, error) {

	var data dto.LoginResponse

	if req.Email == "" || req.Password == "" {
		return nil, &errorhandler.BadRequestError{Message: "incorrectly data"}

	}

	user, err := s.repository.GetUserByEmail(req.Email)

	if err != nil {
		return nil, &errorhandler.NotFoundError{Message: "wrong email or password"}
	}

	if err := helper.VerifyPassword(user.Password, req.Password); err != nil {
		return nil, &errorhandler.NotFoundError{Message: "wrong email or password"}

	}

	token, err := helper.CreateToken(user)

	if err != nil {
		return nil, &errorhandler.InternalServerError{Message: err.Error()}
	}

	data = dto.LoginResponse{
		ID:    user.ID,
		Name:  user.Name,
		Token: token,
	}

	return &data, nil
}
