package entity

import "time"

type Post struct {
	ID             int
	UserID         int
	User           User `gorm:"foreignKey:user_id" json:"user"`
	Posting        string
	PictureUrl     *string
	UploadPostings []UploadPosting `gorm:"foreignKey:PostID"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type UploadPosting struct {
	ID        int
	PostID    uint `gorm:"type:integer" json:"post_id"`
	Post      Post `gorm:"foreignKey:post_id" json:"post"`
	FileUrl   *string
	Format    *string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type LikePosting struct {
	ID        int
	PostID    uint `gorm:"type:integer" json:"post_id"`
	Post      Post `gorm:"foreignKey:post_id" json:"post"`
	UserID    int  `gorm:"type:integer" json:"user_id"`
	User      User `gorm:"foreignKey:user_id" json:"user"`
	Like      int
	Dislike   int
	CreatedAt time.Time
	UpdatedAt time.Time
}
