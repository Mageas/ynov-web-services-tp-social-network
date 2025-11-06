package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"ynov-social-api/internal/api/dto"
	"ynov-social-api/internal/api/middleware"
	"ynov-social-api/internal/api/response"
	"ynov-social-api/internal/domain/post"
	"ynov-social-api/internal/pkg/apperrors"
	"ynov-social-api/internal/pkg/logger"
	postService "ynov-social-api/internal/service/post"
)

// PostHandler handles post endpoints
type PostHandler struct {
	postService *postService.Service
	logger      *logger.Logger
}

// NewPostHandler creates a new post handler
func NewPostHandler(postService *postService.Service, logger *logger.Logger) *PostHandler {
	return &PostHandler{
		postService: postService,
		logger:      logger,
	}
}

// CreatePost handles post creation
func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.Error(w, apperrors.New(http.StatusMethodNotAllowed, "method not allowed"))
		return
	}

	var req dto.CreatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, apperrors.New(http.StatusBadRequest, "invalid JSON"))
		return
	}

	author := middleware.GetUserEmail(r)
	if author == "" {
		response.Error(w, apperrors.ErrUnauthorized)
		return
	}

	post, err := h.postService.CreatePost(r.Context(), author, req.Content)
	if err != nil {
		h.logger.Error("Failed to create post: %v", err)
		response.Error(w, err)
		return
	}

	resp := h.mapPostToDTO(post)
	response.Created(w, resp)
}

// ListPosts handles post listing
func (h *PostHandler) ListPosts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.Error(w, apperrors.New(http.StatusMethodNotAllowed, "method not allowed"))
		return
	}

	query := r.URL.Query()

	var beforeTimestamp int64
	if beforeTsStr := query.Get("beforeTs"); beforeTsStr != "" {
		beforeTimestamp, _ = strconv.ParseInt(beforeTsStr, 10, 64)
	}

	page, _ := strconv.Atoi(query.Get("page"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(query.Get("limit"))
	if limit <= 0 || limit > 50 {
		limit = 10
	}

	posts, err := h.postService.ListPosts(r.Context(), beforeTimestamp, page, limit)
	if err != nil {
		h.logger.Error("Failed to list posts: %v", err)
		response.Error(w, err)
		return
	}

	resp := make([]dto.PostResponse, 0, len(posts))
	for _, p := range posts {
		resp = append(resp, h.mapPostToDTO(p))
	}

	response.OK(w, resp)
}

// HandlePostAction handles post actions (like/unlike)
func (h *PostHandler) HandlePostAction(w http.ResponseWriter, r *http.Request) {
	// Parse URL: /posts/{id}/like or /posts/{id}/unlike
	path := strings.TrimPrefix(r.URL.Path, "/posts/")
	parts := strings.Split(path, "/")

	if len(parts) != 2 || parts[0] == "" {
		response.Error(w, apperrors.New(http.StatusNotFound, "not found"))
		return
	}

	postID := parts[0]
	action := parts[1]

	if action != "like" && action != "unlike" {
		response.Error(w, apperrors.New(http.StatusNotFound, "not found"))
		return
	}

	userEmail := middleware.GetUserEmail(r)
	if userEmail == "" {
		response.Error(w, apperrors.ErrUnauthorized)
		return
	}

	var likesCount int
	var err error
	switch action {
	case "like":
		if r.Method != http.MethodPost {
			response.Error(w, apperrors.New(http.StatusMethodNotAllowed, "method not allowed"))
			return
		}
		likesCount, err = h.postService.LikePost(r.Context(), userEmail, postID)
	case "unlike":
		if r.Method != http.MethodDelete {
			response.Error(w, apperrors.New(http.StatusMethodNotAllowed, "method not allowed"))
			return
		}
		likesCount, err = h.postService.UnlikePost(r.Context(), userEmail, postID)
	default:
		response.Error(w, apperrors.New(http.StatusNotFound, "not found"))
		return
	}

	if err != nil {
		h.logger.Error("Failed to %s post: %v", action, err)
		response.Error(w, err)
		return
	}

	// Return the updated likes count
	resp := dto.LikesCountResponse{
		LikesCount: likesCount,
	}
	response.OK(w, resp)
}

// mapPostToDTO maps a post domain model to DTO
func (h *PostHandler) mapPostToDTO(p *post.Post) dto.PostResponse {
	return dto.PostResponse{
		ID:         p.ID,
		Author:     p.Author,
		Content:    p.Content,
		CreatedAt:  p.CreatedAt.Unix(),
		LikesCount: p.LikesCount,
	}
}
