package http

import (
	"net/http"
	"postupashki-backend-start-hw3/internal/domain"
)

// create godoc
// @Summary Create task
// @Tags tasks
// @Produce json
// @Success 201 {object} map[string]string
// @Router /task [post]
func (h *Server) create(w http.ResponseWriter, _ *http.Request) {
	write(w, http.StatusCreated, map[string]string{"task_id": h.task.Create()})
}

// status godoc
// @Summary Get task status
// @Tags tasks
// @Produce json
// @Param task_id path string true "Task ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /status/{task_id} [get]
func (h *Server) status(w http.ResponseWriter, r *http.Request) {
	task, err := h.task.Get(r.PathValue("task_id"))
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
func (h *Server) result(w http.ResponseWriter, r *http.Request) {
	task, err := h.task.Get(r.PathValue("task_id"))
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
