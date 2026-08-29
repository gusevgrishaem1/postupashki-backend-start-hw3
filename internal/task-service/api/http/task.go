package taskhttp

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"postupashki-backend-start-hw3/internal/task-service/domain"
	"postupashki-backend-start-hw3/internal/task-service/usecases"
)

type createTaskRequest struct {
	Code     string `json:"code"`
	Language string `json:"language"`
	Runtime  string `json:"runtime"`
	Input    string `json:"input"`
}

type commitRequest struct {
	TaskID   string `json:"task_id"`
	Stdout   string `json:"stdout"`
	Stderr   string `json:"stderr"`
	ExitCode int    `json:"exit_code"`
}

// create godoc
// @Summary Create task
// @Tags tasks
// @Produce json
// @Success 201 {object} map[string]string
// @Router /task [post]
func (h *Server) create(w http.ResponseWriter, r *http.Request) {
	var request createTaskRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&request); err != nil {
		write(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	lang := request.Language
	if lang == "" {
		lang = request.Runtime
	}
	if strings.TrimSpace(request.Code) == "" || lang == "" {
		write(w, http.StatusBadRequest, map[string]string{"error": "code and runtime are required"})
		return
	}
	taskID, err := h.task.Create(request.Code, lang, request.Input)
	if err != nil {
		write(w, http.StatusServiceUnavailable, map[string]string{"error": usecases.ErrServiceUnavailable.Error()})
		return
	}
	write(w, http.StatusCreated, map[string]string{"task_id": taskID})
}

func (h *Server) commit(w http.ResponseWriter, r *http.Request) {
	var request commitRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&request); err != nil || request.TaskID == "" {
		write(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	if err := h.task.Commit(request.TaskID, domain.Result{Stdout: request.Stdout, Stderr: request.Stderr, ExitCode: request.ExitCode}); err != nil {
		if errors.Is(err, usecases.ErrNotFound) {
			write(w, http.StatusNotFound, map[string]string{"error": err.Error()})
			return
		}
		write(w, http.StatusServiceUnavailable, map[string]string{"error": usecases.ErrServiceUnavailable.Error()})
		return
	}
	w.WriteHeader(http.StatusNoContent)
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
		writeTaskReadError(w, err)
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
		writeTaskReadError(w, err)
		return
	}
	if task.Status != domain.Ready {
		write(w, http.StatusConflict, map[string]string{"error": "task is not ready"})
		return
	}
	write(w, http.StatusOK, map[string]domain.Result{"result": task.Result})
}

func writeTaskReadError(w http.ResponseWriter, err error) {
	if errors.Is(err, usecases.ErrNotFound) {
		write(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}
	write(w, http.StatusServiceUnavailable, map[string]string{"error": usecases.ErrServiceUnavailable.Error()})
}
