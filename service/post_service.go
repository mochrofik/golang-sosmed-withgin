package service

import (
	"errors"
	"fmt"
	"golang-sosmed-gin/dto"
	"golang-sosmed-gin/entity"
	"golang-sosmed-gin/errorhandler"
	"golang-sosmed-gin/helper"
	"golang-sosmed-gin/repository"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const uploadDir = "./storage/images/"

type PostService interface {
	Posting(req *dto.PostRequest) error
	MyPost(userId int) *[]dto.MyPost
	DeletePost(ID int) error
}

type postService struct {
	repository repository.PostRepository
}

func NewPostService(r repository.PostRepository) *postService {
	return &postService{
		repository: r,
	}
}

const MAX_FILE_SIZE = 5 * 1024 * 1024 // 5MB
func (s *postService) Posting(req *dto.PostRequest) error {

	if userExist := s.repository.UserExist(req.UserID); !userExist {
		return &errorhandler.BadRequestError{Message: "User not found"}
	}

	newPost := entity.Post{
		UserID:  req.UserID,
		Posting: req.Posting,
	}

	post, err := s.repository.Posting(&newPost)

	if err != nil {
		return &errorhandler.InternalServerError{Message: err.Error()}
	}

	if post != nil {
		for _, fileHeader := range req.Files {
			var uploads entity.UploadPosting
			var imagePath string
			if fileHeader.Size > MAX_FILE_SIZE {
				return &errorhandler.InternalServerError{Message: "File terlalu besar" + fileHeader.Filename}
			}
			// Buka file yang di-upload
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
			mimeType := helper.FormatFile(fileBytes, fileHeader)
			if _, seekErr := src.Seek(0, io.SeekStart); seekErr != nil {
				return &errorhandler.BadRequestError{Message: "gagal memutar ulang file stream: " + seekErr.Error()}
			}

			extension := filepath.Ext(fileHeader.Filename)
			originalBaseName := strings.TrimSuffix(fileHeader.Filename, extension)
			cleanedBaseName := strings.ReplaceAll(originalBaseName, " ", "_")
			fileName := fmt.Sprintf("%s_%d%s", cleanedBaseName, time.Now().UnixMilli(), extension)
			imagePath = filepath.Join(uploadDir, fileName)

			newFile, err := os.Create(imagePath)
			if err != nil {
				return &errorhandler.BadRequestError{Message: "failed to create storage file: " + err.Error()}
			}
			defer newFile.Close()

			// Salin konten
			if _, err := io.Copy(newFile, src); err != nil {
				return &errorhandler.BadRequestError{Message: "failed to save file content: " + err.Error()}
			}

			uploads = entity.UploadPosting{
				PostID:  uint(post.ID),
				FileUrl: &imagePath,
				Format:  &mimeType,
			}

			result := s.repository.UploadFiles(&uploads)

			if result != nil {
				return &errorhandler.InternalServerError{Message: result.Error()}

			}

		}

	}

	return nil
}

func (s *postService) MyPost(userID int) *[]dto.MyPost {

	var post []dto.MyPost

	result := s.repository.MyPost(userID)

	for _, v := range *result {

		var files []dto.FilePosting

		for _, file := range v.UploadPostings {

			files = append(files, dto.FilePosting{
				ID:      file.ID,
				FileUrl: file.FileUrl,
				Format:  file.Format,
			})
		}

		post = append(post, dto.MyPost{
			ID:             v.ID,
			Posting:        v.Posting,
			UploadPostings: files,
			UserID:         v.UserID,
			User: dto.User{
				ID:    v.UserID,
				Name:  v.User.Name,
				Email: v.User.Email,
			},
			CreatedAt: helper.FormatDateTimeToString(v.CreatedAt),
			UpdatedAt: helper.FormatDateTimeToString(v.UpdatedAt),
		})

	}

	return &post
}

func (s *postService) DeletePost(ID int) error {

	err := s.repository.DeletePost(ID)
	return err
}
