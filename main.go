package main

import (
	"log"
	"net/http"
	"os"

	"ynov-social-api/internal/handlers"
	"ynov-social-api/internal/storage"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET is not set")
		return
	}

	sqlite, err := storage.OpenSQLite("data.db")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer sqlite.Close()

	srv := handlers.NewServer(sqlite.Users(), sqlite.Posts(), []byte(secret))

	mux := http.NewServeMux()
	mux.HandleFunc("/signup", srv.Signup)
	mux.HandleFunc("/login", srv.Login)
	mux.Handle("/posts", srv.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			srv.ListPosts(w, r)
			return
		}
		srv.CreatePost(w, r)
	})))

	mux.Handle("/posts/", srv.AuthMiddleware(http.HandlerFunc(srv.PostsAction)))

	addr := ":8080"
	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, srv.ErrorJSONMiddleware(mux)); err != nil {
		log.Fatal(err)
	}
}
