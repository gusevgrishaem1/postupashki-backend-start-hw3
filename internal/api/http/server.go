package http

import (
	"encoding/json"
	"net/http"
	"postupashki-backend-start-hw3/internal/usecases"
	swagger "postupashki-backend-start-hw3/pkg"
)

type Server struct {
	task usecases.Task
}

func NewServer(task usecases.Task) *Server {
	return &Server{task: task}
}

func (h *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /task", h.create)
	mux.HandleFunc("GET /status/{task_id}", h.status)
	mux.HandleFunc("GET /result/{task_id}", h.result)
	mux.HandleFunc("GET /swagger.yaml", swagger.SpecificationHandler)
	mux.HandleFunc("GET /swagger", redirectToSwaggerUI)
	mux.HandleFunc("GET /swagger/", swagger.UIHandler)
	return mux
}

func redirectToSwaggerUI(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/swagger/", http.StatusMovedPermanently)
}

func write(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
