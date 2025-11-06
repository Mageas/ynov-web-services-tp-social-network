package dto

// TokenResponse represents the authentication token response
type TokenResponse struct {
	Token string `json:"token"`
}

// PostResponse represents a post in API responses
type PostResponse struct {
	ID         string `json:"id"`
	Author     string `json:"author"`
	Content    string `json:"content"`
	CreatedAt  int64  `json:"createdAt"`
	LikesCount int    `json:"likesCount"`
}

// LikesCountResponse represents the likes count response
type LikesCountResponse struct {
	LikesCount int `json:"likesCount"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Status           int               `json:"status"`
	Message          string            `json:"message"`
	ValidationErrors map[string]string `json:"validationErrors,omitempty"`
}
