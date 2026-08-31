package taskhttp

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"postupashki-backend-start-hw3/internal/task-service/usecases"
)

type credentialsRequest struct {
	Login    string `json:"login"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// register godoc
// @Summary Register user
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body credentialsRequest true "User credentials"
// @Success 201
// @Failure 400 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /register [post]
func (h *Server) register(w http.ResponseWriter, r *http.Request) {
	credentials, ok := decodeCredentials(w, r)
	if !ok {
		return
	}
	if err := h.auth.Register(credentials.Login, credentials.Password); err != nil {
		if errors.Is(err, usecases.ErrUserExists) {
			write(w, http.StatusConflict, map[string]string{"error": err.Error()})
			return
		}
		write(w, http.StatusInternalServerError, map[string]string{"error": "registration failed"})
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// login godoc
// @Summary Login user
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body credentialsRequest true "User credentials"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /login [post]
func (h *Server) login(w http.ResponseWriter, r *http.Request) {
	credentials, ok := decodeCredentials(w, r)
	if !ok {
		return
	}
	token, err := h.auth.Login(credentials.Login, credentials.Password)
	if err != nil {
		write(w, http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
		return
	}
	write(w, http.StatusOK, map[string]string{"token": token})
}

// logout godoc
// @Summary Logout user
// @Tags auth
// @Produce json
// @Security bearerAuth
// @Success 204
// @Failure 401 {object} map[string]string
// @Router /logout [post]
func (h *Server) logout(w http.ResponseWriter, r *http.Request) {
	token, ok := bearerToken(r)
	if !ok || h.auth.Logout(token) != nil {
		write(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func decodeCredentials(w http.ResponseWriter, r *http.Request) (credentialsRequest, bool) {
	const maxBodySize = 4 << 10
	r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)

	var request credentialsRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&request); err != nil {
		write(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return credentialsRequest{}, false
	}
	request.Login = strings.TrimSpace(request.Login)
	request.Password = strings.TrimSpace(request.Password)
	if request.Login == "" {
		request.Login = strings.TrimSpace(request.Username)
	}
	if request.Login == "" || request.Password == "" {
		write(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return credentialsRequest{}, false
	}
	return request, true
}
