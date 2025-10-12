package dto

type UserRequest struct {
	Search *string `json:"search"`
	// Paginate *Paginate `json:"paginate"`
}

type UserResponse struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
