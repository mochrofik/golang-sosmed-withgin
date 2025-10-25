package dto

import "mime/multipart"

type UserRequest struct {
	Search *string `json:"search"`
	// Paginate *Paginate `json:"paginate"`
}

type UserResponse struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Profile string `json:"profile"`
	Gender  string `json:"gender"`
}

type ProfileRequest struct {
	UserID  int                   `form:"id"`
	Name    string                `form:"name"`
	Profile *multipart.FileHeader `form:"profile"`
}
