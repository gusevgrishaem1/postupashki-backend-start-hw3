package http

import (
	"encoding/json"
	"net/http"
	swagger "postupashki-backend-start-hw3/pkg"
	"sync"

	"postupashki-backend-start-hw3/internal/domain"
	"postupashki-backend-start-hw3/internal/usecases"
)

type Task struct {
	usecase  usecases.Task
	mu       sync.RWMutex
	users    map[string]string
	sessions map[string]struct{}
}

func NewTask(usecase usecases.Task) *Task {
	return &Task{usecase: usecase, users: make(map[string]string), sessions: make(map[string]struct{})}
}

func (h *Task) Handler() http.Handler {
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

// create godoc
// @Summary Create task
// @Tags tasks
// @Produce json
// @Success 201 {object} map[string]string
// @Router /task [post]
func (h *Task) create(w http.ResponseWriter, _ *http.Request) {
	write(w, http.StatusCreated, map[string]string{"task_id": h.usecase.Create()})
}

// status godoc
// @Summary Get task status
// @Tags tasks
// @Produce json
// @Param task_id path string true "Task ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /status/{task_id} [get]
func (h *Task) status(w http.ResponseWriter, r *http.Request) {
	task, err := h.usecase.Get(r.PathValue("task_id"))
	if err != nil {
		write(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}
	write(w, http.StatusOK, map[string]string{"status": task.Status})
}

// result godoc
// @Summary Get task result
// @Tags tasks
// @Produce json
// @Param task_id path string true "Task ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Router /result/{task_id} [get]
func (h *Task) result(w http.ResponseWriter, r *http.Request) {
	task, err := h.usecase.Get(r.PathValue("task_id"))
	if err != nil {
		write(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}
	if task.Status != domain.Ready {
		write(w, http.StatusConflict, map[string]string{"error": "task is not ready"})
		return
	}
	write(w, http.StatusOK, map[string]string{"result": task.Result})
}

func write(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
