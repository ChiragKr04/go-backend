package user

import (
	"ChiragKr04/go-backend/types"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserServiceHandlers(t *testing.T) {
	userRepo := &mockUserRepo{}
	handler := NewHandler(userRepo)
	t.Run("fails if payload is invalid", func(t *testing.T) {
		// test payload is invalid
		payload := types.UserRegisterRequest{
			FirstName: "Chirag",
			LastName:  "Kr",
			Email:     "invalid",
			Password:  "password",
		}
		marshalJson, _ := json.Marshal(payload)
		req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(marshalJson))
		if err != nil {
			t.Fatal(err)
		}
		res := httptest.NewRecorder()
		router := http.NewServeMux()
		router.HandleFunc("/register", handler.handleRegister)
		router.ServeHTTP(res, req)
		if res.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, res.Code)
		}
	})

	t.Run("success if user gets created", func(t *testing.T) {
		// test user gets created
		payload := types.UserRegisterRequest{
			FirstName: "Chirag",
			LastName:  "Kr",
			Email:     "chirag@gmail.com",
			Password:  "password",
		}
		marshalJson, _ := json.Marshal(payload)
		req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(marshalJson))
		if err != nil {
			t.Fatal(err)
		}
		res := httptest.NewRecorder()
		router := http.NewServeMux()
		router.HandleFunc("/register", handler.handleRegister)
		router.ServeHTTP(res, req)
		if res.Code != http.StatusCreated {
			t.Errorf("expected status code %d, got %d", http.StatusCreated, res.Code)
		}
	})
}

type mockUserRepo struct {
}

// CreateUser implements types.UserRepository.
func (m *mockUserRepo) CreateUser(user types.User) error {
	return nil
}

// GetUserByEmail implements types.UserRepository.
func (m *mockUserRepo) GetUserByEmail(email string) (*types.User, error) {
	return nil, fmt.Errorf("user not found")
}

// GetUserByID implements types.UserRepository.
func (m *mockUserRepo) GetUserByID(id int) (*types.User, error) {
	return nil, fmt.Errorf("user not found")
}

// UpdateUser implements types.UserRepository.
func (m *mockUserRepo) UpdateUser(user types.User) (*types.User, error) {
	return nil, fmt.Errorf("user not found")
}
