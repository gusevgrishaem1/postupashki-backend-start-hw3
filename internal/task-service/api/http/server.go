package taskhttp

import (
	"encoding/json"
	"net/http"
	"postupashki-backend-start-hw3/internal/task-service/api/http/swagger"
	"strings"

	"postupashki-backend-start-hw3/internal/task-service/usecases"
)

type Server struct {
	task usecases.Task
	auth usecases.Auth
}

func NewServer(task usecases.Task, auth usecases.Auth) *Server {
	return &Server{task: task, auth: auth}
}

func (h *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /register", h.register)
	mux.HandleFunc("POST /login", h.login)
	mux.HandleFunc("POST /logout", h.logout)
	mux.Handle("POST /task", h.requireAuth(http.HandlerFunc(h.create)))
	mux.HandleFunc("POST /commit", h.commit)
	mux.Handle("GET /status/{task_id}", h.requireAuth(http.HandlerFunc(h.status)))
	mux.Handle("GET /result/{task_id}", h.requireAuth(http.HandlerFunc(h.result)))
	mux.HandleFunc("GET /swagger.yaml", swagger.SpecificationHandler)
	mux.HandleFunc("GET /swagger", redirectToSwaggerUI)
	mux.HandleFunc("GET /swagger/", swagger.UIHandler)
	return mux
}

func (h *Server) requireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, ok := bearerToken(r)
		if !ok || h.auth.Authenticate(token) != nil {
			write(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
			return
		}
		next.ServeHTTP(w, r)
	})
}

func bearerToken(r *http.Request) (string, bool) {
	const prefix = "Bearer "
	header := r.Header.Get("Authorization")
	if !strings.HasPrefix(header, prefix) || len(header) == len(prefix) {
		return "", false
	}
	return header[len(prefix):], true
}

func redirectToSwaggerUI(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/swagger/", http.StatusMovedPermanently)
}

func write(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
