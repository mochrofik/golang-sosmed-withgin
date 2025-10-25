package service

import (
	"errors"
	"fmt"
	"golang-sosmed-gin/dto"
	"golang-sosmed-gin/entity"
	"golang-sosmed-gin/errorhandler"
	"golang-sosmed-gin/repository"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type UserService interface {
	GetAllUser(req *dto.UserRequest) *[]entity.User
	GetMyProfile(ID int) *dto.UserResponse
	EditProfile(req *dto.ProfileRequest) error
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

func (s *userService) GetMyProfile(ID int) *dto.UserResponse {
	user := s.repository.GetMyProfile(ID)

	respone := dto.UserResponse{
		ID:      user.ID,
		Email:   user.Email,
		Name:    user.Name,
		Profile: user.Profile,
		Gender:  user.Gender,
	}

	return &respone

}
func (s *userService) EditProfile(req *dto.ProfileRequest) error {

	var userEdit entity.User

	if userExist := s.repository.UserExist(req.UserID); !userExist {
		return &errorhandler.BadRequestError{Message: "User not found"}
	}
	var imagePath *string
	if req.Profile != nil {
		logrus.Info("masuk")
		fileHeader := req.Profile
		if fileHeader.Size > MAX_FILE_SIZE {
			return &errorhandler.InternalServerError{Message: "File terlalu besar" + fileHeader.Filename}
		}

		src, err := fileHeader.Open()
		if err != nil {
			return &errorhandler.BadRequestError{Message: "failed to open uploaded file: " + err.Error()}
		}

		defer src.Close()
		if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
			os.MkdirAll(uploadDir, 0755)
		}
		fileBytes := make([]byte, 512)

		_, err = src.Read(fileBytes)
		if err != nil && !errors.Is(err, io.EOF) {
			return &errorhandler.BadRequestError{Message: "gagal membaca byte file: " + err.Error()}
		}
		if _, seekErr := src.Seek(0, io.SeekStart); seekErr != nil {
			return &errorhandler.BadRequestError{Message: "gagal memutar ulang file stream: " + seekErr.Error()}
		}

		extension := filepath.Ext(fileHeader.Filename)
		originalBaseName := strings.TrimSuffix(fileHeader.Filename, extension)
		cleanedBaseName := strings.ReplaceAll(originalBaseName, " ", "_")
		fileName := fmt.Sprintf("%s_%d%s", cleanedBaseName, time.Now().UnixMilli(), extension)
		imagePath := filepath.Join(uploadDir, fileName)

		newFile, err := os.Create(imagePath)
		if err != nil {
			return &errorhandler.BadRequestError{Message: "failed to create storage file: " + err.Error()}
		}
		defer newFile.Close()

		// Salin konten
		if _, err := io.Copy(newFile, src); err != nil {
			return &errorhandler.BadRequestError{Message: "failed to save file content: " + err.Error()}
		}
	}

	userEdit = entity.User{
		ID:      req.UserID,
		Name:    req.Name,
		Profile: *imagePath,
	}

	result := s.repository.EditProfile(&userEdit)

	if result != nil {
		return &errorhandler.InternalServerError{Message: result.Error()}

	}

	return nil
}
