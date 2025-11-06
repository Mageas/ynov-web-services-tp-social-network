package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"ynov-social-api/internal/auth"
)

type ctxKey string

const ctxKeyUsername ctxKey = "username"

// AuthMiddleware verifies Bearer JWT and injects username into request context.
func (s *Server) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authz := r.Header.Get("Authorization")
		parts := strings.SplitN(authz, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			http.Error(w, "missing or invalid Authorization header", http.StatusUnauthorized)
			return
		}
		username, err := auth.VerifyToken(parts[1], s.JWTSecret)
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), ctxKeyUsername, username)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// errorResponseWriter captures writes to allow JSON formatting of error responses.
type errorResponseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
	buf         bytes.Buffer
}

func (erw *errorResponseWriter) WriteHeader(code int) {
	erw.status = code
	erw.wroteHeader = true
	// For non-error statuses, write the header immediately so handlers can control status codes (e.g., 201)
	if code < 400 {
		erw.ResponseWriter.WriteHeader(code)
	}
}

func (erw *errorResponseWriter) Write(p []byte) (int, error) {
	if !erw.wroteHeader {
		// Default to 200 OK if Write is called first
		erw.WriteHeader(http.StatusOK)
	}
	if erw.status >= 400 {
		// Capture body for error messages; we'll emit JSON later
		return erw.buf.Write(p)
	}
	// For successful responses, header has already been written (in WriteHeader), so write body through
	return erw.ResponseWriter.Write(p)
}

type errorPayload struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// ErrorJSONMiddleware ensures error responses are returned as JSON {"status": xxx, "message": "..."}.
func (s *Server) ErrorJSONMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		erw := &errorResponseWriter{ResponseWriter: w}
		next.ServeHTTP(erw, r)

		// If status wasn't explicitly set, fall back to 200
		if !erw.wroteHeader {
			erw.status = http.StatusOK
		}

		if erw.status >= 400 {
			// Derive message from captured body or from status text
			msg := strings.TrimSpace(erw.buf.String())
			if msg == "" {
				msg = http.StatusText(erw.status)
			}

			// Ensure JSON content type and write header now
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(erw.status)
			_ = json.NewEncoder(w).Encode(errorPayload{Status: erw.status, Message: msg})
			return
		}
	})
}

// UsernameFromContext extracts the username set by AuthMiddleware.
func UsernameFromContext(r *http.Request) string {
	if v := r.Context().Value(ctxKeyUsername); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}
