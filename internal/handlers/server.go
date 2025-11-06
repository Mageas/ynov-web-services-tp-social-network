package handlers

import (
	"ynov-social-api/internal/storage"
)

type Server struct {
	Users     storage.UserStore
	Posts     storage.PostStore
	JWTSecret []byte
}

func NewServer(users storage.UserStore, posts storage.PostStore, jwtSecret []byte) *Server {
	return &Server{Users: users, Posts: posts, JWTSecret: jwtSecret}
}
