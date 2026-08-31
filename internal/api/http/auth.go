package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"postupashki-backend-start-hw3/internal/usecases"
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

// logout godoc
// @Summary Logout user
// @Tags auth
// @Security BearerAuth
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

func readCredentials(w http.ResponseWriter, r *http.Request) (credentials, bool) {
	const maxBodySize = 4 << 10 // 4 KiB

	r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)

	var value credentials
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&value); err != nil {
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			write(w, http.StatusRequestEntityTooLarge, map[string]string{"error": "request body is too large"})
			return credentials{}, false
		}
		write(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return credentials{}, false
	}
	value.Username = strings.TrimSpace(value.Username)
	value.Password = strings.TrimSpace(value.Password)
	if value.Username == "" || value.Password == "" {
		write(w, http.StatusBadRequest, map[string]string{"error": "username and password are required"})
		return credentials{}, false
	}
	return value, true
}
