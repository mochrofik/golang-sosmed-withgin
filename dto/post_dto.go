package dto

import "mime/multipart"

type PostResponse struct {
	ID         int    `json:"id"`
	UserID     int    `json:"-"`
	User       User   `gorm:"foreignKey:UserID" json:"user"`
	Posting    string `json:"posting"`
	PictureUrl string `json:"picture_url"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

type PostRequest struct {
	UserID  int                     `form:"user_id"`
	Posting string                  `form:"posting"`
	Files   []*multipart.FileHeader `form:"files"`
}

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
