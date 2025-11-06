package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"ynov-social-api/internal/storage"
	"ynov-social-api/internal/types"
)

type createPostRequest struct {
	Content string `json:"content"`
}

func (s *Server) CreatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req createPostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	content := strings.TrimSpace(req.Content)
	if content == "" || len([]rune(content)) > 400 {
		http.Error(w, "content too long", http.StatusBadRequest)
		return
	}

	username := UsernameFromContext(r)
	post, err := s.Posts.Create(username, content)
	if err != nil {
		http.Error(w, "failed to create post", http.StatusInternalServerError)
		return
	}

	// Ensure stable JSON response
	resp := types.Post{ID: post.ID, Author: post.Author, Content: post.Content, CreatedAt: post.CreatedAt, LikesCount: 0}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(resp)
}

// ListPosts returns paginated posts ordered by newest first.
// Query params: page (1-based), limit (default 10, max 100)
func (s *Server) ListPosts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	q := r.URL.Query()
	// Cursor-based pagination to avoid duplicates with infinite scroll
	beforeTsStr := q.Get("beforeTs")
	page, _ := strconv.Atoi(q.Get("page"))
	limit, _ := strconv.Atoi(q.Get("limit"))

	if beforeTsStr != "" {
		var beforeTs int64
		if v, err := strconv.ParseInt(beforeTsStr, 10, 64); err == nil {
			beforeTs = v
			posts, err := s.Posts.ListBefore(beforeTs, page, limit)
			if err != nil {
				http.Error(w, "failed to list posts", http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(posts)
			return
		}
	}

	// Default first page: use a cursor of "now" to return newest first
	nowTs := time.Now().Unix()
	posts, err := s.Posts.ListBefore(nowTs+1, page, limit)
	if err != nil {
		http.Error(w, "failed to list posts", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(posts)
}

// PostsAction handles subroutes like /posts/{id}/like
func (s *Server) PostsAction(w http.ResponseWriter, r *http.Request) {
	// Expect: /posts/{id}/like
	if r.Method != http.MethodPost && r.Method != http.MethodDelete {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	path := strings.TrimPrefix(r.URL.Path, "/posts/")
	parts := strings.Split(path, "/")
	if len(parts) != 2 || parts[1] != "like" && parts[1] != "unlike" || parts[0] == "" {
		http.NotFound(w, r)
		return
	}
	postID := parts[0]
	user := UsernameFromContext(r)
	if user == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	if r.Method == http.MethodPost {
		if err := s.Posts.Like(user, postID); err != nil {
			if err == storage.ErrPostNotFound {
				http.Error(w, "post not found", http.StatusNotFound)
				return
			}
			http.Error(w, "failed to like", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if r.Method == http.MethodDelete {
		if err := s.Posts.Unlike(user, postID); err != nil {
			if err == storage.ErrPostNotFound {
				http.Error(w, "post not found", http.StatusNotFound)
				return
			}
			http.Error(w, "failed to unlike", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}
