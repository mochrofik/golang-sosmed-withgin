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
	"time"

	"github.com/sirupsen/logrus"
)

const uploadDir = "./storage/images/"

type PostService interface {
	Posting(req *dto.PostRequest) error
}

type postService struct {
	repository repository.PostRepository
}

func NewPostService(r repository.PostRepository) *postService {
	return &postService{
		repository: r,
	}
}

const MAX_FILE_SIZE = 10 << 20 // 10MB
func (s *postService) Posting(req *dto.PostRequest) error {
	logrus.Println("request")
	logrus.Println(req)

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
			fileBytes := make([]byte, 512)

			_, err = src.Read(fileBytes)
			if err != nil && !errors.Is(err, io.EOF) {
				return &errorhandler.BadRequestError{Message: "gagal membaca byte file: " + err.Error()}
			}
			mimeType := helper.FormatFile(fileBytes, fileHeader)
			// Buat nama file unik dan path tujuan
			fileName := fmt.Sprintf("%d_%s", time.Now().Unix(), fileHeader.Filename)
			imagePath = filepath.Join(uploadDir, fileName)

			if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
				os.MkdirAll(uploadDir, 0755)
			}

			dst, err := os.Create(imagePath)
			if err != nil {
				return &errorhandler.BadRequestError{Message: "failed to create storage file: " + err.Error()}
			}
			defer dst.Close()

			// Salin konten
			if _, err := io.Copy(dst, src); err != nil {
				os.Remove(imagePath) // Bersihkan jika gagal copy
				return &errorhandler.BadRequestError{Message: "failed to save file content: " + err.Error()}
			}

			uploads = entity.UploadPosting{
				PostID:  uint(post.ID),
				FileUrl: &imagePath,
				Format:  &mimeType,
			}

			result := s.repository.UploadFiles(&uploads)

			if result != nil {
				return &errorhandler.InternalServerError{Message: err.Error()}

			}

		}

	}

	return nil
}
