package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"postupashki-backend-start-hw3/internal/repository/inmemory"
	"postupashki-backend-start-hw3/internal/usecases/service"
)

func TestAuthenticationFlow(t *testing.T) {
	taskService := service.NewTask(inmemory.NewTask())
	authService := service.NewAuth(inmemory.NewUser(), inmemory.NewSession())
	handler := NewServer(taskService, authService).Handler()

	request := httptest.NewRequest(http.MethodPost, "/task", nil)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusUnauthorized {
		t.Fatalf("create task without token: got %d, want %d", response.Code, http.StatusUnauthorized)
	}

	credentials := []byte(`{"username":"alice","password":"secret"}`)
	request = httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(credentials))
	response = httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusCreated {
		t.Fatalf("register: got %d, want %d", response.Code, http.StatusCreated)
	}

	request = httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(credentials))
	response = httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("login: got %d, want %d", response.Code, http.StatusOK)
	}
	var body map[string]string
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if body["token"] == "" {
		t.Fatal("login returned an empty token")
	}

	request = httptest.NewRequest(http.MethodPost, "/task", nil)
	request.Header.Set("Authorization", "Bearer "+body["token"])
	response = httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusCreated {
		t.Fatalf("create task with token: got %d, want %d", response.Code, http.StatusCreated)
	}
}

func TestDuplicateRegistrationAndInvalidLogin(t *testing.T) {
	taskService := service.NewTask(inmemory.NewTask())
	authService := service.NewAuth(inmemory.NewUser(), inmemory.NewSession())
	handler := NewServer(taskService, authService).Handler()
	credentials := []byte(`{"username":"alice","password":"secret"}`)

	for index, expected := range []int{http.StatusCreated, http.StatusConflict} {
		request := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(credentials))
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, request)
		if response.Code != expected {
			t.Fatalf("register attempt %d: got %d, want %d", index+1, response.Code, expected)
		}
	}

	request := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(`{"username":"alice","password":"wrong"}`))
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusUnauthorized {
		t.Fatalf("invalid login: got %d, want %d", response.Code, http.StatusUnauthorized)
	}
}
