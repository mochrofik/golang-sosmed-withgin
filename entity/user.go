package entity

import "time"

type User struct {
	ID        int
	Name      string
	Email     string
	Password  string
	Gender    string
	Profile   string
	CreatedAt time.Time
	UpdatedAt time.Time
}
