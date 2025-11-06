package dto

// SignupRequest represents the signup request payload
type SignupRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginRequest represents the login request payload
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// CreatePostRequest represents the create post request payload
type CreatePostRequest struct {
	Content string `json:"content"`
}
