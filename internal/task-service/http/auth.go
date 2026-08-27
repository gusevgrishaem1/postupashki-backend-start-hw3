package taskhttp

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"postupashki-backend-start-hw3/internal/task-service/usecases"
)

type credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// register godoc
// @Summary Register user
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body credentials true "Credentials"
// @Success 201
// @Failure 400 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Router /register [post]
func (h *Server) register(w http.ResponseWriter, r *http.Request) {
	credentials, ok := readCredentials(w, r)
	if !ok {
		return
	}
	if err := h.auth.Register(credentials.Username, credentials.Password); err != nil {
		if errors.Is(err, usecases.ErrUserExists) {
			write(w, http.StatusConflict, map[string]string{"error": err.Error()})
			return
		}
		write(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// login godoc
// @Summary Login user
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body credentials true "Credentials"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /login [post]
func (h *Server) login(w http.ResponseWriter, r *http.Request) {
	credentials, ok := readCredentials(w, r)
	if !ok {
		return
	}
	token, err := h.auth.Login(credentials.Username, credentials.Password)
	if err != nil {
		write(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}
	write(w, http.StatusOK, map[string]string{"token": token})
}

func readCredentials(w http.ResponseWriter, r *http.Request) (credentials, bool) {
	var value credentials
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&value); err != nil {
		write(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return credentials{}, false
	}
	value.Username = strings.TrimSpace(value.Username)
	if value.Username == "" || value.Password == "" {
		write(w, http.StatusBadRequest, map[string]string{"error": "username and password are required"})
		return credentials{}, false
	}
	return value, true
}
