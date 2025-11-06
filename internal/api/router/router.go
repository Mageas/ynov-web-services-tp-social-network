package router

import (
	"net/http"

	"ynov-social-api/internal/api/handler"
	"ynov-social-api/internal/api/middleware"
	"ynov-social-api/internal/api/response"
	"ynov-social-api/internal/pkg/apperrors"
	"ynov-social-api/internal/service/auth"
)

// New creates and configures the application router
func New(authHandler *handler.AuthHandler, postHandler *handler.PostHandler, jwtService *auth.JWTService) http.Handler {
	mux := http.NewServeMux()

	// Public routes
	mux.HandleFunc("/signup", authHandler.Signup)
	mux.HandleFunc("/login", authHandler.Login)

	// Protected routes
	authMiddleware := middleware.Auth(jwtService)

	// Posts routes
	mux.Handle("/posts", authMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			postHandler.ListPosts(w, r)
		case http.MethodPost:
			postHandler.CreatePost(w, r)
		default:
			response.Error(w, apperrors.New(http.StatusMethodNotAllowed, "method not allowed"))
		}
	})))

	// Post actions (like/unlike)
	mux.Handle("/posts/", authMiddleware(http.HandlerFunc(postHandler.HandlePostAction)))

	return mux
}
